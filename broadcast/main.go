package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type topologyMsg struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

type serverBroadcastMsg struct {
	Type          string          `json:"type"`
	Message       int             `json:"message"`
	NotifiedNodes map[string]bool `json:"notified_nodes"`
}

func main() {
	n := maelstrom.NewNode()
	msgs := make(map[int]bool)
	var neighbours []string
	// logger := log.New(os.Stderr, "Knob: ", 0)
	var mu sync.Mutex

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// logger.Println("Message is from: ", msg.Src)

		body := serverBroadcastMsg{
			Type:          "",
			Message:       0,
			NotifiedNodes: make(map[string]bool),
		}

		if msg.Src[0] == 'c' {
			res := make(map[string]string)
			res["type"] = "broadcast_ok"
			err := n.Reply(msg, res)
			if err != nil {
				return err
			}
		}

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if _, ok := msgs[body.Message]; ok {
			return nil
		}

		mu.Lock()
		msgs[body.Message] = true
		mu.Unlock()

		go sendMessage(n, body, neighbours)

		return nil
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		body := make(map[string]any)
		var msgList []int

		mu.Lock()
		for k := range msgs {
			msgList = append(msgList, k)
		}
		mu.Unlock()

		body["type"] = "read_ok"
		body["messages"] = msgList

		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		body := topologyMsg{}

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		neighbours = body.Topology[n.ID()]

		res := make(map[string]string)
		res["type"] = "topology_ok"

		return n.Reply(msg, res)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func sendMessage(n *maelstrom.Node, body serverBroadcastMsg, neighbours []string) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	logger := log.New(os.Stderr, "Knob: ", 0)

	neighboursToMsg := make(map[string]bool)
	for _, neighbour := range neighbours {
		if _, alreadyNotified := body.NotifiedNodes[neighbour]; !alreadyNotified {
			neighboursToMsg[neighbour] = true
		}
		body.NotifiedNodes[neighbour] = true
	}

	body.NotifiedNodes[n.ID()] = true

	// for {
	// mu.Lock()
	for neighbour := range neighboursToMsg {
		wg.Add(1)
		n.RPC(neighbour, body, func(msg maelstrom.Message) error {
			mu.Lock()
			delete(neighboursToMsg, neighbour)
			mu.Unlock()
			wg.Done()
			return nil
		})
	}
	logger.Println(neighboursToMsg)
	// mu.Unlock()

	// time.Sleep(5 * time.Second)

	// mu.Lock()
	// if len(neighboursToMsg) <= 0 {
	// break
	// return
	// }
	// mu.Unlock()
	// }
	wg.Wait()
}

// func remove(l []string, item string) []string {
// 	for i, v := range l {
// 		if v == item {
// 			return append(l[:i], l[i+1:]...)
// 		}
// 	}
// 	return nil
// }
