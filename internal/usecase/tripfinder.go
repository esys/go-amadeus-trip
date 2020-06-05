package usecase

import (
	"amadeus-trip-parser/internal/domain"
	"amadeus-trip-parser/internal/domain/model"
)

type tripFinder struct {
	repo domain.TripRepository
}

func NewTripFinder(repo domain.TripRepository) domain.TripFinder {
	return &tripFinder{repo}
}

func (f tripFinder) GetByReference(ref string) (model.Trip, error) {
	return f.repo.GetOne(model.Trip{Reference: ref})
}

func (f tripFinder) Get() ([]model.Trip, error) {
	return f.repo.GetAll()
}
