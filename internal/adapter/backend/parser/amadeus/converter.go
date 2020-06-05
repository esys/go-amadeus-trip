package amadeus

import (
	"amadeus-trip-parser/internal/domain/model"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type addressConverter interface {
	String(a Address) string
}

type addrConverter struct{}

func NewAddressConverter() addressConverter {
	return &addrConverter{}
}

func (addrConverter) String(a Address) string {
	var txt string
	switch {
	case a.Text != "":
		txt = a.Text
	case len(a.Lines) > 0:
		txt = strings.Join(a.Lines, ", ")
	case a.CityName != "" || a.CountryName != "":
		var s []string
		if a.CityName != "" {
			s = append(s, a.CityName)
		}
		if a.CountryName != "" {
			s = append(s, a.CountryName)
		}
		txt = strings.Join(s, ",")
	default:
		txt = ""
	}
	return txt
}

type tripConverter interface {
	getTrip(resultResponseData) (model.Trip, error)
}

type converter struct {
}

func newTripConverter() tripConverter {
	return &converter{}
}

func (c converter) getTrip(d resultResponseData) (model.Trip, error) {
	conv := newTripStepConverter()
	var steps []model.TripStep
	for _, p := range d.Products {
		ts, err := conv.getTripStep(p)
		if err != nil {
			return model.Trip{}, fmt.Errorf("failed to convert %v to trip: %w", d, err)
		}
		steps = append(steps, ts...)
	}

	return model.Trip{
		ID:        uuid.New().String(),
		Reference: d.Reference,
		Start:     d.Start.DateTime.Time,
		End:       d.End.DateTime.Time,
		TripSteps: steps,
	}, nil
}

type tripStepConverter interface {
	getTripStep(interface{}) ([]model.TripStep, error)
}
type stepConverter struct {
}

func newTripStepConverter() tripStepConverter {
	return &stepConverter{}
}

func (s stepConverter) getTripStep(i interface{}) ([]model.TripStep, error) {
	addrConv := NewAddressConverter()
	switch i.(type) {
	case Product:
		p := i.(Product)
		switch {
		case p.Air != nil:
			return s.getTripStep(p.Air)
		case p.Hotel != nil:
			return s.getTripStep(p.Hotel)
		default:
			return []model.TripStep{}, nil
		}
	case *AirProduct:
		a := i.(*AirProduct)
		return []model.TripStep{
			{
				ID:          uuid.New().String(),
				Type:        model.TripStepTypeFlightStart,
				DateTime:    a.Start.DateTime.Time,
				Location:    addrConv.String(a.Start.Address),
				Description: fmt.Sprintf("Flight start with %s", a.ServiceProvider.Name),
			},
			{
				ID:          uuid.New().String(),
				Type:        model.TripStepTypeFlightEnd,
				DateTime:    a.End.DateTime.Time,
				Location:    addrConv.String(a.End.Address),
				Description: fmt.Sprintf("Flight end with %s", a.ServiceProvider.Name),
			},
		}, nil
	case *HotelProduct:
		h := i.(*HotelProduct)
		return []model.TripStep{{
			ID:          uuid.New().String(),
			Type:        model.TripStepTypeHotel,
			DateTime:    h.Start.DateTime.Time,
			Location:    addrConv.String(h.Start.Address),
			Description: fmt.Sprintf("Hotel at %s", h.ServiceProvider.Name),
		}}, nil
	default:
		return []model.TripStep{}, fmt.Errorf("cannot convert type %T to TripStep", i)
	}
}
