package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:4567")
	if err != nil {
		fmt.Printf("Error with listener.\n")
	}
	conn, err = listener.Accept()
	if err != nil {
		fmt.Printf("Error with connection.\n")
	}
}
