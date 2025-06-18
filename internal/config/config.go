package config

import ()

type Config struct {
	URL  string `json:"db_url"`
	User string `json:"username"`
}

func (config *Config) SetUser() {
	config.User = "Reklai"
}
