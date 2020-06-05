package api

import (
	"amadeus-trip-parser/internal/domain"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	. "net/http"
)

type TripAPI interface {
	Get(c echo.Context) error
}

type tripAPI struct {
	tripFinder domain.TripFinder
}

func NewTripAPI(tripFinder domain.TripFinder) TripAPI {
	return &tripAPI{tripFinder: tripFinder}
}

func (a *tripAPI) Get(c echo.Context) error {
	ref := c.QueryParam("ref")
	if ref == "" {
		trips, err := a.tripFinder.Get()
		if err != nil {
			return err
		}
		return c.JSON(StatusOK, trips)
	}

	trip, err := a.tripFinder.GetByReference(ref)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrorNotFound):
			return echo.NewHTTPError(StatusNotFound, fmt.Sprintf("no trip with reference %s", ref))
		default:
			return echo.NewHTTPError(StatusInternalServerError, err)
		}
	}
	return c.JSON(StatusOK, trip)
}
