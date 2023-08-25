package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/justty/golang-units/configs"
	rest_v1 "github.com/justty/golang-units/internal/handler/rest/v1"
	"github.com/justty/golang-units/internal/service/apod"
	"github.com/justty/golang-units/internal/store"

	"github.com/justty/golang-units/pkg/cron"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s", err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error godotenv load: %w", err)
	}

	dbConnection, err := configs.GetDBConnection()
	if err != nil {
		return fmt.Errorf("error connect to the database: %w", err)
	}
	defer dbConnection.Close()

	cronEnabled, err := strconv.ParseBool(os.Getenv("CRON_ENABLED"))
	if err != nil {
		return fmt.Errorf("error parsing CRON_ENABLED from .env file")
	}

	apodAPI := apod.NewApodAPI()
	apodMetadataStore := store.NewApodMetadataStore(dbConnection)
	apodWorker := apod.NewApodWorker(apodMetadataStore, apodAPI)
	apodHandler := rest_v1.NewApodHandler(apodMetadataStore)

	cronStore := cron.NewCronStore(dbConnection)
	cronManager := cron.NewCron(cronStore)

	if cronEnabled {
		cronManager.AddJob(cron.NewJob().Name(apod.WorkerServiceName).Every(1 * time.Minute).Task(func() { apodWorker.Process() }))
	}

	go cronManager.Start()

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/apod_metadata", apodHandler.GetApodMetadata).Methods("GET")
	http.Handle("/", router)

	fmt.Println("Started server on port 80")
	fmt.Println(http.ListenAndServe(":80", router))

	return nil
}
