package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type Configuration struct {
	Env       string
	LogLevel  string
	BaseUrl   string
	JobCron   string
	Locations []Location
	Db        struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
}

type Location struct {
	Id             string
	Latitude       float32
	Longitude      float32
	Declination    float32
	Azimuth        float32
	MaxPeakPowerKw float32
}

func LoadConfig() *Configuration {
	log.SetOutput(os.Stdout)

	configuration := Configuration{}

	viper.AddConfigPath("./config") //Viper looks here for the files.
	viper.SetConfigType("yaml")     //Sets the format of the config file.
	viper.SetConfigName("default")  // So that Viper loads default.yml.
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Warning could not load configuration: %v", err))
	}

	viper.AutomaticEnv() // Merges any overrides set through env vars.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.Unmarshal(&configuration)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Configuration: %v\n: ", configuration)
	return &configuration
}
