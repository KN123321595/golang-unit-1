package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/justty/golang-units/configs"
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

	cronStore := cron.NewCronStore(dbConnection)
	cronManager := cron.NewCron(cronStore)

	if cronEnabled {
		cronManager.AddJob(cron.NewJob().Name(apod.WorkerServiceName).At(12 * time.Hour).Task(func() { apodWorker.Process() }))
	}

	go cronManager.Start()

	fmt.Println("Started server on port 80")
	fmt.Println(http.ListenAndServe(":80", nil))

	return nil
}
