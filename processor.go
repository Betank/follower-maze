package main

import (
	"strconv"
	"strings"
	"sync"
)

var Sequence = 0
var Type = 1
var FromUserId = 2
var ToUserId = 3

type MessageProcessor struct {
	sync.Mutex
	messageQueue map[int][]byte
	messageCount int
}

func (processor *MessageProcessor) processMessage(message []byte) {
	processor.Lock()
	defer processor.Unlock()
	fragments := splitMessage(message)
	seq, err := strconv.Atoi(fragments[Sequence])
	if err != nil {
		return
	}
	if seq == processor.messageCount+1 {
		processor.sendMessage(message, fragments)
		processor.processQueueMessages()
	} else {
		processor.pushMessageToQueue(seq, message)
	}
}

func (processor *MessageProcessor) processQueueMessages() {
	for {
		if message, ok := processor.messageQueue[processor.messageCount+1]; ok {
			fragments := splitMessage(message)
			processor.sendMessage(message, fragments)
			delete(processor.messageQueue, processor.messageCount)
		} else {
			break
		}
	}
}

func splitMessage(message []byte) []string {
	return strings.Split(strings.Trim(string(message), "\r\n"), "|")
}

func (processor *MessageProcessor) sendMessage(message []byte, fragments []string) {
	sendMessage(message, fragments)
	processor.messageCount++
}

func sendMessage(message []byte, fragments []string) {
	if len(fragments) == 2 {
		registry.sendMessageToAllClients(message)
	} else if len(fragments) == 3 {
		registry.sendMessageToFollower(fragments[FromUserId], message)
	} else {
		switch v := fragments[Type]; v {
		case "F":
			registry.followClient(fragments[ToUserId], fragments[FromUserId])
			registry.sendMessageToClient(fragments[ToUserId], message)
		case "U":
			registry.unfollowClient(fragments[ToUserId], fragments[FromUserId])
		case "P":
			registry.sendMessageToClient(fragments[ToUserId], message)

		}
	}
}

func (processor *MessageProcessor) pushMessageToQueue(sequence int, message []byte) {
	processor.messageQueue[sequence] = append(processor.messageQueue[sequence], message...)
}
