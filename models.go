package main

type Item struct {
	Item       string
	Status     string
	Check_time string
}

type User struct {
	uid   string
	uname string
	item1 string
	item2 string
	item3 string
}

type Report struct {
	Item1 Item
	Item2 Item
	Item3 Item
}
