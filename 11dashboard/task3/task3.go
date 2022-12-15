package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", "151.248.113.144:8030", "http service address")

const url = "https://bmstu.ru/"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func check(url string) string {
	fmt.Println("Проверяем адрес ", url)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Ошибка соединения. %s\n", err)
		return "BAUMAN is not available"
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("Ошибка. http-статус: %s\n", resp.StatusCode)
		return "BAUMAN is not available"
	}
	fmt.Printf("Онлайн. http-статус: %d\n", resp.StatusCode)
	return "BAUMAN is norm"
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/task3", task3)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func task3(w http.ResponseWriter, r *http.Request) {
	c, _ := upgrader.Upgrade(w, r, nil)
	for {
		msg := check(url)
		c.WriteMessage(websocket.TextMessage, []byte(msg))
		time.Sleep(time.Second * 5)
	}

}
