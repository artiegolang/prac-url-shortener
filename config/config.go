package config

import (
	"flag"
	"os"
)

type Options struct {
	FlagRunAddr     string
	BaseURL         string
	FileStoragePath string
}

func ParseFlags() *Options {
	opts := &Options{}

	// Устанавливаем значения по умолчанию
	defaultRunAddr := ":8888"
	defaultBaseURL := "http://localhost:8888"
	defaultFileStoragePath := "D:/practicum-middle/short-url-db.json"

	// Регистрируем флаги командной строки
	flag.StringVar(&opts.FlagRunAddr, "a", defaultRunAddr, "address and port to run server")
	flag.StringVar(&opts.BaseURL, "b", defaultBaseURL, "base URL")
	flag.StringVar(&opts.FileStoragePath, "f", defaultFileStoragePath, "file storage path")

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

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		opts.FileStoragePath = envFileStoragePath
	} else if flag.Lookup("f").Value.String() != defaultFileStoragePath {
		opts.FileStoragePath = flag.Lookup("f").Value.String()
	}

	return opts
}
