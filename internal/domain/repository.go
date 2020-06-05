package domain

import (
	"amadeus-trip-parser/internal/domain/model"
	"errors"
)

var ErrorNotFound = errors.New("trip not found")

type TripRepository interface {
	GetAll() ([]model.Trip, error)
	GetOne(query model.Trip) (model.Trip, error)
	Create(trip *model.Trip) error
}
