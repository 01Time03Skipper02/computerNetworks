package main

import (
	"errors"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os/exec"
	"strings"
)

func main() {
	ssh.Handle(func(session ssh.Session) {
		term := terminal.NewTerminal(session, "> ")
		for {
			line, err := term.ReadLine()
			if err != nil {
				log.Fatalln(errors.New("error of reading line"))
			}
			splited := strings.Split(line, " ")
			cmd := exec.Command(splited[0], splited[1:]...)
			cmd.Stdout = session
			cmd.Stderr = session
			if err := cmd.Run(); err != nil {
				log.Fatalln(errors.New("error of running command"))
			}
		}
	})
	log.Println("Starting ssh server on port 6060...")
	err := ssh.ListenAndServe(":6060", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
