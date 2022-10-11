package main

import (
	"flag"
	"github.com/goftp/file-driver"
	"github.com/goftp/server"
	"log"
)

func main() {
	var (
		rootDir  = flag.String("rootDir", "/mnt/c/Users/1/Desktop/studies/2/отчеты по сетям/6", "rootdir")
		userName = flag.String("userName", "time_skipper", "Username")
		password = flag.String("password", "12345678990", "Password")
		port     = flag.Int("port", 2121, "Port")
		host     = flag.String("host", "localhost", "Host")
	)
	flag.Parse()
	if *rootDir == "" {
		log.Fatalf("no rootDir. Set this one")
	}
	factory := &filedriver.FileDriverFactory{
		RootPath: *rootDir,
		Perm:     server.NewSimplePerm("user", "group"),
	}
	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		Auth:     &server.SimpleAuth{Name: *userName, Password: *password},
	}

	log.Printf("Starting server on %v:%v", opts.Hostname, opts.Port)
	log.Printf("Username %v, Password %v", *userName, *password)
	serv := server.NewServer(opts)
	err := serv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server with error: ", err)
	}
}
