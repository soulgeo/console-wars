package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings" // Add this import
	"time"

	"github.com/soulgeo/console-wars/internal/config"
	"github.com/soulgeo/console-wars/internal/game"
	"github.com/soulgeo/console-wars/internal/messages"
)

type Client struct {
	Conn    net.Conn
	Reader  *bufio.Reader
	Player  game.Player
	MsgChan chan string
}

func matchConnections(c chan net.Conn) {
	var pending *Client
	for {
		newConn := <-c
		msgChan := make(chan string, 10)
		newClient := &Client{
			Conn:    newConn,
			Reader:  bufio.NewReader(newConn),
			MsgChan: msgChan,
			Player:  game.Player{},
		}
		log.Printf("New connection: %s", newConn.RemoteAddr())

		// Read the player's name
		name, err := newClient.Reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading player name: %v", err)
			newClient.Conn.Close()
			close(newClient.MsgChan)
			continue
		}
		newClient.Player.Name = strings.TrimSpace(name)

		go writeMessages(newClient)
		newClient.MsgChan <- messages.Connected

		if pending == nil {
			pending = newClient
			newClient.MsgChan <- messages.Waiting
			log.Printf("Client %s is now waiting...", newConn.RemoteAddr())
			continue
		}

		// IMPORTANT:
		// Check if pending connection is closed when the new connection arrives.
		err = pending.Conn.SetReadDeadline(
			time.Now().Add(10 * time.Millisecond),
		)
		if err != nil {
			log.Printf("Error setting deadline: %v", err)
		}
		_, err = pending.Reader.Peek(1)
		pending.Conn.SetReadDeadline(time.Time{})

		if err == io.EOF {
			log.Printf(
				"Waiting client %s disconnected. Dropping.",
				pending.Conn.RemoteAddr(),
			)
			pending.Conn.Close()
			close(pending.MsgChan)
			pending = newClient
			log.Printf(
				"Client %s promoted to waiting.",
				pending.Conn.RemoteAddr(),
			)
			newClient.MsgChan <- messages.Waiting
			continue
		} else if err != nil {
			// A timeout error is good here. It means "No data, but also no EOF".
			// Check if it's strictly a timeout.
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Alive and silent. Proceed to match.
			} else {
				// Some other error (connection reset, etc). Drop pending.
				log.Printf("Waiting client error (%s). Dropping.", err)
				pending.Conn.Close()
				close(pending.MsgChan)
				pending = newClient
				continue
			}
		}

		// Both clients connected, proceed.
		newClient.MsgChan <- messages.MatchFound
		pending.MsgChan <- messages.MatchFound
		go handleConnections(pending, newClient)
		pending = nil
	}
}

func writeMessages(client *Client) {
	writer := bufio.NewWriter(client.Conn)
	for m := range client.MsgChan {
		_, err := writer.WriteString(m)
		if err != nil {
			log.Printf("error: %s", err.Error())
			return
		}
		writer.Flush()
	}
}

func handleConnections(c1, c2 *Client) {
	defer c1.Conn.Close()
	defer c2.Conn.Close()
	defer close(c1.MsgChan)
	defer close(c2.MsgChan)
	log.Printf(
		"Match found: %s <-> %s",
		c1.Conn.RemoteAddr(),
		c2.Conn.RemoteAddr(),
	)

	logs := make(chan string, 10)
	act1 := make(chan string, 10)
	act2 := make(chan string, 10)

	go game.Play(&c1.Player, &c2.Player, logs, act1, act2)
	go receiveActions(c1.Conn, act1)
	go receiveActions(c2.Conn, act2)

	for l := range logs {
		c1.MsgChan <- l
		c2.MsgChan <- l
	}

	log.Printf(
		"Closing session between %s and %s",
		c1.Conn.RemoteAddr(),
		c2.Conn.RemoteAddr(),
	)
}

func receiveActions(conn net.Conn, action chan string) {
	defer close(action)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		action <- text
	}
}

func main() {
	listener, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("Error with listener: %s", err)
	}
	c := make(chan net.Conn, 100)
	go matchConnections(c)

	log.Printf("Listening on %s", config.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %s", err)
			continue
		}
		c <- conn
	}
}
