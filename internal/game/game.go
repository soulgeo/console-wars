package game

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/soulgeo/console-wars/internal/messages"
)

func Play(p1, p2 *Player, msg, act1, act2 chan string) {
	defer close(msg)
	// act1 and act2 are closed by the sender (server/receiveActions)

	p1.initialize()
	p2.initialize()

	msg <- fmt.Sprintf(messages.Welcome)
	msg <- fmt.Sprintf(messages.GameStart)

	turn := 0
	for p1.Health > 0 && p2.Health > 0 {
		turn++
		msg <- fmt.Sprintf(messages.NewTurn, turn)
		msg <- fmt.Sprintf(messages.CurrentHealth, p1.Name, p1.Health, p2.Name, p2.Health)
		msg <- fmt.Sprintf("%s\n", AwaitingInput)

		var err error
		p1.Action, p2.Action, err = waitForActions(act1, act2)
		if err != nil {
			msg <- err.Error()
			return
		}

		playTurn(p1, p2, msg)
	}
	if p1.Health > 0 {
		msg <- fmt.Sprintf(messages.Victory, p2.Name, p1.Name)
		return
	}
	if p2.Health > 0 {
		msg <- fmt.Sprintf(messages.Victory, p1.Name, p2.Name)
		return
	}
	msg <- fmt.Sprintf(messages.Tie)
}

func waitForActions(
	act1, act2 chan string,
) (string, string, error) {
	timeout := time.After(30 * time.Second)
	var a1, a2 string
	ch1, ch2 := act1, act2

	for ch1 != nil || ch2 != nil {
		select {
		case action, ok := <-ch1:
			if !ok {
				return "", "", fmt.Errorf(messages.PlayerDisconnected)
			}
			a1 = action
			ch1 = nil
		case action, ok := <-ch2:
			if !ok {
				return "", "", fmt.Errorf(messages.PlayerDisconnected)
			}
			a2 = action
			ch2 = nil
		case <-timeout:
			return "", "", errors.New(messages.GameTimeout)
		}
	}
	return a1, a2, nil
}

func playTurn(p1, p2 *Player, msg chan string) {
	processPreparation(p1, msg)
	processPreparation(p2, msg)

	if p1.Action == Attack {
		p1.attack(p2, msg)
	}
	if p2.Action == Attack {
		p2.attack(p1, msg)
	}
}

func processPreparation(p *Player, msg chan string) {
	switch p.Action {
	case Defend:
		p.defend(msg)
	case Charge:
		p.charge(msg)
	case Dodge:
		p.dodge(msg)
	case Heal:
		p.heal(msg)
	case Attack:
		return
	default:
		log.Printf("Warning: Unknown action %s for player %s", p.Action, p.Name)
	}
}
