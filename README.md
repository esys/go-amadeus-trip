# Trip Parser

This project illustrates the use of [Amadeus Trip API](https://developers.amadeus.com/self-service/category/trip/api-doc/trip-parser) to parse emails and extract travel information from it.
The results are then made available through a REST API.

![alt text](https://github.com/esys/go-amadeus-trip/tree/master/doc/flowchart.svg?raw=true)

## Requirements

You will need:
- to [create an AWS Self-Service account](https://developers.amadeus.com/register) and create an API key, please see the [FAQ](https://developers.amadeus.com/support/faq/)
- a GMail API access, follow through the [Go Quickstart](https://developers.google.com/gmail/api/quickstart/go), you should get back a credentials and token JSON files

## Configuration

Configure the application through a file
- this file must be called `config.(yaml|json)` and put at the project root
- a sample `config.sample.yaml` with dummy values is available, copy and rename it to `config.yaml` before editing

Or use environment variables (they take precedence over the configuration file)

|Name               |Description                        |Example                        |
|---                |---                                |---                            |
|api.listen         |server:port on which API listens   |:1323                          |
|parser.key         |Amadeus API key                    |yRveyxreiof83ID2FlldsfgIW95    |
|parser.secret      |Amadeus API secret                 |d5Gtof7Q4pxlI8KGH              |
|parser.url         |Amadeus API endpoint               |https://test.api.amadeus.com   |
|mail.credentials   |GMail client credentials JSON file |client_credentials.json        |
|mail.token         |GMail token JSON file              |gmail_token.json               |
|storage.name       |SQLite database name               |:memory:                       |

## Running

```
go run cmd/parser/main.go
```
