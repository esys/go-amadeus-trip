package repository

import (
	"amadeus-trip-parser/internal/domain"
	"amadeus-trip-parser/internal/domain/model"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type sqliteTripRepo struct {
	db *gorm.DB
}

func NewSQLiteTripRepo(db *sql.DB) (domain.TripRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("failed due to nil DB pointer")
	}
	gdb, err := gorm.Open("sqlite3", db)
	if err != nil {
		return nil, fmt.Errorf("cannot open DB connection: %w", err)
	}
	gdb.AutoMigrate(&model.Trip{})
	gdb.AutoMigrate(&model.TripStep{})
	return &sqliteTripRepo{gdb}, nil
}

func (s *sqliteTripRepo) GetAll() ([]model.Trip, error) {
	var trips []model.Trip
	if dbc := s.db.Preload("TripSteps").Find(&trips); dbc.Error != nil {
		return nil, fmt.Errorf("failed database query for getting all trips: %w", dbc.Error)
	}
	return trips, nil
}

func (s *sqliteTripRepo) GetOne(query model.Trip) (model.Trip, error) {
	var trip model.Trip
	if dbc := s.db.Preload("TripSteps").Where(query).First(&trip); dbc.Error != nil {
		if dbc.Error == gorm.ErrRecordNotFound {
			return model.Trip{}, domain.ErrorNotFound
		}
		return model.Trip{}, fmt.Errorf("failed database query when looking for trip with query %s: %w", query, dbc.Error)
	}
	return trip, nil
}

func (s *sqliteTripRepo) Create(trip *model.Trip) error {
	if dbc := s.db.Create(&trip); dbc.Error != nil {
		return fmt.Errorf("failed creating new trip in repository: %w", dbc.Error)
	}
	return nil
}
