package main

import (
	"testing"
)

func TestProcessFollowMessageClient50(t *testing.T) {
	count = 0
	message := []byte("1|F|60|50\r\n")

	mockReg := &mockRegistry{
		sendMessageCalls:  make([]sendMessageCall, 0),
		followClientCalls: make([]followClientCall, 0),
	}
	registry = mockReg

	processMessage(message)

	expectedCall := sendMessageCall{"50", []byte(message)}
	expectedCall2 := followClientCall{"50", "60"}

	assertSendMessageCalled(mockReg.sendMessageCalls[0], expectedCall, t)
	assertFollowClientCalled(mockReg.followClientCalls[0], expectedCall2, t)
}

func TestProcessFollowMessageClient51(t *testing.T) {
	count = 0
	message := []byte("1|F|60|51\r\n")

	mockReg := &mockRegistry{
		sendMessageCalls:  make([]sendMessageCall, 0),
		followClientCalls: make([]followClientCall, 0),
	}
	registry = mockReg

	processMessage(message)

	expectedCall := sendMessageCall{"51", []byte(message)}
	expectedCall2 := followClientCall{"51", "60"}

	assertSendMessageCalled(mockReg.sendMessageCalls[0], expectedCall, t)
	assertFollowClientCalled(mockReg.followClientCalls[0], expectedCall2, t)
}

func TestProcessFollowMessageWithOrder(t *testing.T) {
	count = 0
	message1 := []byte("1|F|60|50\r\n")
	message2 := []byte("2|F|60|50\r\n")
	message3 := []byte("3|F|60|50\r\n")

	mockReg := &mockRegistry{
		sendMessageCalls:  make([]sendMessageCall, 0),
		followClientCalls: make([]followClientCall, 0),
	}
	registry = mockReg

	processMessage(message3)
	processMessage(message2)
	processMessage(message1)

	expectedCall1 := sendMessageCall{"50", []byte(message1)}
	expectedCall2 := sendMessageCall{"50", []byte(message2)}
	expectedCall3 := sendMessageCall{"50", []byte(message3)}

	expectedCall4 := followClientCall{"50", "60"}
	expectedCall5 := followClientCall{"50", "60"}
	expectedCall6 := followClientCall{"50", "60"}

	assertSendMessageCalled(mockReg.sendMessageCalls[0], expectedCall1, t)
	assertSendMessageCalled(mockReg.sendMessageCalls[1], expectedCall2, t)
	assertSendMessageCalled(mockReg.sendMessageCalls[2], expectedCall3, t)

	assertFollowClientCalled(mockReg.followClientCalls[0], expectedCall4, t)
	assertFollowClientCalled(mockReg.followClientCalls[1], expectedCall5, t)
	assertFollowClientCalled(mockReg.followClientCalls[2], expectedCall6, t)
}

func TestProcessBroadcastMessage(t *testing.T) {
	count = 0
	message := []byte("1|B\r\n")

	mockReg := &mockRegistry{sendMessageCalls: make([]sendMessageCall, 0)}
	registry = mockReg

	processMessage(message)

	if string(mockReg.sendMessageToAllCalls[0]) != string(message) {
		t.Errorf("should call method with %v but called with %v", string(message), string(mockReg.sendMessageToAllCalls[0]))
	}
}

func TestProcessStatusUpdateMessage(t *testing.T) {
	count = 0
	message := []byte("1|S|50\r\n")

	mockReg := &mockRegistry{sendMessageToFollowerCalls: make([]sendMessageCall, 0)}
	registry = mockReg

	processMessage(message)

	expectedCall := sendMessageCall{"50", []byte(message)}

	assertSendMessageCalled(mockReg.sendMessageToFollowerCalls[0], expectedCall, t)
}

func TestProcessUnfollowMessage(t *testing.T) {
	count = 0
	message := []byte("1|U|60|50\r\n")

	mockReg := &mockRegistry{
		followClientCalls:   make([]followClientCall, 0),
		unfollowClientCalls: make([]unfollowClientCall, 0)}
	registry = mockReg

	processMessage(message)

	expectedCall := unfollowClientCall{"50", "60"}

	assertUnfollowClientCalled(mockReg.unfollowClientCalls[0], expectedCall, t)
}

func assertSendMessageCalled(call, expected sendMessageCall, t *testing.T) {
	if call.id != expected.id {
		t.Errorf("should call method with %v but called with %v", expected.id, call.id)
	}

	if string(call.message) != string(expected.message) {
		t.Errorf("should call method with %v but called with %v", expected.message, call.message)
	}
}

func assertFollowClientCalled(call, expected followClientCall, t *testing.T) {
	if call.id != expected.id {
		t.Errorf("should call method with %v but called with %v", expected.id, call.id)
	}

	if call.followerID != expected.followerID {
		t.Errorf("should call method with %v but called with %v", expected.id, call.id)
	}
}

func assertUnfollowClientCalled(call, expected unfollowClientCall, t *testing.T) {
	if call.id != expected.id {
		t.Errorf("should call method with %v but called with %v", expected.id, call.id)
	}

	if call.followerID != expected.followerID {
		t.Errorf("should call method with %v but called with %v", expected.id, call.id)
	}
}

type mockRegistry struct {
	followClientCalls          []followClientCall
	unfollowClientCalls        []unfollowClientCall
	sendMessageCalls           []sendMessageCall
	sendMessageToAllCalls      [][]byte
	sendMessageToFollowerCalls []sendMessageCall
}

type followClientCall struct {
	id         string
	followerID string
}

type unfollowClientCall followClientCall

type sendMessageCall struct {
	id      string
	message []byte
}

func (registry *mockRegistry) addClient(client *Client) {
}

func (registry *mockRegistry) followClient(id string, followerID string) {
	registry.followClientCalls = append(registry.followClientCalls, followClientCall{id, followerID})
}

func (registry *mockRegistry) unfollowClient(id string, followerID string) {
	registry.unfollowClientCalls = append(registry.unfollowClientCalls, unfollowClientCall{id, followerID})
}

func (registry *mockRegistry) sendMessageToClient(id string, message []byte) {
	registry.sendMessageCalls = append(registry.sendMessageCalls, sendMessageCall{id, message})
}

func (registry *mockRegistry) sendMessageToAllClients(message []byte) {
	registry.sendMessageToAllCalls = append(registry.sendMessageToAllCalls, message)
}

func (registry *mockRegistry) sendMessageToFollower(id string, message []byte) {
	registry.sendMessageToFollowerCalls = append(registry.sendMessageToFollowerCalls, sendMessageCall{id, message})
}
