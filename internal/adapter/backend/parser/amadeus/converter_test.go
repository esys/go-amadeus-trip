package amadeus

import (
	"amadeus-trip-parser/internal/domain/model"
	"encoding/json"
	"testing"
	"time"
)

func readResponseData(file string, t *testing.T) resultResponseData {
	byt := readTestData(file, t)
	var r resultResponse
	if err := json.Unmarshal(byt, &r); err != nil {
		t.Errorf("failed to load data to convert: %s", err)
		return resultResponseData{}
	}
	return r.Data
}

func compareTrip(a model.Trip, b model.Trip, t *testing.T) bool {
	if a.Reference != b.Reference {
		t.Logf("reference are different: %s != %s", a.Reference, b.Reference)
		return false
	}
	if len(a.TripSteps) != len(b.TripSteps) {
		t.Logf("not the same number of steps: %d != %d", len(a.TripSteps), len(b.TripSteps))
		return false
	}
	for i, ats := range a.TripSteps {
		bts := b.TripSteps[i]
		if ats.Type != bts.Type {
			t.Logf("step types are different: %s != %s", ats.Type, bts.Type)
			return false
		}
		if ats.DateTime != bts.DateTime {
			t.Logf("step times are different: %s != %s", ats.DateTime, bts.DateTime)
			return false
		}
	}
	return true
}

func Test_addrConverter_String(t *testing.T) {
	type args struct {
		a Address
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "address with text",
			args: args{
				Address{
					Text:        "ADDRESS",
					Lines:       []string{"LINE1, LINE2"},
					CityName:    "PARIS",
					CountryName: "FRANCE",
				},
			},
			want: "ADDRESS",
		},
		{
			name: "address without text but with address",
			args: args{
				Address{
					Lines:       []string{"LINE1, LINE2"},
					CityName:    "PARIS",
					CountryName: "FRANCE",
				},
			},
			want: "LINE1, LINE2",
		},
		{
			name: "address with only city or country name",
			args: args{
				Address{
					CountryName: "FRANCE",
				},
			},
			want: "FRANCE",
		},
		{
			name: "address with both city and country name",
			args: args{
				Address{
					CityName:    "PARIS",
					CountryName: "FRANCE",
				},
			},
			want: "PARIS, FRANCE",
		},
		{
			name: "empty address",
			args: args{
				Address{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := addrConverter{}
			if got := a.String(tt.args.a); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_converter_getTrip(t *testing.T) {
	type args struct {
		d resultResponseData
	}
	tests := []struct {
		name    string
		args    args
		want    model.Trip
		wantErr bool
	}{
		{
			name: "air product",
			args: args{
				readResponseData("testdata/air.json", t),
			},
			want: model.Trip{
				Reference: "JKW499",
				TripSteps: []model.TripStep{
					{
						Type:     model.TripStepTypeFlightStart,
						DateTime: time.Date(2020, 04, 06, 16, 10, 00, 0, time.UTC),
					},
					{
						Type:     model.TripStepTypeFlightEnd,
						DateTime: time.Date(2020, 04, 06, 17, 45, 00, 0, time.UTC),
					},
					{
						Type:     model.TripStepTypeFlightStart,
						DateTime: time.Date(2020, 04, 12, 11, 55, 00, 0, time.UTC),
					},
					{
						Type:     model.TripStepTypeFlightEnd,
						DateTime: time.Date(2021, 04, 12, 15, 30, 00, 0, time.UTC),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := converter{}
			got, err := c.getTrip(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareTrip(got, tt.want, t) {
				t.Errorf("getTrip() got = %v, want %v", got, tt.want)
			}
		})
	}
}
