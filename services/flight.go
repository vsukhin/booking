package services

import (
	"errors"
	"strings"
	"time"

	gorp "gopkg.in/gorp.v2"

	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/models"
	"github.com/vsukhin/booking/persistence/sqldb"
)

// FlightService is a flight service
type FlightService struct {
	db           sqldb.DBInterface
	blockService BlockServiceInterface
	seatService  SeatServiceInterface
}

// FlightServiceInterface is an interface for flight service methods
type FlightServiceInterface interface {
	Create(flight *models.Flight) error
	Retrieve(id int64) (*models.Flight, error)
	Delete(flight *models.Flight) error
	ListAll(filtering string, sorting string, limitation string) ([]models.Flight, error)
	GetMeta(filtering string) (*models.FlightMeta, error)
	SetBlocks(trans *gorp.Transaction, id int64, blocks []models.Block) error
}

// NewFlightService is a constructor for flight service
func NewFlightService(db sqldb.DBInterface, blockService BlockServiceInterface,
	seatService SeatServiceInterface) FlightServiceInterface {
	db.AddTableWithName(models.Flight{}, "flights").SetKeys(true, "ID")

	return &FlightService{db: db, blockService: blockService, seatService: seatService}
}

// Create creates flight
func (flightService *FlightService) Create(flight *models.Flight) error {
	trans, err := flightService.db.Begin()
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":  err,
			"flight": *flight,
		}).Error("Error creating transaction")
		return err
	}

	err = flightService.db.Insert(trans, flight)
	if err != nil {
		trErr := flightService.db.Rollback(trans)
		if trErr != nil {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":  trErr,
				"flight": *flight,
			}).Error("Error rollbacking transaction")
		}

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":  err,
			"flight": *flight,
		}).Error("Error creating flight")
		return err
	}

	err = flightService.SetBlocks(trans, flight.ID, flight.Blocks)
	if err != nil {
		trErr := flightService.db.Rollback(trans)
		if trErr != nil {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":  trErr,
				"flight": *flight,
			}).Error("Error rollbacking transaction")
		}
		return err
	}

	index := 0
	for _, block := range flight.Blocks {
		for i := 0; i < block.Rows; i++ {
			var seats []models.Seat

			line := 65
			for j := 0; j < block.SideSeatNumbers[0]; j++ {
				var seatType models.SeatType

				if j == block.SideSeatNumbers[0]-1 {
					seatType = models.SeatTypeAisle
				} else if j == 0 {
					seatType = models.SeatTypeWindow
				} else {
					seatType = models.SeatTypeMiddle
				}

				seats = append(seats, models.Seat{
					FlightID:  flight.ID,
					Index:     index + 1,
					Type:      seatType,
					Row:       i + 1,
					Line:      string(line),
					CreatedAt: time.Now().Unix(),
				})

				line++
				index++
			}

			for _, number := range block.MiddleSeatNumbers {
				for j := 0; j < number; j++ {
					var seatType models.SeatType

					if j == 0 || j == number-1 {
						seatType = models.SeatTypeAisle
					} else {
						seatType = models.SeatTypeMiddle
					}

					seats = append(seats, models.Seat{
						FlightID:  flight.ID,
						Index:     index + 1,
						Type:      seatType,
						Row:       i + 1,
						Line:      string(line),
						CreatedAt: time.Now().Unix(),
					})

					line++
					index++
				}
			}

			for j := 0; j < block.SideSeatNumbers[1]; j++ {
				var seatType models.SeatType

				if j == 0 {
					seatType = models.SeatTypeAisle
				} else if j == block.SideSeatNumbers[1]-1 {
					seatType = models.SeatTypeWindow
				} else {
					seatType = models.SeatTypeMiddle
				}

				seats = append(seats, models.Seat{
					FlightID:  flight.ID,
					Index:     index + 1,
					Type:      seatType,
					Row:       i + 1,
					Line:      string(line),
					CreatedAt: time.Now().Unix(),
				})

				line++
				index++
			}

			for j := range seats {
				err = flightService.seatService.Create(trans, &seats[j])
				if err != nil {
					trErr := flightService.db.Rollback(trans)
					if trErr != nil {
						logging.Log.WithFields(logging.DepthModerate, logging.Fields{
							"error":  trErr,
							"flight": *flight,
						}).Error("Error rollbacking transaction")
					}
					return err
				}
			}
		}
	}

	err = flightService.db.Commit(trans)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":  err,
			"flight": *flight,
		}).Error("Error committing transaction")
		return err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"flight": *flight,
	}).Debug("Flight successfully created")
	return nil
}

