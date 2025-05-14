package configs

import (
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
	IntelligenceKeyBCS string

	Logger LoggerConfig
}

type LoggerConfig struct {
	LogDir      string
	LogFileName string
	LogLevel    string
	LogKeepDays int
	LogToFile   bool
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		MongoURI:           os.Getenv("MONGO_URI"),
		MongoDBName:        os.Getenv("MONGO_DB_NAME"),
		MongoColl:          os.Getenv("MONGO_COLL"),
		RateLimit:          2,
		MaxRetries:         3,
		IntelligenceHost:   os.Getenv("INTELLIGENCE_HOST"),
		IntelligenceKeyCIB: os.Getenv("INTELLIGENCE_KEY_CIB"),
		IntelligenceKeyBCS: os.Getenv("INTELLIGENCE_KEY_BCS"),
		Logger: LoggerConfig{
			LogDir:      "./logs",
			LogFileName: "api_runner.log",
			LogLevel:    "info",
			LogKeepDays: 30,
			LogToFile:   true,
		},
	}
}
