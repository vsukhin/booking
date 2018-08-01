package models

import (
	"fmt"
)

// BlockType is block type
type BlockType int

const (
	// BlockTypeSide is side block type
	BlockTypeSide BlockType = iota + 1
	// BlockTypeMiddle is middle block type
	BlockTypeMiddle
)

// Block is seat block
type Block struct {
	ID                int64 `json:"-"                   db:"id"`
	FlightID          int64 `json:"-"                   db:"flight_id"`
	Rows              int   `json:"rows"                db:"rows"`
	SideSeatNumbers   []int `json:"side_seat_numbers"   db:"-"`
	MiddleSeatNumbers []int `json:"middle_seat_numbers" db:"-"`
}

// Validate validates block data
func (block *Block) Validate() []Error {
	var errs []Error

	if block.Rows <= 0 {
		errs = append(errs, Error{
			Code:    "rows.TooSmall",
			Message: "Rows must be more than zero",
			Field:   "rows",
		})
	}

	if len(block.SideSeatNumbers) != 2 {
		errs = append(errs, Error{
			Code:    "side_seat_numbers.Invalid",
			Message: "Must be precisely two side seats",
			Field:   "side_seat_numbers",
		})
	}

	lines := 0
	for _, number := range block.SideSeatNumbers {
		if number <= 0 {
			errs = append(errs, Error{
				Code:    "side_seat_numbers.TooSmall",
				Message: "Side seats must be more than zero",
				Field:   "side_seat_numbers",
			})
		}

		lines += number
	}

	for _, number := range block.MiddleSeatNumbers {
		if number <= 0 {
			errs = append(errs, Error{
				Code:    "middle_seat_numbers.TooSmall",
				Message: "Middle seats must be more than zero",
				Field:   "middle_seat_numbers",
			})
		}

		lines += number
	}

	if lines > maxLines {
		errs = append(errs, Error{
			Code:    "seat_number.TooLarge",
			Message: fmt.Sprintf("Seat number must be less than %v", maxLines),
			Field:   "side_seat_numbers,middle_seat_numbers",
		})
	}

	return errs
}
