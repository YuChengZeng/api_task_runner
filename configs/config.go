package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI           string
	MongoDBName        string
	MongoColl          string
	RateLimit          int
	MaxRetries         int
	IntelligenceHost   string
	IntelligenceKeyCIB string

	LoggerConfig LoggerConfig
}

type LoggerConfig struct {
	LogDir      string
	LogFileName string
	LogLevel    string
	LogKeepDays int
	LogToFile   bool
}

func LoadConfig() *Config {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}


	if err := godotenv.Load(envFile); err != nil {
		fmt.Printf("Warning: Failed to load env file %s\n", envFile)
	}

	_ = godotenv.Overload(".env.example")

	return &Config{
		LoggerConfig: LoggerConfig{
			LogDir:      "./logs",
			LogFileName: "api_runner.log",
			LogLevel:    "debug",
			LogKeepDays: 30,
			LogToFile:   true,
		},
		MongoURI:           os.Getenv("MONGO_URI"),
		MongoDBName:        os.Getenv("MONGO_DB_NAME"),
		MongoColl:          os.Getenv("MONGO_COLL"),
		RateLimit:          1,
		MaxRetries:         10,
		IntelligenceHost:   os.Getenv("INTELLIGENCE_HOST"),
		IntelligenceKeyCIB: os.Getenv("INTELLIGENCE_KEY_CIB"),
	}
}
