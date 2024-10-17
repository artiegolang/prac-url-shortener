package config

import (
	"os"
)

type Options struct {
	ServerAddress string
	BaseURL       string
}

func ParseConfig() *Options {
	opt := &Options{
		ServerAddress: ":8085",
		BaseURL:       "http://localhost:8085",
	}

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		opt.ServerAddress = envRunAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		opt.BaseURL = envBaseURL
	}

	return opt
}
