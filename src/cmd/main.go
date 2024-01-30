package main

import (
	"my-little-olap/internal/db"
	"my-little-olap/internal/srv"
	"my-little-olap/internal/utils"
	"os"
)

func main() {
	logger := utils.NewLogger()
	ds, err := db.NewClickhouseDB(newDBConfig(&logger))
	if err != nil {
		panic(err)
	}
	s := srv.NewOLAPServer(ds, &logger)
	err = s.Run()
	if err != nil {
		logger.Error.Printf("Server running error: %s\n", err)
	}
}

func newDBConfig(logger *utils.Logger) db.Config {
	var password string
	pFile, exist := os.LookupEnv("MY_LITTLE_OLAP_DB_PASSWORD_FILE")
	if exist {
		data, err := os.ReadFile(pFile)
		if err == nil {
			password = string(data)
		}
	}
	return db.Config{
		Host:     os.Getenv("MY_LITTLE_OLAP_DB_HOST"),
		DBName:   os.Getenv("MY_LITTLE_OLAP_DB_DBNAME"),
		User:     os.Getenv("MY_LITTLE_OLAP_DB_USER"),
		Password: password,
		Logger:   logger,
	}
}
