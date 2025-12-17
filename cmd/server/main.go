package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const Port = ":4567"

func handleConnections(conn1 net.Conn, conn2 net.Conn) {
	defer conn1.Close()
	defer conn2.Close()
	defer fmt.Printf("Connections closed.\n")

	fmt.Printf("Chat initiated.\n")

	scan1 := bufio.NewScanner(conn1)
	scan2 := bufio.NewScanner(conn2)
	go scanAndSend(*scan1, conn2)
	scanAndSend(*scan2, conn1)
}

func scanAndSend(scanner bufio.Scanner, conn net.Conn) {
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("read from 1: %s\n", text)
		writer := bufio.NewWriter(conn)
		_, err := writer.WriteString(text + "\n")
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		writer.Flush()
	}
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", Port)
	if err != nil {
		log.Fatalf("Error resolving address: %v", err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Error with listener: %v", err)
	}
	fmt.Printf("Listening on %s\n", Port)
	for {
		conn1, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Printf("Connection accepted (1/2)\n")
		conn2, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Printf("Connection accepted (2/2)\n")
		go handleConnections(conn1, conn2)
	}
}
