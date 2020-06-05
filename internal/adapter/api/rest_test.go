package api

import (
	"amadeus-trip-parser/internal/domain"
	"amadeus-trip-parser/internal/domain/mocks"
	"amadeus-trip-parser/internal/domain/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var trip = []model.Trip{
	{
		ID:        "ID0",
		Reference: "REF0",
		TripSteps: []model.TripStep{
			{
				ID:          "IDS00",
				TripID:      "ID0",
				Type:        model.TripStepTypeFlightStart,
				Location:    "PARIS",
				Description: "DESC00",
			},
			{
				ID:          "IDS01",
				TripID:      "ID0",
				Type:        model.TripStepTypeFlightEnd,
				Location:    "ROME",
				Description: "DESC01",
			},
		},
	},
}

var tripJSON = `[{"ID":"ID0","Reference":"REF0","Start":"0001-01-01T00:00:00Z","End":"0001-01-01T00:00:00Z","TripSteps":[{"ID":"IDS00","TripID":"ID0","Type":"flight-start","DateTime":"0001-01-01T00:00:00Z","Location":"PARIS","Description":"DESC00"},{"ID":"IDS01","TripID":"ID0","Type":"flight-end","DateTime":"0001-01-01T00:00:00Z","Location":"ROME","Description":"DESC01"}]}]
`

func Test_tripAPI_Get(t *testing.T) {
	mockFinder := &mocks.TripFinder{}
	mockFinder.On("Get").Return(trip, nil)
	mockFinder.On("GetByReference", "1111").Return(model.Trip{}, domain.ErrorNotFound)
	e := echo.New()

	type fields struct {
		tripFinder domain.TripFinder
	}
	tests := []struct {
		name    string
		fields  fields
		path    string
		code    int
		body    string
		wantErr bool
	}{
		{
			"get all",
			fields{mockFinder},
			"/trip",
			200,
			tripJSON,
			false,
		},
		{
			"get 404",
			fields{mockFinder},
			"/trip?ref=1111",
			404,
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &tripAPI{
				tripFinder: tt.fields.tripFinder,
			}
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath(tt.path)

			err := a.Get(ctx)
			if tt.wantErr {
				assert.IsType(t, &echo.HTTPError{}, err)
				herr := err.(*echo.HTTPError)
				// rec will will say code 200, so test expected code on HTTPError directly
				assert.Equal(t, herr.Code, tt.code)
			} else {
				assert.Equal(t, tt.code, rec.Code)
				assert.Equal(t, tt.body, rec.Body.String())
			}

		})
	}
}
