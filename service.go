package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex
var registry Registry
var count int
var messageStore = make(map[int][]byte)

func main() {
	registry = &clientReqistry{clients: make(map[string]*Client), followerMap: make(map[string][]string)}
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
			processMessage(message)
			message = message[:0]
		}
	}
}

func processMessage(message []byte) {
	mu.Lock()
	defer mu.Unlock()
	fragments := strings.Split(strings.Trim(string(message), "\r\n"), "|")
	seq, err := strconv.Atoi(fragments[0])
	if err != nil {
		return
	}
	if seq == count+1 {
		sendMessage(message, fragments)
		count++
		for {
			if message, ok := messageStore[count+1]; ok {
				fragments = strings.Split(strings.Trim(string(message), "\r\n"), "|")
				sendMessage(message, fragments)
				delete(messageStore, count+1)
				count++
			} else {
				break
			}
		}
	} else {
		messageStore[seq] = append(messageStore[seq], message...)
	}
}

func sendMessage(message []byte, fragments []string) {
	if len(fragments) == 2 {
		registry.sendMessageToAllClients(message)
	} else if len(fragments) == 3 {
		registry.sendMessageToFollower(fragments[2], message)
	} else {
		switch v := fragments[1]; v {
		case "F":
			registry.followClient(fragments[3], fragments[2])
			registry.sendMessageToClient(fragments[3], message)
		case "U":
			registry.unfollowClient(fragments[3], fragments[2])
		case "P":
			registry.sendMessageToClient(fragments[3], message)

		}
	}
}
