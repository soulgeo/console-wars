package game

import (
	"fmt"
	"log"

	"github.com/soulgeo/console-wars/internal/messages"
)

func PlayGame(p1, p2 *Player, msg, act1, act2 chan string) {
	defer close(msg)
	defer close(act1)
	defer close(act2)
	p1.initialize()
	p2.initialize()
	msg <- fmt.Sprintf(messages.GameStart)
	turn := 0
	for p1.Health > 0 && p2.Health > 0 {
		turn++
		msg <- fmt.Sprintf(messages.NewTurn, turn)
		msg <- fmt.Sprintf(messages.CurrentHealth, p1.Name, p1.Health, p2.Name, p2.Health)
		msg <- fmt.Sprintf("%s\n", AwaitingInput)

		p1.Action, p2.Action = <-act1, <-act2

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
