package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq" // the underscore allows for import without explicit reference
	"github.com/satori/go.uuid"
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
	http.ListenAndServe(":80"), r)
}

func cleanSQL(s string) string {
	prohib := []string{"(", "{", ";", "&", "@", "^", "%", ",", ":", "}", ")", "'", "\""}
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

func checkUserExists(uname string) error {
	uname = cleanSQL(uname)
	u := User{
		Uid:   "",
		Uname: uname,
		Items: []string{
			"",
			"",
			"",
		},
	}
	err := db.QueryRow("select * from users where uname = $1", uname).Scan(&u.Uid, &u.Uname, &u.Items[0], &u.Items[1], &u.Items[2])
	switch {
	case err != sql.ErrNoRows:
		return errors.New("User already exists")
	case err != nil:
		fmt.Printf("query error: %v\n", err)
	default:
		fmt.Printf("username is %q", uname)
	}
	return nil
}

func PostRegister(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var err error
	req.ParseForm()
	vals := []interface{}{
		uuid.Must(uuid.NewV4(), err),
		cleanSQL(req.Form["uname"][0]),
		cleanSQL(req.Form["item1"][0]),
		cleanSQL(req.Form["item2"][0]),
		cleanSQL(req.Form["item3"][0]),
	}
	err = checkUserExists(cleanSQL(req.Form["uname"][0]))
	fmt.Println("error inside the post route handler:", err)
	if err != nil {
		fmt.Println("error:", err)
		_ = tpl.ExecuteTemplate(res, "error.gohtml", err)
		return
	}
	insertUser := "insert into users (uid, uname, item1, item2, item3) VALUES ($1, $2, $3, $4, $5)"
	createUserTable := fmt.Sprintf("create table if not exists %s(item varchar(255), status varchar(255), check_time varchar(255))", cleanSQL(req.Form["uname"][0]))
	starterReport := fmt.Sprintf("insert into %s (item, status, check_time) values ($1, $2, $3)", cleanSQL(req.Form["uname"][0]))
	_, err = db.Exec(createUserTable)
	if err != nil {
		fmt.Println("error creating user table:", err)
	}
	for i, v := range vals {
		if i > 1 {
			_, err := db.Exec(starterReport, v, "unchecked", "unchecked")
			if err != nil {
				fmt.Println(i, "error creating stock table:", err)
			}
		}
	}
	_, err = db.Exec(insertUser, vals...)
	if err != nil {
		fmt.Println("error performing query:", err)
	}
	err = tpl.ExecuteTemplate(res, "success.gohtml", cleanSQL(req.Form["uname"][0]))
	if err != nil {
		fmt.Println("error rendering template: ", err)
		_ = tpl.ExecuteTemplate(res, "error.gohtml", err)
		return
	}
}

func PostLogin(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	req.ParseForm()
	items := []Item{}
	query := fmt.Sprintf("select * from %s", cleanSQL(req.Form["uname"][0]))
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("error performing query:", err)
		_ = tpl.ExecuteTemplate(res, "error.gohtml", err)
		return
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
	report := Report{items}
	err = tpl.ExecuteTemplate(res, "show.gohtml", report)
	if err != nil {
		_ = tpl.ExecuteTemplate(res, "error.gohtml", err)
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

type Item struct {
	Item       string
	Status     string
	Check_time string
}

type User struct {
	Uid   string
	Uname string
	Items []string
}

type Report struct {
	ReportItems []Item
}
