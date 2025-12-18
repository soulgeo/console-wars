package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/soulgeo/console-wars/internal/config"
	"github.com/soulgeo/console-wars/internal/game"
)

func readFromServer(conn net.Conn) {
	connScanner := bufio.NewScanner(conn)
	for connScanner.Scan() {
		text := connScanner.Text()
		if text != game.AwaitingInput {
			fmt.Printf("%s\n", text)
			continue
		}
		writeToServer(conn)
	}
}

func writeToServer(conn net.Conn) {
	for {
		connWriter := bufio.NewWriter(conn)
		inReader := bufio.NewReader(os.Stdin)
		input, _ := inReader.ReadString('\n')
		err := filterInput(input)
		if err != nil {
			fmt.Printf("Invalid input, try again.\n")
			continue
		}
		_, err = connWriter.WriteString(input)
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		connWriter.Flush()
	}
}

func filterInput(input string) error {
	if input != game.Attack && input != game.Defend && input != game.Charge &&
		input != game.Dodge &&
		input != game.Heal {
		return errors.New("invalid input")
	}
	return nil
}

func main() {
	conn, err := net.Dial("tcp", config.Port)
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	readFromServer(conn)
}
