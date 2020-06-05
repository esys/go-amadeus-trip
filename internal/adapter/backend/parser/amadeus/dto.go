package amadeus

import (
	"strings"
	"time"
)

// Common
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	str := strings.TrimSuffix(string(data), "\"")
	str = strings.TrimPrefix(str, "\"")
	tt, err := time.Parse("2006-01-02T15:04:05", str)
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

type CommonResponse struct {
	Warnings []apiError `json:"warnings"`
	Errors   []apiError `json:"errors"`
}

type CommonResponseData struct {
	TypeP string `json:"type"`
	ID    string `json:"id"`
}

// Authorize
type authorizeResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"` // seconds
}

// Create
type createRequest struct {
	Data createRequestData `json:"data"`
}

type createRequestData struct {
	TypeP   string `json:"type"`
	Content string `json:"content"`
}

type createResponse struct {
	CommonResponse
	Data createResponseData `json:"data"`
}

type createResponseData struct {
	CommonResponseData
	Status Status `json:"status"`
}

// Status
type statusResponse struct {
	CommonResponse
	Data statusResponseData `json:"data"`
}

type statusResponseData struct {
	CommonResponseData
	Status Status `json:"status"`
	Detail string `json:"detail"`
}

type apiError struct {
	Status int    `json:"status"`
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// Result
type resultResponse struct {
	CommonResponse
	Data resultResponseData `json:"data"`
}

type resultResponseData struct {
	CommonResponseData
	Reference    string    `json:"reference"`
	Start        TripPoint `json:"start"`
	End          TripPoint `json:"end"`
	Stakeholders []struct {
		PTC         string    `json:"PTC"`
		DateOfBirth time.Time `json:"dateOfBirth"`
		Roles       []string  `json:"roles"`
		Names       []struct {
			LastName  string `json:"lastName"`
			FirstName string `json:"firstName"`
		} `json:"names"`
	} `json:"stakeholders"`
	Products []Product `json:"products"`
}

type TripPoint struct {
	DateTime     Time    `json:"dateTime"`
	LocationName string  `json:"locationName"`
	LocationCode string  `json:"locationCode"`
	Address      Address `json:"Address"`
}

type Address struct {
	CountryCode string   `json:"countryCode"`
	CountryName string   `json:"countryName"`
	CityName    string   `json:"cityName"`
	Lines       []string `json:"lines"`
	Text        string   `json:"text"`
}

type Product struct {
	Air   *AirProduct   `json:"air"`
	Hotel *HotelProduct `json:"hotel"`
}

type BaseProduct struct {
	BookingChannel struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"bkgChannel"`
	Description string `json:"description"`
	Status      string `json:"status"`
	NIP         int    `json:"NIP"`
	ConfirmNbr  string `json:"confirmNbr"`
}

type AirProduct struct {
	BaseProduct
	ServiceProvider struct {
		Code              string `json:"code"`
		Name              string `json:"name"`
		BaggagePolicyLink string `json:"baggagePolicyLink"`
	} `json:"serviceProvider"`
	Start struct {
		TripPoint
		Terminal    string `json:"terminal"`
		AirportName string `json:"airportName"`
		CityCode    string `json:"cityCode"`
		CountryCode string `json:"countryCode"`
		RegionCode  string `json:"regionCode"`
	} `json:"start"`
	End struct {
		TripPoint
		Terminal    string `json:"terminal"`
		AirportName string `json:"airportName"`
		CityCode    string `json:"cityCode"`
		CountryCode string `json:"countryCode"`
		VisaAlert   bool   `json:"visaAlert"`
	} `json:"end"`
	Duration string `json:"duration"`
}

type HotelProduct struct {
	BaseProduct
	AccomodationType   string `json:"accomodationType"`
	AdditionalServices string `json:"additionalServices"`
	CancelPolicies     string `json:"cancelPolicies"`
	ServiceProvider    struct {
		Name string
	} `json:"serviceProvider"`
	Start struct {
		DateTime     Time    `json:"dateTime"`
		LocationCode string  `json:"locationCode"`
		Address      Address `json:"Address"`
		Contact      struct {
			Phone string `json:"phone"`
		} `json:"contact"`
	} `json:"start"`
	End struct {
		DateTime Time `json:"dateTime"`
	} `json:"end"`
	CheckInEndTime  Time `json:"checkInEndTime"`
	CheckOutEndTime Time `json:"checkOutEndTime"`
	Rate            struct {
		Description string `json:"description"`
		Code        string `json:"code"`
		Inclusions  string `json:"inclusions"`
	} `json:"rate"`
}
