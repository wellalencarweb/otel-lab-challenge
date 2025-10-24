package dto

type GetTemperaturesByZipCodeOutput struct {
	City       string  `json:"city"`
	Celcius    float32 `json:"temp_C"`
	Fahrenheit float32 `json:"temp_F"`
	Kelvin     float32 `json:"temp_K"`
}
