package main

import (
	utils "7WebSocket"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8023", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Данные о телах: %s", message)

		bodies := make([]utils.Body, 5)
		err = json.Unmarshal(message, &bodies)
		if err != nil {
			log.Fatalln(err)
		}
		var res utils.Result
		var massSum, xmas, ymas, zmas float64
		for i := 0; i < len(bodies); i++ {
			massSum += bodies[i].M
			xmas += bodies[i].X * bodies[i].M
			ymas += bodies[i].Y * bodies[i].M
			zmas += bodies[i].Z * bodies[i].M
		}
		res.X = xmas / massSum
		res.Y = ymas / massSum
		res.Z = zmas / massSum

		message, err = json.Marshal(res)
		if err != nil {
			log.Fatalln(err)
		}

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
