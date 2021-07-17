package main

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
