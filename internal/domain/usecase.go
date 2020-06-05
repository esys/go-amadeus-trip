package domain

import "amadeus-trip-parser/internal/domain/model"

type EmailProcessor interface {
	Process()
	Stop()
}

type TripFinder interface {
	Get() ([]model.Trip, error)
	GetByReference(ref string) (model.Trip, error)
}