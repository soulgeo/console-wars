package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/soulgeo/console-wars/internal/config"
	"github.com/soulgeo/console-wars/internal/game"
)

type Client struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Player game.Player
}

func matchConnections(c chan net.Conn) {
	var pending *Client
	for {
		newConn := <-c
		newClient := &Client{
			Conn:   newConn,
			Reader: bufio.NewReader(newConn),
		}
		fmt.Printf("New connection: %s\n", newConn.RemoteAddr())

		if pending == nil {
			pending = newClient
			fmt.Printf("Client %s is now waiting...\n", newConn.RemoteAddr())
			continue
		}

		// IMPORTANT:
		// Check if pending connection is closed when the new connection arrives.
		err := pending.Conn.SetReadDeadline(
			time.Now().Add(10 * time.Millisecond),
		)
		if err != nil {
			log.Printf("Error setting deadline: %v", err)
		}
		_, err = pending.Reader.Peek(1)
		pending.Conn.SetReadDeadline(time.Time{})

		if err == io.EOF {
			fmt.Printf(
				"Waiting client %s disconnected. Dropping.\n",
				pending.Conn.RemoteAddr(),
			)
			pending.Conn.Close()
			pending = newClient
			fmt.Printf(
				"Client %s promoted to waiting.\n",
				pending.Conn.RemoteAddr(),
			)
			continue
		} else if err != nil {
			// A timeout error is good here. It means "No data, but also no EOF".
			// Check if it's strictly a timeout.
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Alive and silent. Proceed to match.
			} else {
				// Some other error (connection reset, etc). Drop pending.
				fmt.Printf("Waiting client error (%s). Dropping.\n", err)
				pending.Conn.Close()
				pending = newClient
				continue
			}
		}

		// Both clients connected, proceed.
		go handleConnections(pending, newClient)
		pending = nil
	}
}

func handleConnections(c1, c2 *Client) {
	defer c1.Conn.Close()
	defer c2.Conn.Close()
	log.Printf(
		"Match found: %s <-> %s\n",
		c1.Conn.RemoteAddr(),
		c2.Conn.RemoteAddr(),
	)

	done := make(chan struct{})

	logs := make(chan string, 10)
	played := make(chan byte)
	c1.Player = game.Player{}
	c2.Player = game.Player{}

	go game.PlayGame(&c1.Player, &c2.Player, logs, played)

	writer1 := bufio.NewWriter(c1.Conn)
	writer2 := bufio.NewWriter(c2.Conn)
	for l := range logs {
		_, err := writer1.WriteString(l)
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		_, err = writer2.WriteString(l)
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		writer1.Flush()
		writer2.Flush()
	}

	<-done
	fmt.Printf(
		"Closing chat between %s and %s\n",
		c1.Conn.RemoteAddr(),
		c2.Conn.RemoteAddr(),
	)
}

func scanAndSend(r io.Reader, w net.Conn, done chan struct{}) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := fmt.Fprintf(w, "%s\n", text)
		if err != nil {
			break
		}
	}
	select {
	case done <- struct{}{}:
	default:
		// Prevent blocking if both fail simultaneously
	}
}

func main() {
	listener, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("Error with listener: %s", err)
	}
	c := make(chan net.Conn, 100)
	go matchConnections(c)

	fmt.Printf("Listening on %s\n", config.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Accept error: %s\n", err)
		}
		c <- conn
	}
}
