package federation

import (
	"fmt"
	"log"
	"net/http"

	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/stellar/federation/db"
	"github.com/stellar/federation/db/drivers/mysql"
	"github.com/stellar/federation/db/drivers/postgres"
	"github.com/stellar/federation/db/drivers/sqlite3"
)

type Database interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

type App struct {
	config Config
	driver db.Driver
}

// NewApp constructs an new App instance from the provided config.
func NewApp(config Config) (*App, error) {
	log.Println("Starting app")

	var driver db.Driver
	switch config.Database.Type {
	case "mysql":
		driver = &mysql.MysqlDriver{}
	case "postgres":
		driver = &postgres.PostgresDriver{}
	case "sqlite3":
		driver = &sqlite3.Sqlite3Driver{}
	default:
		return nil, fmt.Errorf("%s database has no driver.", config.Database.Type)
	}

	err := driver.Init(config.Database.Url)
	if err != nil {
		return nil, err
	}

	result := &App{config: config, driver: driver}
	return result, nil
}

func (a *App) Serve() {
	requestHandler := &RequestHandler{config: &a.config, driver: a.driver}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET"},
	})

	handler := alice.New(
		c.Handler,
		headersMiddleware(),
	).Then(http.HandlerFunc(requestHandler.Main))

	http.Handle("/federation", handler)
	http.Handle("/federation/", handler)

	portString := fmt.Sprintf(":%d", a.config.Port)

	log.Println("Starting server on port: ", a.config.Port)

	var err error
	if a.config.Tls.CertificateFile != "" && a.config.Tls.PrivateKeyFile != "" {
		err = http.ListenAndServeTLS(
			portString,
			a.config.Tls.CertificateFile,
			a.config.Tls.PrivateKeyFile,
			nil,
		)
	} else {
		err = http.ListenAndServe(portString, nil)
	}

	if err != nil {
		log.Fatal(err)
	}
}
