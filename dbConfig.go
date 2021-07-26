package main

import (
	"os"
)

type dbConfig struct {
	host     string
	user     string
	password string
	port     string
	dbname   string
}

var myDbConfig = dbConfig{
	os.Getenv("DBHOST"),
	"postgres",
	os.Getenv("DBPASS"),
	"5432",
	"lego",
}
