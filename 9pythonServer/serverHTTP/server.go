package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func handlePython(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	tasks := strings.Split(url, "/")

	if len(tasks) > 1 {
		name := []string{}
		quer := []string{}

		if len(tasks) > 2 {
			name = strings.Split(tasks[1], ".")
			quer = tasks[2:]
		} else {
			name = strings.Split(tasks[1], ".")
			query := r.URL.Query()
			for _, v := range query {
				quer = append(quer, v[0])
			}
		}

		switch name[1] {
		case "png":
			message, err := ioutil.ReadFile("./" + name[1] + "/" + tasks[1])
			if err != nil {
				log.Fatalln(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "image/png")
			w.Write(message)
		case "jpg":
			message, err := ioutil.ReadFile("./" + name[1] + "/" + tasks[1])
			if err != nil {
				log.Fatalln(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "image/jpg")
			w.Write(message)
		case "gif":
			message, err := ioutil.ReadFile("./" + name[1] + "/" + tasks[1])
			if err != nil {
				log.Fatalln(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "image/gif")
			w.Write(message)
		case "txt":
			message, err := ioutil.ReadFile("./" + name[1] + "/" + tasks[1])
			if err != nil {
				log.Fatalln(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/txt; charset=utf-8")
			w.Write(message)
		case "html":
			message, err := ioutil.ReadFile("./" + name[1] + "/" + tasks[1])
			if err != nil {
				log.Fatalln(err)
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, string(message))
		case "py":
			command := exec.Command("python3", "./"+name[1]+"/"+tasks[1])

			var res bytes.Buffer
			command.Stdout = &res
			stdin, err := command.StdinPipe()
			if err != nil {
				log.Fatalln(err)
			}

			go func() {
				defer stdin.Close()
				for i := 0; i < len(quer); i++ {
					io.WriteString(stdin, quer[i]+"\n")
				}
			}()

			compileError := command.Run()

			if compileError != nil {
				log.Fatalln(compileError)
			}

			w.Write(res.Bytes())
		}
	}
}

func main() {
	http.HandleFunc("/", handlePython)
	fmt.Println("Server started listen on 8080 port")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
