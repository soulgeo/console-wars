package game

import (
	"fmt"
	"log"

	"github.com/soulgeo/console-wars/internal/messages"
)

func playGame(p1, p2 *Player) {
	p1.initialize()
	p2.initialize()
	for p1.Health > 0 && p2.Health > 0 {
		playTurn(p1, p2)
		fmt.Printf(
			messages.CurrentHealth,
			p1.Name,
			p1.Health,
			p2.Name,
			p2.Health,
		)
	}
	if p1.Health > 0 {
		fmt.Printf(messages.Victory, p2.Name, p1.Name)
		return
	}
	if p2.Health > 0 {
		fmt.Printf(messages.Victory, p1.Name, p2.Name)
		return
	}
	fmt.Printf(messages.Tie)
}

func playTurn(p1, p2 *Player) {
	switch p1.Action {
	case Attack:
		p1.attack(p2)
	case Defend:
		p1.defend()
	case Charge:
		p1.charge()
	case Dodge:
		p1.dodge()
	case Heal:
		p1.heal()
	default:
		log.Fatal("Error")
	}
	switch p2.Action {
	case Attack:
		p2.attack(p1)
	case Defend:
		p2.defend()
	case Charge:
		p2.charge()
	case Dodge:
		p2.dodge()
	case Heal:
		p2.heal()
	default:
		log.Fatal("Error")
	}
}
