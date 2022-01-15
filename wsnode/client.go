package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gorilla/websocket"
)

type Client struct {
	uid string

	hub  *Hub
	conn *websocket.Conn
	send chan Message
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Legge dai client e scrive su Kafka, scrive anche heartbeat
func (c *Client) Reader() {
	defer c.hub.Unregister(c)

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	connectMessage := Message{UID: c.uid, MessageType: "connect"}
	c.SendKafkaMessage(connectMessage)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var m Message
		if err := json.Unmarshal(message, &m); err == nil {
			// Scrivi messaggio su kafka
			m.UID = c.uid
			c.SendKafkaMessage(m)
			fmt.Println("Messaggio inbound inviato.")
		}
	}

	disconnectMessage := Message{UID: c.uid, MessageType: "disconnect"}
	c.SendKafkaMessage(disconnectMessage)
}

// Invia ai client
func (c *Client) Writer() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendKafkaMessage(data interface{}) {
	md, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("[WSNODE] Marshalling error while sending in Kafka: %v\n", err)
		return
	}

	kafkaMessage := sarama.ProducerMessage{
		Topic: KAFKA_INBOUND_TOPIC,
		Key:   nil,
		Value: sarama.ByteEncoder(md),
	}
	c.hub.kafkaProducer.Input() <- &kafkaMessage
}
