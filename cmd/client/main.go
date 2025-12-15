package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const Port = ":4567"

func writeToServer(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')

		writer := bufio.NewWriter(conn)
		_, err := writer.WriteString(text)
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		writer.Flush()
	}
}

func readFromServer(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Printf("Server wrote: %s\n", scanner.Text())
	}
}

func main() {
	conn, err := net.Dial("tcp", Port)
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	go readFromServer(conn)
	writeToServer(conn)
}
