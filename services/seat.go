package services

import (
	"database/sql"
	"time"

	gorp "gopkg.in/gorp.v2"

	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/models"
	"github.com/vsukhin/booking/persistence/sqldb"
)

// SeatService is a seat service
type SeatService struct {
	db sqldb.DBInterface
}

// SeatServiceInterface is an interface for seat service methods
type SeatServiceInterface interface {
	Create(trans *gorp.Transaction, seat *models.Seat) error
	Assign(flightID int64) (*models.Seat, error)
	Update(seat *models.Seat) error
	DeleteAll(trans *gorp.Transaction, flightID int64) error
	Retrieve(flightID int64, index int64) (*models.Seat, error)
	Find(flightID int64, row int, line string) (*models.Seat, error)
	ListAll(flightID int64, filtering string, sorting string, limitation string) ([]models.Seat, error)
	GetMeta(flightID int64, filtering string) (*models.SeatMeta, error)
}

// NewSeatService is a constructor for seat service
func NewSeatService(db sqldb.DBInterface) SeatServiceInterface {
	db.AddTableWithName(models.Seat{}, "seats").SetKeys(true, "ID")

	return &SeatService{db: db}
}

// Create creates seat
func (seatService *SeatService) Create(trans *gorp.Transaction, seat *models.Seat) error {
	err := seatService.db.Insert(trans, seat)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error": err,
			"seat":  *seat,
		}).Error("Error creating seat")
		return err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"seat": *seat,
	}).Debug("Seat successfully created")
	return nil
}

// Assign assignes seat
func (seatService *SeatService) Assign(flightID int64) (*models.Seat, error) {
	var seat models.Seat

	trans, err := seatService.db.Begin()
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":    err,
			"flightID": flightID,
		}).Error("Error creating transaction")
		return nil, err
	}

	err = seatService.db.SelectOne(nil, &seat, "SELECT * FROM seats WHERE flight_id = ? AND assigned = false "+
		"ORDER BY row ASC, type ASC, line ASC LIMIT 1", flightID)
	if err != nil {
		trErr := seatService.db.Rollback(trans)
		if trErr != nil {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":    trErr,
				"flightID": flightID,
			}).Error("Error rollbacking transaction")
		}

		if err == sql.ErrNoRows {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"flightID": flightID,
			}).Error("Seat not found")
			return nil, nil
		}

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":    err,
			"flightID": flightID,
		}).Error("Error returning seat")
		return nil, err
	}

	seat.Assigned = true
	seat.UpdatedAt = time.Now().Unix()
	_, err = seatService.db.Update(trans, &seat)
	if err != nil {
		trErr := seatService.db.Rollback(trans)
		if trErr != nil {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":    trErr,
				"flightID": flightID,
			}).Error("Error rollbacking transaction")
		}

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":    err,
			"flightID": flightID,
		}).Error("Error reating seat")
		return nil, err
	}

	err = seatService.db.Commit(trans)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":    err,
			"flightID": flightID,
		}).Error("Error committing transaction")
		return nil, err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"flightID": flightID,
		"seat":     seat,
	}).Debug("Seat successfully assigned")
	return &seat, nil
}

// Update updates seat
func (seatService *SeatService) Update(seat *models.Seat) error {
	_, err := seatService.db.Update(nil, seat)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error": err,
			"seat":  *seat,
		}).Error("Error updating seat")
		return err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"seat": *seat,
	}).Debug("Seat successfully updated")
	return nil
}

// DeleteAll deletes all seats
func (seatService *SeatService) DeleteAll(trans *gorp.Transaction, flightID int64) error {
	_, err := seatService.db.Exec(trans, "DELETE FROM seats WHERE flight_id = ?", flightID)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":    err,
			"flightID": flightID,
		}).Error("Error deleting all seats")
		return err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"flightID": flightID,
	}).Debug("All seats successfully deleted")
	return nil
}

// Retrieve retrieves seat
func (seatService *SeatService) Retrieve(flightID int64, index int64) (*models.Seat, error) {
	var seat models.Seat

	err := seatService.db.SelectOne(nil, &seat, "SELECT * FROM seats WHERE flight_id = ? AND `index` = ?",
		flightID, index)
	if err != nil {
		if err == sql.ErrNoRows {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"flightID": flightID,
				"index":    index,
			}).Error("Seat not found")
			return nil, nil
		}

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":    err,
			"flightID": flightID,
			"index":    index,
		}).Error("Error returning seat")
		return nil, err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"flightID": flightID,
		"index":    index,
		"seat":     seat,
	}).Debug("Seat successfully retrieved")
	return &seat, nil
}

// Find finds seat
func (seatService *SeatService) Find(flightID int64, row int, line string) (*models.Seat, error) {
	var seat models.Seat

	err := seatService.db.SelectOne(nil, &seat, "SELECT * FROM seats WHERE flight_id = ? AND row = ? AND line = ?",
		flightID, row, line)
	if err != nil {
		if err == sql.ErrNoRows {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"flightID": flightID,
				"row":      row,
				"line":     line,
			}).Error("Seat not found")
			return nil, nil
		}

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":    err,
			"flightID": flightID,
			"row":      row,
			"line":     line,
		}).Error("Error returning seat")
		return nil, err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"flightID": flightID,
		"row":      row,
		"line":     line,
		"seat":     seat,
	}).Debug("Seat successfully found")
	return &seat, nil
}

// ListAll list all seats according filtering, sorting, limitation parameters
func (seatService *SeatService) ListAll(flightID int64, filtering string, sorting string,
	limitation string) ([]models.Seat, error) {
	var seats []models.Seat

	_, err := seatService.db.Select(&seats, "SELECT * FROM seats WHERE flight_id = ?"+filtering+sorting+limitation,
		flightID)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":      err,
			"flightID":   flightID,
			"filtering":  filtering,
			"sorting":    sorting,
			"limitation": limitation,
		}).Error("Error returning seats")
		return nil, err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"flightID":   flightID,
		"filtering":  filtering,
		"sorting":    sorting,
		"limitation": limitation,
		"seats":      seats,
	}).Debug("Seats successfully returned")
	return seats, nil
}

// GetMeta gets metadata about seat list according filtering parameters
func (seatService *SeatService) GetMeta(flightID int64, filtering string) (*models.SeatMeta, error) {
	count, err := seatService.db.SelectInt("SELECT COUNT(*) FROM seats WHERE flight_id = ?"+filtering, flightID)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":     err,
			"flightID":  flightID,
			"filtering": filtering,
		}).Error("Error returning seat metadata")
		return nil, err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"flightID":  flightID,
		"filtering": filtering,
		"count":     count,
	}).Debug("Seat metadata successfully returned")
	return &models.SeatMeta{
		TotalRecords: count,
	}, nil
}
