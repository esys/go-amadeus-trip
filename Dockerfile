FROM golang:1.14

WORKDIR /app
COPY . .
RUN go mod tidy
RUN make build
CMD ["./bin/amadeus-trip-parser"]