package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/skorobogatov/input"
	"net"
)

import "proto"

// interact - функция, содержащая цикл взаимодействия с сервером.
func interact(conn *net.TCPConn) {
	defer conn.Close()
	encoder, decoder := json.NewEncoder(conn), json.NewDecoder(conn)
	for {
		// Чтение команды из стандартного потока ввода
		fmt.Printf("command = ")
		command := input.Gets()

		// Отправка запроса.
		switch command {
		case "quit":
			send_request(encoder, "quit", nil)
			return
		case "name":
			send_request(encoder, "name", nil)
		case "ip":
			send_request(encoder, "ip", nil)
		case "port":
			send_request(encoder, "port", nil)
		case "nextIp":
			send_request(encoder, "nextIp", nil)
		case "nextPort":
			send_request(encoder, "nextPort", nil)
		case "invisible":
			send_request(encoder, "invisible", nil)
		case "uninvisible":
			send_request(encoder, "uninvisible", nil)
		default:
			fmt.Printf("error: unknown command\n")
			continue
		}

		// Получение ответа.
		var resp proto.Response
		if err := decoder.Decode(&resp); err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}

		// Вывод ответа в стандартный поток вывода.
		switch resp.Status {
		case "ok":
			fmt.Printf("ok\n")
		case "failed":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var errorMsg string
				if err := json.Unmarshal(*resp.Data, &errorMsg); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Printf("failed: %s\n", errorMsg)
				}
			}
		case "resultName":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var name string
				if err := json.Unmarshal(*resp.Data, &name); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Println(name)
				}
			}
		case "resultPort":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var port string
				if err := json.Unmarshal(*resp.Data, &port); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Println(port)
				}
			}
		case "resultIp":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var ip string
				if err := json.Unmarshal(*resp.Data, &ip); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Println(ip)
				}
			}
		case "resultNextPort":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var port string
				if err := json.Unmarshal(*resp.Data, &port); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Println(port)
				}
			}
		case "resultNextIp":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var ip string
				if err := json.Unmarshal(*resp.Data, &ip); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Println(ip)
				}
			}
		case "resultOfInvisible":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var res string
				if err := json.Unmarshal(*resp.Data, &res); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Println(res)
				}
			}
		case "resultOfUninvisible":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var res string
				if err := json.Unmarshal(*resp.Data, &res); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Println(res)
				}
			}
		default:
			fmt.Printf("error: server reports unknown status %q\n", resp.Status)
		}
	}
}

// send_request - вспомогательная функция для передачи запроса с указанной командой
// и данными. Данные могут быть пустыми (data == nil).
func send_request(encoder *json.Encoder, command string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	encoder.Encode(&proto.Request{command, &raw})
}

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	flag.StringVar(&addrStr, "addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	// Разбор адреса, установка соединения с сервером и
	// запуск цикла взаимодействия с сервером.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		fmt.Printf("error: %v\n", err)
	} else if conn, err := net.DialTCP("tcp", nil, addr); err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		interact(conn)
	}
}
