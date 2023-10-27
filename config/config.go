package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Database Config
type DatabaseConfig struct {
	ServerPort int
	DBPort     int
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
}

// Initial Config untuk Load Config diawal
func InitConfig() *DatabaseConfig {
	var res = new(DatabaseConfig)
	res = loadConfig()
	if res == nil {
		logrus.Fatal("Config: Tidak Dapat Terkoneksi Ke Database")
		return nil
	}
	return res
}

// Load Config dari Env
func loadConfig() *DatabaseConfig {
	var res = new(DatabaseConfig)

	// Load Env
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatal("Config: Tidak Dapat Meload File Config")
	}

	// Get Server Value
	if val, found := os.LookupEnv("SERVER"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Fatal("Config: Port Server Tidak Valid")
		}
		res.ServerPort = port
	}

	// Get DB Port
	if val, found := os.LookupEnv("DBPORT"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Fatal("Config: Port Database Tidak Valid")
		}
		res.DBPort = port
	}

	// Get DB Host Value
	if val, found := os.LookupEnv("DBHOST"); found {
		res.DBHost = val
	}

	// Get DB User Value
	if val, found := os.LookupEnv("DBUSER"); found {
		res.DBUser = val
	}

	// Get DB User Password Value
	if val, found := os.LookupEnv("DBPASS"); found {
		res.DBPassword = val
	}

	// Get DB Name Value
	if val, found := os.LookupEnv("DBNAME"); found {
		res.DBName = val
	}
	return res
}
