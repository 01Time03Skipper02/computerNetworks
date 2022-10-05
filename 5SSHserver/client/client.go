package main

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
)

const (
	IP       = "151.248.113.144"
	PORT     = "443"
	LOGIN    = "test"
	PASSWORD = "SDHBCXdsedfs222"
)

func main() {
	cfg := &ssh.ClientConfig{
		User: LOGIN,
		Auth: []ssh.AuthMethod{
			ssh.Password(PASSWORD),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	conn, err := ssh.Dial("tcp", IP+":"+PORT, cfg)
	if err != nil {
		log.Fatalln(errors.New("error connection"))
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalln(errors.New("error creation of session"))
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalln(errors.New("error stdin path"))
	}

	go io.Copy(stdin, os.Stdin)
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err = session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatalln(errors.New("error request to terminal"))
	}

	if err = session.Shell(); err != nil {
		log.Fatalln(errors.New("error start shell"))
	}

	if err = session.Wait(); err != nil {
		log.Fatalln(errors.New("error of returning"))
	}
}
