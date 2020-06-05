package model

import "time"

type TripStepType string

const (
	TripStepTypeFlightStart = "flight-start"
	TripStepTypeFlightEnd   = "flight-end"
	TripStepTypeHotel       = "hotel"
)

type TripStep struct {
	ID          string
	TripID      string
	Type        TripStepType
	DateTime    time.Time
	Location    string
	Description string
}

type Trip struct {
	ID        string
	Reference string
	Start     time.Time
	End       time.Time
	TripSteps []TripStep
}
