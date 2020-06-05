package main

import (
	"amadeus-trip-parser/internal/adapter/api"
	"amadeus-trip-parser/internal/adapter/backend/mail/gmail"
	"amadeus-trip-parser/internal/adapter/backend/parser/amadeus"
	"amadeus-trip-parser/internal/adapter/repository"
	"amadeus-trip-parser/internal/domain"
	"amadeus-trip-parser/internal/usecase"
	"database/sql"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func loadConfig() {
	//load environment variables: e.g. 'parser.url' config key matches 'PARSER_URL' environment variable
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(rootDir())

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic().Msgf("error reading config file: %s", err)
		//TODO set default configuration
	}
	log.Debug().Msgf("config keys: %s", viper.AllKeys())
}

func initMailClient() domain.EmailProvider {
	mc, err := gmail.NewGMailClient(
		viper.GetString("mail.credentials"),
		viper.GetString("mail.token"))
	if err != nil {
		log.Panic().Msgf("when creating mail client: %s", err)
	}
	return mc
}

func initMailParser() domain.EmailParser {
	p, err := amadeus.NewAmadeusTripAPI(
		viper.GetString("parser.url"),
		viper.GetString("parser.key"),
		viper.GetString("parser.secret"))
	if err != nil {
		log.Panic().Msgf("when creating parser: %s", err)
	}
	return p
}

func initRepository() domain.TripRepository {
	name := viper.GetString("repository.name")
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Panic().Msgf("cannot open DB connection with db name %s: %s", name, err)
	}
	repo, err := repository.NewSQLiteTripRepo(db)
	if err != nil {
		log.Panic().Msgf("cannot open repository with db name %s: %s", name, err)
	}
	return repo
}

func runServer(repo domain.TripRepository) {
	finder := usecase.NewTripFinder(repo)
	api := api.NewTripAPI(finder)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	e.GET("/trip", api.Get)

	e.Start(viper.GetString("api.listen"))
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	loadConfig()

	repo := initRepository()
	mail := initMailClient()
	parser := initMailParser()
	proc := usecase.NewEmailProcessor(mail, parser, repo)
	proc.Process()

	runServer(repo)
	proc.Stop()
}
