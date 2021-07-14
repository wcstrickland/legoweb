package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq" // the underscore allows for import without explicit reference
)

var db *sql.DB
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	db = connect(myDbConfig)
	defer db.Close()

	r := httprouter.New()

	r.GET("/", GetHome)
	r.GET("/register", GetRegister)
	r.POST("/register", PostRegister)
	r.GET("/login", GetLogin)
	r.POST("/login", PostLogin)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Listening on http://127.0.0.1:8080 ")
	http.ListenAndServe(":8080", r)
}

func cleanSQL(s string) string {
	prohib := []string{"(", "{", ";", "&", "@", "^", "%", ",", ":", "}", ")"}
	newString := []string{}
	for _, v := range s {
		if strings.Contains(strings.Join(prohib, ""), string(v)) {
			nv := strings.ReplaceAll(string(v), string(v), "")
			newString = append(newString, nv)
		} else {
			newString = append(newString, string(v))
		}
	}
	return strings.Join(newString, "")
}

//TODO refactor to range over values instead of hardcoding based on the number of items

func PostRegister(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	vals := []interface{}{}
	req.ParseForm()
	uid := "jassss"
	vals = append(vals, uid)
	uname := req.Form["uname"][0]
	uname = cleanSQL(uname)
	vals = append(vals, uname)
	item1 := req.Form["item1"][0]
	item1 = cleanSQL(item1)
	vals = append(vals, item1)
	item2 := req.Form["item2"][0]
	item2 = cleanSQL(item2)
	vals = append(vals, item2)
	item3 := req.Form["item3"][0]
	item3 = cleanSQL(item3)
	vals = append(vals, item3)
	query := "insert into users (uid, uname, item1, item2, item3) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.Exec(query, vals...)
	if err != nil {
		fmt.Println("error performing query:", err)
	}
	err = tpl.ExecuteTemplate(res, "success.gohtml", uname)
	if err != nil {
		fmt.Fprintln(res, err)
	}
}

func PostLogin(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	items := []Item{}
	req.ParseForm()
	uname := req.Form["uname"][0]
	uname = cleanSQL(uname)
	query := fmt.Sprintf("select * from %s", uname)
	rows, err := db.Query(query)
	if err != nil {
		err = tpl.ExecuteTemplate(res, "error.gohtml", nil)
		if err != nil {
			fmt.Println("error showing error page. how ironic:", err)
		}
		return
		fmt.Println("error performing query:", err)
	}
	defer rows.Close()
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.Item, &item.Status, &item.Check_time)
		if err != nil {
			fmt.Println("error getting row:", err)
		}
		items = append(items, item)
	}
	report := Report{
		Item1: items[0],
		Item2: items[1],
		Item3: items[2],
	}
	err = tpl.ExecuteTemplate(res, "show.gohtml", report)
	if err != nil {
		err = tpl.ExecuteTemplate(res, "error.gohtml", nil)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}

func connect(dc dbConfig) *sql.DB {
	configString := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable",
		dc.host, dc.user, dc.password, dc.port, dc.dbname)
	db, err := sql.Open("postgres", configString)
	if err != nil {
		log.Fatalln("something went wrong connecting to the database : ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln("Error pinging db: ", err)
	}
	return db
}

func GetHome(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	err := tpl.ExecuteTemplate(res, "index.gohtml", nil)
	if err != nil {
		fmt.Fprintln(res, "something went wrong")
	}
}

func GetRegister(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	err := tpl.ExecuteTemplate(res, "register.gohtml", nil)
	if err != nil {
		fmt.Fprintln(res, "something went wrong")
	}
}

func GetLogin(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	err := tpl.ExecuteTemplate(res, "login.gohtml", nil)
	if err != nil {
		fmt.Fprintln(res, "something went wrong")
	}
}
