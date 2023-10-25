package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/justty/golang-units/app/internal/config"
	rest_v1 "github.com/justty/golang-units/app/internal/handler/rest/v1"
	"github.com/justty/golang-units/app/internal/service/apod"
	"github.com/justty/golang-units/app/internal/store"
	"github.com/justty/golang-units/app/pkg/cron"
	"github.com/justty/golang-units/app/pkg/db/postgresql"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s", err)
	}
}

func run() error {
	log.Println("Initializing config")
	config := config.GetConfig()

	log.Println("Initializing database connection")
	configDB := config.Database
	dbConnection, err := postgresql.GetDBConnection(configDB.User, configDB.Password, configDB.Host, configDB.Port, configDB.DBname)
	if err != nil {
		return fmt.Errorf("error connect to the database: %w", err)
	}
	defer dbConnection.Close()

	apodAPI := apod.NewApodAPI(config.ApodApiKey)
	apodMetadataStore := store.NewApodMetadataStore(dbConnection)
	apodWorker := apod.NewApodWorker(apodMetadataStore, apodAPI)
	apodHandler := rest_v1.NewApodHandler(apodMetadataStore)

	if config.CronEnabled {
		cronManager := cron.NewCron()

		cronManager.AddJob(cron.NewJob(apod.WorkerServiceName, 24*time.Hour, apodWorker.Process))

		cronManager.Start()
	}

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	apiRouter.HandleFunc("/apod_metadata", apodHandler.GetApodMetadata).Methods("GET")
	http.Handle("/", router)

	log.Println("Started server on port 80")
	log.Println(http.ListenAndServe(":8080", router))

	return nil
}
