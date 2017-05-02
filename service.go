package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var messageProcessor *MessageProcessor
var registry Registry

func main() {
	registry = &clientReqistry{clients: make(map[string]*Client), followerMap: make(map[string][]string)}
	messageProcessor = &MessageProcessor{messageQueue: make(map[int][]byte)}
	eventsourceListener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}
	clientListener, err := net.Listen("tcp", ":9099")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go handleListener(eventsourceListener, handleEventSource)
	go handleListener(clientListener, handleClient)
	wg.Wait()

}

func handleListener(listener net.Listener, connectionHandler func(conn net.Conn)) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		connectionHandler(conn)
	}
}

func handleEventSource(conn net.Conn) {
	request := make([]byte, 128)
	defer conn.Close()

	input := make(chan byte)
	go chunk(input)
	for {
		length, err := conn.Read(request)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, c := range request[:length] {
			input <- c
		}
	}
}

func handleClient(conn net.Conn) {
	go func() {
		request := make([]byte, 128)
		length, err := conn.Read(request)
		if err != nil {
			fmt.Println(err)
		}
		newClient := &Client{
			connection: conn,
			id:         strings.Trim(string(request[:length]), "\r\n"),
		}
		registry.addClient(newClient)
	}()
}

func chunk(input chan byte) {
	message := make([]byte, 0)

	for in := range input {
		message = append(message, in)
		if in == 10 {
			messageProcessor.processMessage(message)
			message = message[:0]
		}
	}
}