// Retrieve retrieves flight
func (flightService *FlightService) Retrieve(id int64) (*models.Flight, error) {
	obj, err := flightService.db.Get(nil, models.Flight{}, id)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error": err,
			"id":    id,
		}).Error("Error retrieving flight")
		return nil, err
	}

	if obj == nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"id": id,
		}).Error("Flight not found")
		return nil, nil
	}

	flight, ok := obj.(*models.Flight)
	if !ok {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"id":  id,
			"obj": obj,
		}).Error("Error returning flight")
		return nil, errors.New("Flight not valid")
	}

	flight.Blocks, err = flightService.blockService.ListAll(id)
	if err != nil {
		return nil, err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"id":     id,
		"flight": *flight,
	}).Debug("Flight successfully retrieved")
	return flight, nil
}

// Delete deletes flight
func (flightService *FlightService) Delete(flight *models.Flight) error {
	trans, err := flightService.db.Begin()
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":  err,
			"flight": *flight,
		}).Error("Error creating transaction")
		return err
	}

	err = flightService.seatService.DeleteAll(trans, flight.ID)
	if err != nil {
		trErr := flightService.db.Rollback(trans)
		if trErr != nil {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":  trErr,
				"flight": *flight,
			}).Error("Error rollbacking transaction")
		}
		return err
	}

	for i := range flight.Blocks {
		err = flightService.blockService.Delete(trans, &flight.Blocks[i])
		if err != nil {
			trErr := flightService.db.Rollback(trans)
			if trErr != nil {
				logging.Log.WithFields(logging.DepthModerate, logging.Fields{
					"error":  trErr,
					"flight": *flight,
				}).Error("Error rollbacking transaction")
			}
			return err
		}
	}

	_, err = flightService.db.Delete(trans, flight)
	if err != nil {
		trErr := flightService.db.Rollback(trans)
		if trErr != nil {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":  trErr,
				"flight": *flight,
			}).Error("Error rollbacking transaction")
		}

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":  err,
			"flight": *flight,
		}).Error("Error deleting flight")
		return err
	}

	err = flightService.db.Commit(trans)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":  err,
			"flight": *flight,
		}).Error("Error committing transaction")
		return err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"flight": *flight,
	}).Debug("Flight successfully deleted")
	return nil
}

// ListAll list all flights according filtering, sorting, limitation parameters
func (flightService *FlightService) ListAll(filtering string, sorting string,
	limitation string) ([]models.Flight, error) {
	var flights []models.Flight

	if filtering != "" {
		filtering = " WHERE " + strings.TrimPrefix(filtering, " AND ")
	}

	_, err := flightService.db.Select(&flights, "SELECT * FROM flights"+filtering+sorting+limitation)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":      err,
			"filtering":  filtering,
			"sorting":    sorting,
			"limitation": limitation,
		}).Error("Error returning flights")
		return nil, err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"filtering":  filtering,
		"sorting":    sorting,
		"limitation": limitation,
		"flights":    flights,
	}).Debug("Flights successfully returned")
	return flights, nil
}

// GetMeta gets metadata about flight list according filtering parameters
func (flightService *FlightService) GetMeta(filtering string) (*models.FlightMeta, error) {
	if filtering != "" {
		filtering = " WHERE " + strings.TrimPrefix(filtering, " AND ")
	}

	count, err := flightService.db.SelectInt("SELECT COUNT(*) FROM flights" + filtering)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":     err,
			"filtering": filtering,
		}).Error("Error returning flight metadata")
		return nil, err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"filtering": filtering,
		"count":     count,
	}).Debug("Flight metadata successfully returned")
	return &models.FlightMeta{
		TotalRecords: count,
	}, nil
}

// SetBlocks sets flight blocks
func (flightService *FlightService) SetBlocks(trans *gorp.Transaction, id int64, blocks []models.Block) error {
	for i := range blocks {
		blocks[i].FlightID = id
		err := flightService.blockService.Create(trans, &blocks[i])
		if err != nil {
			return err
		}
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"id":     id,
		"blocks": blocks,
	}).Debug("Flight blocks successfully inserted")
	return nil
}
