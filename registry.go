package main

import (
	"net"
	"sync"
)

type Registry interface {
	addClient(client *Client)
	followClient(id string, followerID string)
	unfollowClient(id string, followerID string)
	sendMessageToClient(id string, message []byte)
	sendMessageToAllClients(message []byte)
	sendMessageToFollower(id string, message []byte)
}

type Client struct {
	connection net.Conn
	id         string
}

type clientReqistry struct {
	sync.Mutex
	clients     map[string]*Client
	followerMap map[string][]string
}

func (registry *clientReqistry) addClient(client *Client) {
	registry.Lock()
	defer registry.Unlock()
	registry.clients[client.id] = client
}

func (registry *clientReqistry) followClient(id string, followerID string) {
	registry.Lock()
	defer registry.Unlock()
	follower := registry.followerMap[id]
	if !containsFollower(followerID, follower) {
		registry.followerMap[id] = append(registry.followerMap[id], followerID)
	}

}

func containsFollower(id string, followers []string) bool {
	for _, follower := range followers {
		if follower == id {
			return true
		}
	}
	return false
}

func (registry *clientReqistry) unfollowClient(id string, followerID string) {
	registry.Lock()
	defer registry.Unlock()
	for i, follower := range registry.followerMap[id] {
		if follower == followerID {
			registry.followerMap[id] = removeArrayEntry(i, registry.followerMap[id])
			return
		}
	}
}

func removeArrayEntry(position int, array []string) []string {
	array[position] = array[len(array)-1]
	return array[:len(array)-1]
}

func (registry *clientReqistry) sendMessageToClient(id string, message []byte) {
	registry.Lock()
	defer registry.Unlock()

	if client, ok := registry.clients[id]; ok {
		client.connection.Write(message)
	}
}

func (registry *clientReqistry) sendMessageToAllClients(message []byte) {
	registry.Lock()
	defer registry.Unlock()

	for _, client := range registry.clients {
		client.connection.Write(message)
	}
}

func (registry *clientReqistry) sendMessageToFollower(id string, message []byte) {
	registry.Lock()
	defer registry.Unlock()

	for _, follower := range registry.followerMap[id] {
		if followerClient, ok := registry.clients[follower]; ok {
			followerClient.connection.Write(message)
		}
	}
}
