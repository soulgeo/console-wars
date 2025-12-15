package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const Port = ":4567"

func handleConnection(conn net.Conn) {
	defer conn.Close()
	defer fmt.Printf("Connection closed.\n")
	fmt.Printf("Connection accepted.\n")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("read: %s\n", text)
		writer := bufio.NewWriter(conn)
		_, err := writer.WriteString("Received.\n")
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		writer.Flush()
	}
}

func main() {
	listener, err := net.Listen("tcp", Port)
	if err != nil {
		fmt.Printf("Error with listener.\n")
	}
	fmt.Printf("Listening on %s\n", Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		go handleConnection(conn)
	}
}
