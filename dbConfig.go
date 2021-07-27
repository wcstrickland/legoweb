package main

type dbConfig struct {
	host     string
	user     string
	password string
	port     string
	dbname   string
}

var myDbConfig = dbConfig{
	"lego.ceurumcah93i.us-east-1.rds.amazonaws.com",
	"postgres",
	"ThirtyFour12",
	"5432",
	"lego",
}
