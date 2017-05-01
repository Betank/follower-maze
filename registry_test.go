package main

import "testing"

func TestFollowClient(t *testing.T) {
	registry := &clientReqistry{clients: make(map[string]*Client), followerMap: make(map[string][]string)}
	client := &Client{id: "1"}
	registry.addClient(client)

	registry.followClient("1", "2")

	if registry.followerMap["1"][0] != "2" {
		t.Errorf("client with id %s should have follower with id %s", "1", "2")
	}
}

func TestSameClientFollows2Times(t *testing.T) {
	registry := &clientReqistry{clients: make(map[string]*Client), followerMap: make(map[string][]string)}
	client := &Client{id: "1"}
	registry.addClient(client)

	registry.followClient("1", "2")
	registry.followClient("1", "2")

	if len(registry.followerMap["1"]) > 1 {
		t.Errorf("client should only have %d follower, but has %d", 1, len(registry.followerMap["1"]))
	}
}

func TestUnfollowClient(t *testing.T) {
	registry := &clientReqistry{clients: make(map[string]*Client), followerMap: make(map[string][]string)}
	client := &Client{id: "1"}
	registry.addClient(client)

	registry.followClient("1", "2")
	registry.unfollowClient("1", "2")

	if len(registry.followerMap["1"]) != 0 {
		t.Errorf("client should only have %d follower, but has %d", 0, len(registry.followerMap["1"]))
	}
}
