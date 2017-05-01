package main

import (
	"net"
	"testing"
	"time"
)

func TestFollowWithClientAvailable(t *testing.T) {
	event := "1|F|60|50\r\n"

	eventsource := connect("tcp", "localhost:9090", t)
	defer eventsource.Close()

	client := connect("tcp", "localhost:9099", t)
	defer client.Close()

	client.Write([]byte("50\r\n"))
	time.Sleep(50 * time.Millisecond)

	eventsource.Write([]byte(event))

	message := awaitMessage(client, 5)
	if message != event {
		t.Errorf("got %s message but want %s", message, event)
	}
}

func TestFollowWithClientNotAvailable(t *testing.T) {
	event := "2|F|60|50\r\n"

	eventsource := connect("tcp", "localhost:9090", t)
	defer eventsource.Close()

	client := connect("tcp", "localhost:9099", t)
	defer client.Close()

	client.Write([]byte("1\r\n"))
	time.Sleep(50 * time.Millisecond)

	eventsource.Write([]byte(event))

	message := awaitMessage(client, 2)
	if message != "" {
		t.Errorf("should not get message but got %s", message)
	}
}

func TestMultipleFollowMessages(t *testing.T) {
	event1 := "3|F|60|50\r\n"
	event2 := "4|F|60|51\r\n"

	eventsource := connect("tcp", "localhost:9090", t)
	defer eventsource.Close()

	client1 := connect("tcp", "localhost:9099", t)
	defer client1.Close()
	client2 := connect("tcp", "localhost:9099", t)
	defer client2.Close()

	client1.Write([]byte("50\r\n"))
	time.Sleep(50 * time.Millisecond)
	client2.Write([]byte("51\r\n"))
	time.Sleep(50 * time.Millisecond)

	eventsource.Write([]byte(event1))
	eventsource.Write([]byte(event2))

	message1 := awaitMessage(client1, 2)
	if message1 != event1 {
		t.Errorf("got %s message but want %s", message1, event1)
	}

	message2 := awaitMessage(client2, 2)
	if message2 != event2 {
		t.Errorf("got %s message but want %s", message2, event2)
	}
}

func TestBroadcastMessage(t *testing.T) {
	event := "5|B\r\n"

	eventsource := connect("tcp", "localhost:9090", t)
	defer eventsource.Close()

	client1 := connect("tcp", "localhost:9099", t)
	defer client1.Close()
	client2 := connect("tcp", "localhost:9099", t)
	defer client2.Close()

	client1.Write([]byte("50\r\n"))
	time.Sleep(50 * time.Millisecond)
	client2.Write([]byte("51\r\n"))
	time.Sleep(50 * time.Millisecond)

	eventsource.Write([]byte(event))

	message1 := awaitMessage(client1, 2)
	if message1 != event {
		t.Errorf("got %s message but want %s", message1, event)
	}

	message2 := awaitMessage(client2, 2)
	if message2 != event {
		t.Errorf("got %s message but want %s", message2, event)
	}
}

func TestSendPrivateMessage(t *testing.T) {
	event := "6|P|60|50\r\n"

	eventsource := connect("tcp", "localhost:9090", t)
	defer eventsource.Close()

	client := connect("tcp", "localhost:9099", t)
	defer client.Close()

	client.Write([]byte("50\r\n"))
	time.Sleep(50 * time.Millisecond)

	eventsource.Write([]byte(event))

	message := awaitMessage(client, 5)
	if message != event {
		t.Errorf("got %s message but want %s", message, event)
	}
}

func TestStatusUpdateMessage(t *testing.T) {
	event1 := "7|F|60|50\r\n"
	event2 := "8|S|50\r\n"

	eventsource := connect("tcp", "localhost:9090", t)
	defer eventsource.Close()

	client1 := connect("tcp", "localhost:9099", t)
	defer client1.Close()
	client2 := connect("tcp", "localhost:9099", t)
	defer client2.Close()

	client1.Write([]byte("50\r\n"))
	time.Sleep(50 * time.Millisecond)
	client2.Write([]byte("60\r\n"))
	time.Sleep(50 * time.Millisecond)

	eventsource.Write([]byte(event1))
	eventsource.Write([]byte(event2))

	awaitMessage(client1, 2)
	message := awaitMessage(client2, 2)
	if message != event2 {
		t.Errorf("got %s message but want %s", message, event2)
	}
}

func TestUnfollowMessage(t *testing.T) {
	event1 := "9|F|60|50\r\n"
	event2 := "10|U|60|50\r\n"
	event3 := "11|S|50\r\n"

	eventsource := connect("tcp", "localhost:9090", t)
	defer eventsource.Close()

	client1 := connect("tcp", "localhost:9099", t)
	defer client1.Close()
	client2 := connect("tcp", "localhost:9099", t)
	defer client2.Close()

	client1.Write([]byte("50\r\n"))
	time.Sleep(50 * time.Millisecond)
	client2.Write([]byte("60\r\n"))
	time.Sleep(50 * time.Millisecond)

	eventsource.Write([]byte(event1))
	eventsource.Write([]byte(event2))
	eventsource.Write([]byte(event3))

	awaitMessage(client1, 2)
	message := awaitMessage(client2, 2)
	if message != "" {
		t.Errorf("should not get message but got %s", message)
	}
}

func connect(network, address string, t *testing.T) net.Conn {
	client, err := net.Dial(network, address)
	if err != nil {
		t.Error(err)
	}
	return client
}

func awaitMessage(client net.Conn, timeout time.Duration) string {
	stop := make(chan string, 1)

	go func() {
		message := make([]byte, 128)
		length, err := client.Read(message)
		if err != nil {
			stop <- ""
		}
		stop <- string(message[:length])
	}()
	go func() {
		time.Sleep(timeout * time.Second)
		stop <- ""
	}()

	return <-stop
}
