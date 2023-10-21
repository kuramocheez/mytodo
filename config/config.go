package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type DatabaseConfig struct{
	ServerPort int
	DBPort int
	DBHost string
	DBUser string
	DBPassword string
	DBName string
}

func InitConfig() *DatabaseConfig{
	res := new(DatabaseConfig)
	res = loadConfig()
	if res == nil{
		logrus.Fatal("Config: Tidak dapat terkoneksi ke database")
		return nil
	}
	return res
}

func loadConfig() *DatabaseConfig{
	res := new(DatabaseConfig)

	err := godotenv.Load(".env")
	if err != nil{
		logrus.Fatal("Config: Tidak Dapat Meload File Config")
	}

	if val, found := os.LookupEnv("SERVER"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Fatal("Config: Port Server Tidak Valid")
		}
		res.ServerPort = port
	}

	if val, found := os.LookupEnv("DBHOST"); found {
		res.DBHost = val
	}

	if val, found := os.LookupEnv("DBUSER"); found {
		res.DBUser = val
	}

	if val, found := os.LookupEnv("DBPASS"); found {
		res.DBPassword = val
	}

	if val, found := os.LookupEnv("DBNAME"); found {
		res.DBName = val
	}
	return res
}