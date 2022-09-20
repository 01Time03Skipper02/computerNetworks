package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mmcdole/gofeed"
)

const (
	password string = "Je2dTYr6"
	login    string = "iu9networkslabs"
	host     string = "students.yss.su"
	dbname   string = "iu9networkslabs"
)

func main() {
	db, err := sql.Open("mysql", login+":"+password+"@tcp("+host+")/"+dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://briansk.ru/rss20_briansk.xml")
	for _, item := range feed.Items {
		_, err = db.Exec("insert ignore into iu9networkslabs.iu9Shvets (title, link, description, category, pubDate) values (?, ?, ?, ?, ?)",
			item.Title, item.Link, item.Description, item.Categories[0], item.PublishedParsed)
		if err != nil {
			panic(err)
		}
	}
}
