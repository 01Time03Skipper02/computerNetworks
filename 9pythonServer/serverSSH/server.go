package main

import (
	"bufio"
	"github.com/gliderlabs/ssh"
)

func handler(s ssh.Session) {
	for {
		scanner := bufio.NewScanner(s)
		for scanner.Scan() {
			line := scanner.Text()

		}
	}
}

func main() {

}
