package Weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Weather struct {
	Name       string  `json:"name"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Country    string  `json:"country"`

	Current struct {
		WeatherCode    int     `json:"weather_code"`
		Temperature2M  float64 `json:"temperature_2m"`
		Feeling        float64 `json:"apparent_temperature"`
		WindSpeed10M   float64 `json:"wind_speed_10m"`
		Precipitation  float64 `json:"precipitation"`
		Probability    int     `json:"precipitation_probability"`
	} `json:"current"`	

	Hourly    struct {
		WeatherCode    []int     `json:"weather_code"`
		Temperature2M  []float64 `json:"temperature_2m"`
		Feeling        []float64 `json:"apparent_temperature"`
		WindSpeed10M   []float64 `json:"wind_speed_10m"`
		Precipitation  []float64 `json:"precipitation"`
		Probability    []int     `json:"precipitation_probability"`
	} `json:"hourly"`
}

func (w *Weather) Weather(conf *Config, cityIdx *int) {
	w.getCityInfo(&conf.CitiesList[*cityIdx])
	w.getCurrentWeather()
	w.getFullWeather(conf)

	w.PrintWeather()
}

func (w *Weather) getCityInfo(city *string) {
	apiURL := "https://api.api-ninjas.com/v1/city"
	params := url.Values{"name": {*city}}

	req, err := http.NewRequest("GET", apiURL + "?" + params.Encode(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("X-Api-Key", getToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var temp []Weather
	err = json.Unmarshal(body, &temp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	*w = temp[0]
}

func getToken() string {
	file, err := os.ReadFile("token_ninjas.txt")
	if err != nil {
		fmt.Println("Error: File token_ninjas.txt is not open!")
		os.Exit(1)
	}
	return string(file)
}

func (w *Weather) getCurrentWeather() {
	apiURL := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{"latitude": {FtoS(w.Latitude)},
					    "longitude": {FtoS(w.Longitude)},
					    "timezone": {"auto"},
					   	"current": {"weather_code", 
									"temperature_2m",
					   				"apparent_temperature",
					   				"wind_speed_10m",
					   				"precipitation",
					   				"precipitation_probability"}}

	req, err := http.NewRequest("GET", apiURL + "?" + params.Encode(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	err = json.Unmarshal(body, w)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
}

func (w *Weather) getFullWeather(conf *Config) {
	apiURL := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{"latitude": {FtoS(w.Latitude)},
					    "longitude": {FtoS(w.Longitude)},
					    "timezone": {"auto"},
					   	"hourly": {"weather_code", 
								   "temperature_2m",
					   			   "apparent_temperature",
					   			   "wind_speed_10m",
					   			   "precipitation",
					   			   "precipitation_probability"},
						"forecast_days": {ItoS(conf.CountDays)}}

	req, err := http.NewRequest("GET", apiURL + "?" + params.Encode(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	err = json.Unmarshal(body, w)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
}

func centerString(str string, width int) string {
	if len(str) >= width {
		return str
	}
	padding := width - len(str)
	leftPadding := padding / 2
	rightPadding := padding - leftPadding
	return strings.Repeat(" ", leftPadding) + str + strings.Repeat(" ", rightPadding)
}

func getWeatherCode(code int) string {
	switch code {
	case 0:
		return "Clear"
	case 1:
		return "Mainly Clear"
	case 2:
		return "Partly Cloudy"
	case 3:
		return "Overcast"
	case 45:
		return "Fog"
	case 48:
		return "Depositing Rime Fog"
	case 51:
		return "Light Drizzle"
	case 53:
		return "Moderate Drizzle"
	case 55:
		return "Dense Drizzle"
	case 56:
		return "Light Freezing Drizzle"
	case 57:
		return "Dense Freezing Drizzle"
	case 61:
		return "Slight Rain"
	case 63:
		return "Moderate Rain"
	case 65:
		return "Heavy Rain"
	case 66:
		return "Light Freezing Rain"
	case 67:
		return "Heavy Freezing Rain"
	case 71:
		return "Slight Snow"
	case 73:
		return "Moderate Snow"
	case 75:
		return "Heavy Snow"
	case 77:
		return "Snow Grains"
	case 80:
		return "Slight Rain Showers"
	case 81:
		return "Moderate Rain Showers"
	case 82:
		return "Violent Rain Showers"
	case 85:
		return "Slight Snow Showers"
	case 86:
		return "Heavy Snow Showers"
	case 95:
		return "Thunderstorm"
	case 96:
		return "Thunderstorm with Slight Hail"
	case 99:
		return "Thunderstorm with Heavy Hail"
	default:
		return "Unknown"
	}
}

func FtoS(number float64) string {
	return strconv.FormatFloat(number, 'f', 1, 64)
}

func ItoS(number int) string {
	return strconv.Itoa(number)
}

func (w *Weather) PrintWeather() {
	fmt.Printf("City: %s, %s\n", w.Name, w.Country)
	w.PrintCurrentWeather()
}

func (w *Weather) PrintCurrentWeather() {
    fmt.Printf("┌────────────────────┐\n")
	fmt.Printf("│%s├───────────────────┐\n", centerString("Now", 20))
	fmt.Printf("├────────────────────┘                   │\n")
	fmt.Printf("│%s│\n", centerString(getWeatherCode(w.Current.WeatherCode), 40))
	fmt.Printf("│%s│\n", centerString(FtoS(w.Current.Temperature2M) + "(" + FtoS(w.Current.Feeling) + ") °C", 41))
	fmt.Printf("│%s│\n", centerString(FtoS(w.Current.WindSpeed10M) + " km/h", 40))
	fmt.Printf("│%s│\n", centerString(FtoS(w.Current.Precipitation) + " mm | " + ItoS(w.Current.Probability) + " %", 40))
	fmt.Printf("└────────────────────────────────────────┘\n")
}