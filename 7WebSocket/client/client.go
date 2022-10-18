package main

import (
	utils "7WebSocket"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8023", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	bodies := make([]utils.Body, 5)
	for i := 1; i < 6; i++ {
		fmt.Printf("Введите координату X для %d объекта: ", i)
		fmt.Scanf("%f", &bodies[i-1].X)
		fmt.Printf("Введите координату Y для %d объекта: ", i)
		fmt.Scanf("%f", &bodies[i-1].Y)
		fmt.Printf("Введите координату Z для %d объекта: ", i)
		fmt.Scanf("%f", &bodies[i-1].Z)
		fmt.Printf("Введите массу %d объекта: ", i)
		fmt.Scanf("%f", &bodies[i-1].M)
	}

	message, err := json.Marshal(bodies)
	if err != nil {
		log.Fatalln(err)
	}

	err = c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("write:", err)
		return
	}

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("Центр масс: %s", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
