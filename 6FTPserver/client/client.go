package main

import (
	"bytes"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jlaffaye/ftp"
	"github.com/mmcdole/gofeed"
	"github.com/skorobogatov/input"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

const (
	password string = "Je2dTYr6"
	login    string = "iu9networkslabs"
	host     string = "students.yss.su"
	dbname   string = "iu9networkslabs"
)

func main() {
	c, err := ftp.Dial("localhost:2121", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatalln("cannot connect to server")
	}
	err = c.Login("time_skipper", "12345678990")
	if err != nil {
		log.Fatalln(err)
	}
waitCommand:
	print(">")
	command := input.Gets()
	switch command {
	case "add":
		data := bytes.NewBufferString("Hello World again")
		err = c.Stor(strconv.FormatInt(time.Now().Unix(), 10)+".txt", data)
		if err != nil {
			log.Println(err)
		}
	case "read":
		print("print name of file what you want to read: ")
		fileName := input.Gets()
		r, err := c.Retr(fileName)
		if err != nil {
			log.Println(err)
		}
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			println(err.Error())
		}
		println(string(buf))
		err = r.Close()
		if err != nil {
			println(err.Error())
		}
	case "mkdir":
		print("print name of dir what you want to make: ")
		newDir := input.Gets()
		err = c.MakeDir(newDir)
		if err != nil {
			println(err)
		}
	case "rm":
		print("print name of file what you want to delete: ")
		rm := input.Gets()
		err = c.Delete(rm)
		if err != nil {
			println(err)
		}
	case "dir":
		names, err := c.List("")
		if err != nil {
			println(err)
		}
		for i := 0; i < len(names); i++ {
			println(names[i].Name)
		}
	case "news":
		db, err := sql.Open("mysql", login+":"+password+"@tcp("+host+")/"+dbname)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL("https://briansk.ru/rss20_briansk.xml")
		news := []string{}
		for _, item := range feed.Items {
			news = append(news, item.Title, item.Link, item.Description, item.Categories[0], "\n")
			_, err = db.Exec("insert ignore into iu9networkslabs.iu9Shvets (title, link, description, category, pubDate) values (?, ?, ?, ?, ?)",
				item.Title, item.Link, item.Description, item.Categories[0], item.PublishedParsed)
			if err != nil {
				panic(err)
			}
		}
		for _, str := range news {
			data := bytes.NewBufferString(str)
			err = c.Stor("Shvets_Alexander"+strconv.FormatInt(time.Now().Unix(), 10)+".txt", data)
			if err != nil {
				panic(err)
			}
		}
	case "quit":
		goto end
	}
	goto waitCommand
end:
	if err := c.Quit(); err != nil {
		log.Fatalln(err)
	}
}
