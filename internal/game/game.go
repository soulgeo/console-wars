package game

import (
	"fmt"
	"log"

	"github.com/soulgeo/console-wars/internal/messages"
)

func playGame(p1, p2 *Player, c chan string) {
	defer close(c)
	p1.initialize()
	p2.initialize()
	for p1.Health > 0 && p2.Health > 0 {
		playTurn(p1, p2, c)
		c <- fmt.Sprintf(
			messages.CurrentHealth,
			p1.Name,
			p1.Health,
			p2.Name,
			p2.Health,
		)
	}
	if p1.Health > 0 {
		c <- fmt.Sprintf(messages.Victory, p2.Name, p1.Name)
		return
	}
	if p2.Health > 0 {
		c <- fmt.Sprintf(messages.Victory, p1.Name, p2.Name)
		return
	}
	c <- fmt.Sprintf(messages.Tie)
}

func playTurn(p1, p2 *Player, c chan string) {
	processPreparation(p1, c)
	processPreparation(p2, c)

	if p1.Action == Attack {
		p1.attack(p2, c)
	}
	if p2.Action == Attack {
		p2.attack(p1, c)
	}
}

func processPreparation(p *Player, c chan string) {
	switch p.Action {
	case Defend:
		p.defend(c)
	case Charge:
		p.charge(c)
	case Dodge:
		p.dodge(c)
	case Heal:
		p.heal(c)
	case Attack:
		return
	default:
		log.Printf("Warning: Unknown action %s for player %s", p.Action, p.Name)
	}
}
