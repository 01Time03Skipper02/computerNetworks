package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func serveClient(response http.ResponseWriter, request *http.Request) {
	req, err := http.Get("https://kibers.com/courses.html")
	if err != nil {
		log.Fatalln(err)
	}
	defer req.Body.Close()
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	htmlCode := string(b)
	start, stop := 0, 0
	keyword, word := "tbody", ""
	flag := false
	for i := 0; i < len(htmlCode); i++ {
		letter := string(htmlCode[i])
		word += letter
		if strings.HasPrefix(keyword, word) {
			if strings.Compare(word, keyword) == 0 {
				if !flag {
					start = i - 5
					flag = true
				} else {
					stop = i + 1
					flag = false
				}
				word = ""
			} else {
				continue
			}
		} else {
			word = ""
			continue
		}
	}

	wallet, price, cnt := "", "", 0
	keyword1, keyword2, word := `<div class="courses_table_name">`, `<div class="courses_table_cost">`, ""
	htmlCode = htmlCode[start:stop]
	start, stop = 0, 0
	for i := 0; i < len(htmlCode); i++ {
		letter := string(htmlCode[i])
		word += letter
		if strings.HasPrefix(keyword1, word) {
			if strings.Compare(word, keyword1) == 0 {
				i++
				for ; string(htmlCode[i]) != "<"; i++ {
					wallet += string(htmlCode[i])
				}
				cnt++
				word = ""
			} else {
				continue
			}
		} else if strings.HasPrefix(keyword2, word) {
			if strings.Compare(word, keyword2) == 0 {
				i++
				for ; string(htmlCode[i]) != "â"; i++ {
					price += string(htmlCode[i])
				}
				cnt++
				word = ""
			}
		} else {
			word = ""
			continue
		}

		if cnt == 2 {
			res := wallet + " = " + price + "RUB" + "\n"
			fmt.Println(price)
			response.Write([]byte(res))
			cnt = 0
			wallet, price = "", ""
		}
	}
}

func main() {
	http.HandleFunc("/", serveClient)
	log.Println("starting listener")
	log.Fatalln("listener failed", "error", http.ListenAndServe("127.0.0.1:6060", nil))
}
