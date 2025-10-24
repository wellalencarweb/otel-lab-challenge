package config

import "github.com/spf13/viper"

type Conf struct {
	LogLevel                         string `mapstructure:"LOG_LEVEL"`
	InputServiceWebServerPort        int    `mapstructure:"INPUT_SERVICE_WEB_SERVER_PORT"`
	OrchestratorServiceWebServerPort int    `mapstructure:"ORCHESTRATOR_SERVICE_WEB_SERVER_PORT"`
	HttpClientTimeout                int    `mapstructure:"HTTP_CLIENT_TIMEOUT_MS"`
	ViaCepApiBaseUrl                 string `mapstructure:"VIACEP_API_BASE_URL"`
	WeatherApiBaseUrl                string `mapstructure:"WEATHER_API_BASE_URL"`
	WeatherApiKey                    string `mapstructure:"WEATHER_API_KEY"`
	OrchestratorServiceHost          string `mapstructure:"ORCHESTRATOR_SERVICE_HOST"`
	OtelCollectorURL                 string `mapstructure:"OTEL_COLLECTOR_URL"`
}

func LoadConfig(path string) (*Conf, error) {
	var c *Conf

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	return c, nil
}
