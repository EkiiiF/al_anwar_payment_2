package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppName           string
	AppPort           string
	AppEnv            string
	AppURL            string
	AppExternalURL    string
	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPass            string
	CORSOrigins       string
	MidtransServerKey string
	MidtransClientKey string
	MidtransIsProd    bool
	JWTSecret         string
	ServerPort        string
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("CORS_ORIGINS", "http://localhost:5173,http://localhost:3000,http://127.0.0.1:5173")
	viper.SetDefault("APP_PORT", "3000")

	if err := viper.ReadInConfig(); err != nil {
		panic("cannot read .env file: " + err.Error())
	}
	return &Config{
		AppName:           viper.GetString("APP_NAME"),
		AppPort:           viper.GetString("APP_PORT"),
		AppEnv:            viper.GetString("APP_ENV"),
		AppURL:            viper.GetString("APP_URL"),
		AppExternalURL:    viper.GetString("APP_EXTERNAL_URL"),
		DBHost:            viper.GetString("DB_HOST"),
		DBPort:            viper.GetString("DB_PORT"),
		DBName:            viper.GetString("DB_NAME"),
		DBUser:            viper.GetString("DB_USER"),
		DBPass:            viper.GetString("DB_PASS"),
		CORSOrigins:       viper.GetString("CORS_ORIGINS"),
		MidtransServerKey: viper.GetString("MIDTRANS_SERVER_KEY"),
		MidtransClientKey: viper.GetString("MIDTRANS_CLIENT_KEY"),
		MidtransIsProd:    viper.GetBool("MIDTRANS_IS_PRODUCTION"),
		ServerPort:        viper.GetString("SERVER_PORT"),
		JWTSecret:         viper.GetString("JWT_SECRET"),
	}
}
