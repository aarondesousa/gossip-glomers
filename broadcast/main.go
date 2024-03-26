package main

import (
	"encoding/json"
	"log"
	"os"

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
	var msgs []int
	var neighbours []string
	// messagesAndReceivers := make(map[int][]string)
	logger := log.New(os.Stderr, "Knob: ", 0)

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		logger.Println("Message is from: ", msg.Src)

		body := serverBroadcastMsg{
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

		body.NotifiedNodes[n.ID()] = true

		msgs = append(msgs, body.Message)

		// Send to neighbour server nodes
		for _, neighbour := range neighbours {
			_, msgSent := body.NotifiedNodes[neighbour]

			body.NotifiedNodes[neighbour] = true

			if !msgSent {
				err := n.Send(neighbour, body)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		body := make(map[string]any)
		body["type"] = "read_ok"
		body["messages"] = msgs
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
