package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/soulgeo/console-wars/internal/config"
	"github.com/soulgeo/console-wars/internal/game"
	"github.com/soulgeo/console-wars/internal/messages"
)

func readFromServer(conn net.Conn, inReader *bufio.Reader) {
	connScanner := bufio.NewScanner(conn)
	for connScanner.Scan() {
		text := connScanner.Text()
		if text != game.AwaitingInput {
			fmt.Printf("%s\n", text)
			continue
		}
		fmt.Printf(messages.AwaitAction)
		writeToServer(conn, inReader)
	}
}

func writeToServer(conn net.Conn, inReader *bufio.Reader) {
	connWriter := bufio.NewWriter(conn)

	for {
		input, err := inReader.ReadString('\n')
		if err != nil {
			// EOF or other error, stop trying to read
			fmt.Printf("Input error: %v\n", err)
			os.Exit(0)
		}
		input = strings.TrimSpace(input)
		err = filterInput(input)
		if err != nil {
			fmt.Printf("Invalid input, try again.\n")
			continue
		}
		_, err = connWriter.WriteString(input + "\n")
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		connWriter.Flush()
		break
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
	inReader := bufio.NewReader(os.Stdin)
	readFromServer(conn, inReader)
}
