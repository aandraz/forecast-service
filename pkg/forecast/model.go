package forecast

import "time"

type Location struct {
	Id             string  `json:"id"`
	Latitude       float32 `json:"latitude"`
	Longitude      float32 `json:"longitude"`
	Declination    float32 `json:"declination"`
	Azimuth        float32 `json:"azimuth"`
	MaxPeakPowerKw float32 `json:"max_peak_power_kw"`
}

type Info struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance"`
	Place     string  `json:"place"`
	Timezone  string  `json:"timezone"`
	Time      string  `json:"time"`
	TimeUTC   string  `json:"time_utc"`
}

type Message struct {
	Code      int    `json:"code"`
	Type      string `json:"type"`
	Text      string `json:"text"`
	PID       string `json:"pid"`
	Info      Info   `json:"info"`
	RateLimit struct {
		Zone      string `json:"zone"`
		Period    int    `json:"period"`
		Limit     int    `json:"limit"`
		Remaining int    `json:"remaining"`
	} `json:"ratelimit"`
}
type Result struct {
	Watts           map[string]float32 `json:"watts"`
	WattHoursPeriod map[string]float32 `json:"watt_hours_period"`
	WattHours       map[string]float32 `json:"watt_hours"`
	WattHoursDay    map[string]float32 `json:"watt_hours_day"`
}

type Response struct {
	Result  Result  `json:"result"`
	Message Message `json:"message"`
}

type Forecasts struct {
	Ts         time.Time `json:"ts"`
	LocationId string    `json:"location_id"`
	Type       string    `json:"Type"`
	Value      float32   `json:"value"`
	UpdatedAt  time.Time `json:"updated_at"`
}
