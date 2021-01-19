package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Client ...
type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

// Message ...
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Type: messageType, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received:%+v\n", message)
	}
}

// Start ...
func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true
			fmt.Println("Size of Connection Pool:", len(p.Clients))
			for client, _ := range p.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined ..."})
			}
			break
		case client := <-p.Unregister:
			delete(p.Clients, client)
			fmt.Println("Size of Connection Pool:", len(p.Clients))
			for client, _ := range p.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected ..."})
			}
			break
		case message := <-p.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range p.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
