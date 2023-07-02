package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	tcpListener "github.com/goodstemy/redis-on-go/interfaces/tcp"
	memoryStore "github.com/goodstemy/redis-on-go/store/memory"
)

type Message struct {
	Operation string      `json:"operation"`
	Hash      string      `json:"hash"`
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
}

func main() {
	fmt.Println("Started...")

	handler := make(chan bytes.Buffer)
	response := make(chan []byte)
	defer close(response)

	memoryStore.Init()
	go tcpListener.Listen(handler, response)

	for msg := range handler {
		m := Message{}
		err := json.Unmarshal(msg.Bytes(), &m)

		if err != nil {
			response <- []byte(fmt.Sprintf("error on parse message: %v", err))
			continue
		}

		if m.Operation == "" || m.Hash == "" || m.Key == "" || m.Value == nil {
			response <- []byte(fmt.Sprintf("error on parse message: %v", "You should specify operation, key and value"))
			continue
		}

		m.Operation = strings.ToLower(m.Operation)

		err = execute(m, response)

		if err != nil {
			response <- []byte(fmt.Sprintf("%v", err))
		}
	}
}

func execute(msg Message, snd chan []byte) error {
	switch op := msg.Operation; op {
	case "hset":
		memoryStore.Hset(msg.Hash, msg.Key, msg.Value)

		snd <- []byte{49} // response 1 just like in redis if success
	case "hget":
		val := memoryStore.HGet(msg.Hash, msg.Key)

		if val == nil {
			snd <- []byte{48} // response 0 just like in redis if success
			return nil
		}

		valAsBytes, ok := val.(string)

		if !ok {
			return errors.New(fmt.Sprintf("Cannot parse operation response"))
		}

		snd <- []byte(valAsBytes)
	default:
		return errors.New(fmt.Sprintf("Operation %v not found", msg.Operation))
	}

	return nil
}
