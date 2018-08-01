package services

import (
	"fmt"

	gorp "gopkg.in/gorp.v2"

	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/models"
	"github.com/vsukhin/booking/persistence/sqldb"
)

// BlockService is a block service
type BlockService struct {
	db sqldb.DBInterface
}

// BlockServiceInterface is an interface for block service methods
type BlockServiceInterface interface {
	Create(trans *gorp.Transaction, block *models.Block) error
	Delete(trans *gorp.Transaction, block *models.Block) error
	ListAll(flightID int64) ([]models.Block, error)
}

// NewBlockService is a constructor for block service
func NewBlockService(db sqldb.DBInterface) BlockServiceInterface {
	db.AddTableWithName(models.Block{}, "blocks").SetKeys(true, "ID")

	return &BlockService{db: db}
}

// Create creates block
func (blockService *BlockService) Create(trans *gorp.Transaction, block *models.Block) error {
	err := blockService.db.Insert(trans, block)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error": err,
			"block": *block,
		}).Error("Error creating block")
		return err
	}

	var seatNumbers = []struct {
		numbers   []int
		blockType models.BlockType
	}{
		{
			block.SideSeatNumbers,
			models.BlockTypeSide,
		},
		{
			block.MiddleSeatNumbers,
			models.BlockTypeMiddle,
		},
	}

	for _, seats := range seatNumbers {
		if len(seats.numbers) > 0 {
			statement := ""
			for _, number := range seats.numbers {
				if statement != "" {
					statement += " UNION ALL"
				}
				statement += fmt.Sprintf(" SELECT %v, %v, %v", block.ID, seats.blockType, number)
			}

			_, err = blockService.db.Exec(trans, "INSERT INTO seat_numbers (block_id, type, number)"+statement)
			if err != nil {
				logging.Log.WithFields(logging.DepthModerate, logging.Fields{
					"error": err,
					"block": *block,
				}).Error("Error inserting seat numbers")
				return err
			}
		}
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"block": block,
	}).Debug("Block successfully created")
	return nil
}

// Delete deletes block
func (blockService *BlockService) Delete(trans *gorp.Transaction, block *models.Block) error {
	_, err := blockService.db.Exec(trans, "DELETE FROM seat_numbers WHERE block_id = ?", block.ID)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error": err,
			"block": *block,
		}).Error("Error deleting seat numbers")
		return err
	}

	_, err = blockService.db.Delete(trans, block)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error": err,
			"block": *block,
		}).Error("Error deleting block")
		return err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"block": *block,
	}).Debug("Block successfully deleted")
	return nil
}

// ListAll list all blocks according filtering, sorting, limitation parameters
func (blockService *BlockService) ListAll(flightID int64) ([]models.Block, error) {
	var blocks []models.Block

	_, err := blockService.db.Select(&blocks, "SELECT * FROM blocks WHERE flight_id = ? ORDER BY id", flightID)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":    err,
			"flightID": flightID,
		}).Error("Error returning blocks")
		return nil, err
	}

	for i := range blocks {
		_, err = blockService.db.Select(&blocks[i].SideSeatNumbers, "SELECT number FROM seat_numbers "+
			"WHERE block_id = ? AND type = ? ORDER BY id", blocks[i].ID, models.BlockTypeSide)
		if err != nil {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":    err,
				"flightID": flightID,
			}).Error("Error returning side seat number")
			return nil, err
		}

		_, err = blockService.db.Select(&blocks[i].MiddleSeatNumbers, "SELECT number FROM seat_numbers "+
			"WHERE block_id = ? AND type = ? ORDER BY id", blocks[i].ID, models.BlockTypeMiddle)
		if err != nil {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":    err,
				"flightID": flightID,
			}).Error("Error returning middle seat number")
			return nil, err
		}
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"flightID": flightID,
		"blocks":   blocks,
	}).Debug("Blocks successfully returned")
	return blocks, nil
}
