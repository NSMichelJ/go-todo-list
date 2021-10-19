package config

import "os"

type Config struct{
	DB *DBConfig
}

type DBConfig struct{
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	DBName     string
	Charset  string
}

func GetDBConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "mysql",
			Host:     "127.0.0.1",
			Port:     3306,
			Username: os.Getenv("USERNAME"),
			Password: os.Getenv("PASSWORD"),
			DBName:    os.Getenv("DBNAME") ,
			Charset:  "utf8",
		},
	}
}