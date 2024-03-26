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
	Type      string   `json:"type"`
	Message   int      `json:"message"`
	Receivers []string `json:"receivers"`
}

func main() {
	n := maelstrom.NewNode()
	var msgs []int
	var neighbours []string
	// messagesAndReceivers := make(map[int][]string)
	logger := log.New(os.Stderr, "Knob: ", 0)

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		logger.Println("Message is from: ", msg.Src)

		body := serverBroadcastMsg{}

		if msg.Src[0] == 'c' {
			// var body map[string]any

			// if err := json.Unmarshal(msg.Body, &body); err != nil {
			// 	return err
			// }

			// body["type"] = "broadcast_ok"
			// delete(body, "message")

			res := make(map[string]string)
			res["type"] = "broadcast_ok"
			n.Reply(msg, res)
		}

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body.Receivers = append(body.Receivers, n.ID())
		// messagesAndReceivers[body.Message] = append(messagesAndReceivers[body.Message], body.Receivers...)

		// msgs = append(msgs, int(body["message"].(float64)))
		msgs = append(msgs, body.Message)

		// res.MessagesAndReceivers[]
		// res[int(body["message"].(float64))] = append(res[int(body["message"].(float64))], msg.Src)

		// Send to neighbour server nodes
		shouldSend := true
		receivers := body.Receivers
		body.Receivers = append(body.Receivers, neighbours...)

		for _, neighbour := range neighbours {
			shouldSend = true
			for _, receiver := range receivers {
				if neighbour == receiver {
					shouldSend = false
					break
				}
			}

			if shouldSend {
				n.Send(neighbour, body)
			}

			// if msg.Src != neighbour {
			// 	n.Send(neighbour, body)
			// }
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
