package forecast

import (
	"encoding/json"
	"forecast/config"
	"github.com/go-pg/pg/v10"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	GetForecast()
	RegisterLocations(id string, latitude, longitude, declination, azimut, maxPeakPowerKw float32)
}

type service struct {
	pgdb      *pg.DB
	conf      *config.Configuration
	locations map[string]Location // key location id
}

func NewService(pgdb *pg.DB, config *config.Configuration) Service {
	return &service{
		pgdb:      pgdb,
		conf:      config,
		locations: make(map[string]Location),
	}
}

func (s *service) RegisterLocations(id string, latitude, longitude, declination, azimut, maxPeakPowerKw float32) {
	s.locations[id] = Location{
		Id:             id,
		Latitude:       latitude,
		Longitude:      longitude,
		Declination:    declination,
		Azimuth:        azimut,
		MaxPeakPowerKw: maxPeakPowerKw,
	}
	log.Printf("Registered location: %v\n", s.locations[id])
}

func (s *service) GetForecast() {

	for _, location := range s.locations {
		forecasts, err := s.fetchDataFromAPI(buildURL(s.conf.BaseUrl, location.Latitude, location.Longitude, location.Declination, location.Azimuth, location.MaxPeakPowerKw))
		if err != nil {
			log.Printf("error fetching data from api: %v\n", err)
			return
		}
		if s.conf.LogLevel == "debug" {
			log.Printf("Forecast: %v\n", forecasts.Result)
		}

		go s.generate("watts", "1", forecasts.Result.Watts)
		go s.generate("watt_hours_period", "1", forecasts.Result.WattHoursPeriod)
		go s.generate("watt_hours", "1", forecasts.Result.WattHours)
	}
}

func (s *service) generate(fType, locationId string, data map[string]float32) {

	var forecasts []Forecasts

	for t, value := range data {

		tt, err := parseTimestampToUTC(t)
		if err != nil {
			log.Printf("error parsing timestamp: %v, with error: %v\n", t, err)
			continue
		}

		forecasts = append(forecasts, Forecasts{
			Ts:         tt,
			LocationId: locationId,
			Type:       fType,
			Value:      value,
			UpdatedAt:  time.Now(),
		})
	}

	// upsert the forecasts
	_, err := s.pgdb.Model(&forecasts).
		OnConflict("(ts, location_id, type) DO UPDATE").
		Set("value = EXCLUDED.value").
		Insert()
	if err != nil {
		log.Printf("error persisting forecasts: %v\n", err)
	}

	if s.conf.LogLevel == "debug" {
		log.Printf("Forecasts persisted: %v\n", forecasts)
	}
}

func (s *service) fetchDataFromAPI(apiURL string) (*Response, error) {
	log.Printf("Fetching data from API: %v\n", apiURL)
	//Make a GET request to the API
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}

	if s.conf.LogLevel == "debug" {
		log.Printf("Response: %v\n", response)
	}

	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func float32toString(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', 4, 32)
}

func buildURL(baseUrl string, latitude, longitude, declination, azimuth, maxPeakPower float32) string {
	url := baseUrl + float32toString(latitude) + "/" + float32toString(longitude) + "/" + float32toString(declination) + "/" + float32toString(azimuth) + "/" + float32toString(maxPeakPower) + "?time=utc"
	return url
}

func parseTimestampToUTC(timestamp string) (time.Time, error) {
	// Parse the timestamp string
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}, err
	}

	// Convert to UTC
	utcTime := parsedTime.UTC()

	return utcTime, nil
}
