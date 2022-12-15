package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"time"
)

var address = flag.String("addr", "151.248.113.144:8021", "http service address")

func contains(t time.Time, news []feed) bool {
	for _, n := range news {
		if t == n.pubDate {
			return true
		}
	}
	return false
}
func clear() {
	db, _ := sql.Open("mysql", login+":"+password+"@tcp("+host+")/"+dbname+"?parseTime=true")
	db.Query("TRUNCATE iu9networkslabs.iu9Shvets")
}
func save() {
	db, err := sql.Open("mysql", login+":"+password+"@tcp("+host+")/"+dbname+"?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fp := gofeed.NewParser()
	feeds, _ := fp.ParseURL("https://briansk.ru/rss20_briansk.xml")

	rows, error := db.Query("select * from iu9networkslabs.iu9Shvets")
	if err != nil {
		panic(error)
	}
	defer rows.Close()
	var news []feed
	for rows.Next() {
		n := feed{}
		err := rows.Scan(&n.title, &n.link, &n.description, &n.category, &n.pubDate)
		fmt.Println(n.pubDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		news = append(news, n)
	}
	for _, item := range feeds.Items {

		if !contains(item.PublishedParsed.UTC(), news) {
			_, err = db.Exec("insert ignore into iu9networkslabs.iu9Shvets (title, link, description, category, pubDate) values (?, ?, ?, ?, ?)",
				item.Title, item.Link, item.Description, item.Categories[0], item.PublishedParsed)
			if err != nil {
				panic(err)
			}
		}

	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/task2_1", task2_1)
	log.Fatal(http.ListenAndServe(*address, nil))
}

func task2_1(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		msg := string(message)
		switch msg {
		case "saveNews":
			save()
		case "clear":
			clear()
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
