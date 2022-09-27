package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
)

import "proto"

type ClientPool struct {
	Clients []*Client
}

// Client - состояние клиента.
type Client struct {
	logger         log.Logger    // Объект для печати логов
	conn           *net.TCPConn  // Объект TCP-соединения
	enc            *json.Encoder // Объект для кодирования и отправки сообщений
	name           string
	ip             string
	port           string
	nextClient     *Client
	previousClient *Client
	pool           *ClientPool // Пулл клиентов
}

type ClientConfig struct {
	name string
	ip   string
	port string
	pool *ClientPool
}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.
func NewClient(conn *net.TCPConn, cfg *ClientConfig) *Client {
	client := Client{
		logger: *log.New(os.Stdout, "RES", log.LstdFlags),
		conn:   conn,
		enc:    json.NewEncoder(conn),
		name:   cfg.name,
		ip:     cfg.ip,
		port:   cfg.port,
		pool:   cfg.pool,
	}
	client.pool.Clients = append(client.pool.Clients, &client)
	if len(client.pool.Clients) == 2 {
		client.pool.Clients[0].nextClient = client.pool.Clients[1]
		client.pool.Clients[0].previousClient = client.pool.Clients[1]
		client.pool.Clients[1].nextClient = client.pool.Clients[0]
		client.pool.Clients[1].previousClient = client.pool.Clients[0]
	}
	if len(client.pool.Clients) > 2 {
		client.pool.Clients[len(client.pool.Clients)-1].nextClient = client.pool.Clients[0]
		client.pool.Clients[len(client.pool.Clients)-1].previousClient = client.pool.Clients[len(client.pool.Clients)-2]
		client.pool.Clients[len(client.pool.Clients)-1].previousClient.nextClient = client.pool.Clients[len(client.pool.Clients)-1]
		client.pool.Clients[0].previousClient = client.pool.Clients[len(client.pool.Clients)-1]
	}
	return &client
}

// serve - метод, в котором реализован цикл взаимодействия с клиентом.
// Подразумевается, что метод serve будет вызаваться в отдельной go-программе.
func (client *Client) serve() {
	defer client.conn.Close()
	decoder := json.NewDecoder(client.conn)
	for {
		var req proto.Request
		if err := decoder.Decode(&req); err != nil {
			client.logger.Println("cannot decode message", "reason", err)
			break
		} else {
			client.logger.Println("received command", "command", req.Command)
			if client.handleRequest(&req) {
				client.logger.Println("shutting down connection")
				break
			}
		}
	}
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func (client *Client) handleRequest(req *proto.Request) bool {
	switch req.Command {
	case "quit":
		client.respond("ok", nil)
		return true
	case "name":
		client.respond("resultName", client.name)
	case "ip":
		client.respond("resultIp", client.ip)
	case "port":
		client.respond("resultPort", client.port)
	case "nextIp":
		client.respond("resultNextIp", client.nextClient.ip)
	case "nextPort":
		client.respond("resultNextPort", client.nextClient.port)
	case "invisible":
		client.previousClient.nextClient = client.nextClient
		client.nextClient.previousClient = client.previousClient
		client.respond("resultOfInvisible", "client invisibled")
	case "uninvisible":
		client.previousClient.nextClient = client
		client.nextClient.previousClient = client
		client.respond("resultOfUninvisible", "client uninvisibled")
	default:
		client.logger.Println("unknown command")
		client.respond("failed", "unknown command")
	}
	return false
}

// respond - вспомогательный метод для передачи ответа с указанным статусом
// и данными. Данные могут быть пустыми (data == nil).
func (client *Client) respond(status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	client.enc.Encode(&proto.Response{status, &raw})
}

func ParseAddress(address string) (string, string) {
	ip := address[:9]
	port := address[9:]
	return ip, port
}

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	flag.StringVar(&addrStr, "addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	clientPool := ClientPool{Clients: nil}

	// Разбор адреса, строковое представление которого находится в переменной addrStr.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		log.Fatalln("address resolution failed", "address", addrStr)
	} else {
		log.Println("resolved TCP address", "address", addr.String())

		// Инициация слушания сети на заданном адресе.
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			log.Fatalln("listening failed", "reason", err)
		} else {
			// Цикл приёма входящих соединений.
			for {
				if conn, err := listener.AcceptTCP(); err != nil {
					log.Fatalln("cannot accept connection", "reason", err)
				} else {
					log.Println("accepted connection", "address", conn.RemoteAddr().String())
					address := conn.RemoteAddr().String()

					name := strconv.Itoa(len(clientPool.Clients))
					ip, port := ParseAddress(address)
					clientConfig := ClientConfig{ip: ip, name: name, port: port, pool: &clientPool}
					// Запуск go-программы для обслуживания клиентов.
					go NewClient(conn, &clientConfig).serve()
				}
			}
		}
	}
}
