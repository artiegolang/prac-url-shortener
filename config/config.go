package config

import (
	"flag"
	"os"
)

type Options struct {
	FlagRunAddr string
	BaseURL     string
}

func ParseFlags() *Options {
	opts := &Options{}

	// Устанавливаем значения по умолчанию
	defaultRunAddr := ":8888"
	defaultBaseURL := "http://localhost:8888"

	// Регистрируем флаги командной строки
	flag.StringVar(&opts.FlagRunAddr, "a", defaultRunAddr, "address and port to run server")
	flag.StringVar(&opts.BaseURL, "b", defaultBaseURL, "base URL")

	// Парсим флаги
	flag.Parse()

	// Приоритет: переменные окружения -> флаги -> значения по умолчанию
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		opts.FlagRunAddr = envRunAddr
	} else if flag.Lookup("a").Value.String() != defaultRunAddr {
		// Если переменная окружения не указана, используем флаг, если он был изменён
		opts.FlagRunAddr = flag.Lookup("a").Value.String()
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		opts.BaseURL = envBaseURL
	} else if flag.Lookup("b").Value.String() != defaultBaseURL {
		// Если переменная окружения не указана, используем флаг, если он был изменён
		opts.BaseURL = flag.Lookup("b").Value.String()
	}

	return opts
}
