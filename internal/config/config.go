package config

import "os"

type Config struct {
	GatewayPort string
	UserPort    string
	HRPort      string
	UserDB      string
	HRDB        string
	UserAddress string
	HRAddress   string
}

func Load() Config {
	c := Config{
		GatewayPort: Env("GATEWAY_PORT", "8000"),
		UserPort:    Env("USER_PORT", "8080"),
		HRPort:      Env("HR_PORT", "9090"),
		UserDB:      Env("USER_DB", "data/users.db"),
		HRDB:        Env("HR_DB", "data/hr.db"),
	}
	c.UserAddress = Env("USER_SERVICE_ADDR", "127.0.0.1:"+c.UserPort)
	c.HRAddress = Env("HR_SERVICE_ADDR", "127.0.0.1:"+c.HRPort)
	return c
}

func Env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
