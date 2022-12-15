package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	host     = "students.yss.su"
	password = "3Ru7yOTA"
	user     = "ftpiu8"
)

var addr = flag.String("addr", "151.248.113.144:8010", "http service address")
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ls(conn *ftp.ServerConn) []string {
	curDir, err := conn.CurrentDir()
	fmt.Println(curDir)
	if err != nil {
		panic(nil)
	}
	filesList, err := conn.NameList(curDir)
	if err != nil {
		panic(err)
	}
	return filesList
}
func get(conn *ftp.ServerConn, pathToFile string) []byte {
	resp, err := conn.Retr(pathToFile)
	if err != nil {
		panic(err)
	}
	defer func(resp *ftp.Response) {
		err := resp.Close()
		if err != nil {
			panic(err)
		}
	}(resp)
	buf, err := ioutil.ReadAll(resp)
	if err != nil {
		panic(err)
	}
	return buf
}
func task1(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	c, err := ftp.Dial(host + ":ftp")
	if err != nil {
		log.Fatal(err)
	}
	err = c.Login(user, password)
	if err != nil {
		log.Fatal(err)
	}
	f := true
	for {
		files := ls(c)
		find := false
		for _, file := range files {
			if file == "/achtung.txt" {
				find = true
				break
			}
		}
		if find {
			if f {
				f = !f
				message := get(c, "/achtung.txt")
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("write:", err)
					break
				}
			} else {
				continue
			}
		} else {
			if !f {
				f = !f
				err = conn.WriteMessage(websocket.TextMessage, []byte("norm"))
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}
func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/task1", task1)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
