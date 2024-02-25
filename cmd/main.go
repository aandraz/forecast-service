package main

import (
	"fmt"
	"forecast/config"
	"forecast/db"
	"forecast/pkg/forecast"
	"github.com/robfig/cron"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	conf := config.LoadConfig()

	//start the db
	pgdb, err := db.StartDB(conf)
	if err != nil {
		log.Printf("error starting the database %v\n", err)
	}

	// init the forecast service
	forecastService := forecast.NewService(pgdb, conf)

	// register locations set in the configuration
	for _, location := range conf.Locations {
		forecastService.RegisterLocations(location.Id, location.Latitude, location.Longitude, location.Declination, location.Azimuth, location.MaxPeakPowerKw)
	}

	// start the cron job. By default, it will run @hourly
	c := cron.New()
	err = c.AddFunc(conf.JobCron, forecastService.GetForecast)
	if err != nil {
		log.Panicf("error adding cron job: %v\n", err)
	}
	log.Printf("Starting the cron jobs. Trigger jobs: %s\n", conf.JobCron)
	c.Start()

	// Wait for a signal to exit the program gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Received signal to exit")

	c.Stop()
}
