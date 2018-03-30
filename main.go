package main

import (
	"go-http2-streaming/http2"
	"io/ioutil"
	"log"
)

func main() {
	// Waitc is used to hold the main function
	// from returning before the client connects to the server.
	waitc := make(chan bool)

	// Reads a file into memory
	data, err := ioutil.ReadFile("./test.json")
	if err != nil {
		log.Println(err)
		return
	}

	// HTTP2 Client
	go func() {
		client := new(http2.Client)
		client.Dial()
		client.Post(data)
	}()

	// HTTP2 Server
	server := new(http2.Server)
	err = server.Initialize()
	if err != nil {
		log.Println(err)
		return
	}

	// Waits forever
	<-waitc
}
