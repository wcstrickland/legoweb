package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // the underscore allows for import without explicit reference
	"html/template"
	//    "net/http"
	//	"log"
	//	"os"
	//	"strings"
	//	"time"
)

var db *sql.DB
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	configString := fmt.Sprintf("host=%s user=% password=%s port=%s dbname%s sslmode=disable",
		dbConfig.host, dbConfig.user, dbConfig.password, dbConfig.port, dbConfig.dbname)
	db, err := sql.Open("postgres")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging db: ", err)
	}

	fmt.Println("helloworld")
}
