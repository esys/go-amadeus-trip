# Trip Parser

This project illustrates the use of [Amadeus Trip API](https://developers.amadeus.com/self-service/category/trip/api-doc/trip-parser) to parse emails and extract travel information from it.

  * [Requirements](#requirements)
  * [Configuration](#configuration)
  * [Running](#running)
  * [Code](#code)
  * [Dependencies](#dependencies)

## How does it work

The go process retrieve mails by polling a Gmail inbox and then process them using the [Amadeus TRIP API](https://developers.amadeus.com/self-service/category/trip/api-doc/trip-parser), 
first by creating a parsing job, then by querying the status to eventually get the results.

Extracted travel information is stored into a SQLite3 database and made available through a REST API.

![alt text](doc/flowchart.svg?raw=true)
  
## Requirements

You will need:
- to [create an Amadeus Self-Service account](https://developers.amadeus.com/register) and create an API key, please see the [FAQ](https://developers.amadeus.com/support/faq/)
- a GMail API access, follow through the [Go Quickstart](https://developers.google.com/gmail/api/quickstart/go), you should get back a credentials and token JSON files

## Configuration

Configure the application through a file
- this file must be called `config.(yaml|json)` and put at the project root
- a sample `config.sample.yaml` with dummy values is available, copy and rename it to `config.yaml` before editing

Or use environment variables (they take precedence over the configuration file)

|Name               |Description                        |Example                        |
|---                |---                                |---                            |
|API_LISTEN         |server:port on which API listens   |:1323                          |
|PARSER_KEY         |Amadeus API key                    |yRveyxreiof83ID2FlldsfgIW95    |
|PARSER_SECRET      |Amadeus API secret                 |d5Gtof7Q4pxlI8KGH              |
|PARSER_URL         |Amadeus API endpoint               |https://test.api.amadeus.com   |
|MAIL_CREDENTIALS   |GMail client credentials JSON file |client_credentials.json        |
|MAIL_TOKEN         |GMail token JSON file              |gmail_token.json               |
|STORAGE_NAME       |SQLite database name               |:memory:                       |

## Running

Build and run it directly
```
$ make run
```

Or Build and run it as a docker image
```
$ make docker-build docker-run
```

Check logs in stdout to see if processing is going OK. 
```
2:02PM DBG trip 95ed6a4c-3910-4bce-8f06-0d2b2ea1d344 (ref: UCFRMZ) written in repository
```

To check for results, query the API
```
$ curl "http://localhost:1323"
[{"ID":"95ed6a4c-3910-4bce-8f06-0d2b2ea1d344","Reference":"UCFRMZ" .... }]

$ curl "http://localhost:1323/trip?ref=UCFRMZ"
{"ID":"95ed6a4c-3910-4bce-8f06-0d2b2ea1d344","Reference":"UCFRMZ","Start":"0001-01-01T00:00:00Z","End":"0001-01-01T00:00:00Z","TripSteps":[{"ID":"1ef923bb-10af-44b7-8022-65c7aae805b3","TripID":"95ed6a4c-3910-4bce-8f06-0d2b2ea1d344","Type":"flight-end","DateTime":"2020-04-12T15:30:00Z","Location":"PARIS","Description":"Flight end with TRANSAVIA FRANCE"},{"ID":"49c85155-4ad5-493e-973b-ec4a939b7a18","TripID":"95ed6a4c-3910-4bce-8f06-0d2b2ea1d344","Type":"flight-start","DateTime":"2020-04-12T11:55:00Z","Location":"TUNIS","Description":"Flight start with TRANSAVIA FRANCE"},{"ID":"df9b3b22-e797-467c-9e71-7ee6d13bb787","TripID":"95ed6a4c-3910-4bce-8f06-0d2b2ea1d344","Type":"flight-end","DateTime":"2020-04-06T17:45:00Z","Location":"TUNIS","Description":"Flight end with TRANSAVIA FRANCE"},{"ID":"ef983fa6-1222-4e2f-9315-9355895570a3","TripID":"95ed6a4c-3910-4bce-8f06-0d2b2ea1d344","Type":"flight-start","DateTime":"2020-04-06T16:10:00Z","Location":"PARIS","Description":"Flight start with TRANSAVIA FRANCE"}]}
```

## Code

Code organization follows the [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) guidelines.
- `domain` contains only models and interfaces
- `usecase` contains application logic
- `adapter` is for interface implementations and communication with external world

```
.
├── adapter
│   ├── api
│   │   ├── rest.go
│   │   └── rest_test.go
│   ├── backend
│   │   ├── mail
│   │   │   └── gmail
│   │   │       ├── gmail.go
│   │   │       └── gmail_test.go
│   │   └── parser
│   │       └── amadeus
│   │           ├── amadeus.go
│   │           ├── amadeus_test.go
│   │           ├── converter.go
│   │           ├── converter_test.go
│   │           ├── dto.go
│   │           ├── dto_test.go
│   │           └── testdata
│   │               ├── air.json
│   │               ├── hotel.json
│   │               └── msg-encoded
│   └── repository
│       ├── sqlite.go
│       └── sqlite_test.go
├── domain
│   ├── backend.go
│   ├── mocks
│   │   ├── TripFinder.go
│   │   └── TripRepository.go
│   ├── model
│   │   ├── email.go
│   │   └── trip.go
│   ├── repository.go
│   └── usecase.go
└── usecase
    ├── processor.go
    ├── tripfinder.go
    └── tripfinder_test.go

```


## Dependencies

This project makes use of:
- [Echo](https://github.com/labstack/echo) : minimalist Go web framework
- [Gorm](https://github.com/go-gorm/gorm) : ORM library
- [Zerolog](https://github.com/rs/zerolog) : JSON Logger
- [Viper](https://github.com/spf13/viper) : configuration
- [Testify](https://github.com/stretchr/testify) : test assertion and mocks
- [Mockery](https://github.com/vektra/mockery): mock object generator
