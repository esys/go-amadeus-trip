package usecase

import (
	"amadeus-trip-parser/internal/domain"
	"amadeus-trip-parser/internal/domain/mocks"
	"amadeus-trip-parser/internal/domain/model"
	"errors"
	"reflect"
	"testing"
	"time"
)

var trip = []model.Trip{
	{
		ID:        "ID0",
		Reference: "REF0",
		Start:     time.Now(),
		End:       time.Now().Add(36 * time.Hour),
		TripSteps: []model.TripStep{
			{
				ID:          "IDS00",
				TripID:      "ID0",
				Type:        model.TripStepTypeFlightStart,
				DateTime:    time.Now(),
				Location:    "PARIS",
				Description: "DESC00",
			},
			{
				ID:          "IDS01",
				TripID:      "ID0",
				Type:        model.TripStepTypeFlightEnd,
				DateTime:    time.Now().Add(36 * time.Hour),
				Location:    "ROME",
				Description: "DESC01",
			},
		},
	},
}

func Test_tripFinder_Get(t *testing.T) {
	mockRepo := &mocks.TripRepository{}
	mockRepo.On("GetAll").Return(trip, nil)

	mockRepoErr := &mocks.TripRepository{}
	mockRepoErr.On("GetAll").Return(nil, errors.New("error"))

	type fields struct {
		repo domain.TripRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Trip
		wantErr bool
	}{
		{
			"get all trips",
			fields{mockRepo},
			trip,
			false,
		},
		{
			"get all trips with error",
			fields{mockRepoErr},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tripFinder{
				repo: tt.fields.repo,
			}
			got, err := f.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tripFinder_GetByReference(t *testing.T) {
	mockRepo := &mocks.TripRepository{}
	mockRepo.On("GetOne", model.Trip{Reference: trip[0].Reference}).Return(trip[0], nil)

	mockRepoErr := &mocks.TripRepository{}
	mockRepoErr.On("GetOne", model.Trip{Reference: trip[0].Reference}).
		Return(model.Trip{}, errors.New("error"))

	type fields struct {
		repo domain.TripRepository
	}
	type args struct {
		ref string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Trip
		wantErr bool
	}{
		{
			"get a trip by reference",
			fields{mockRepo},
			args{trip[0].Reference},
			trip[0],
			false,
		},
		{
			"get a trip with error",
			fields{mockRepoErr},
			args{trip[0].Reference},
			model.Trip{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tripFinder{
				repo: tt.fields.repo,
			}
			got, err := f.GetByReference(tt.args.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByReference() got = %v, want %v", got, tt.want)
			}
		})
	}
}
