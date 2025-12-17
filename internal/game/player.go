// Package game
package game

import (
	"fmt"
	"math/rand"

	"github.com/soulgeo/console-wars/internal/messages"
)

type Player struct {
	Name    string
	Action  string
	Health  int
	Attack  float32
	Armor   int
	Dodging bool
}

const (
	MaxAttackRoll = 20
	MaxDamageRoll = 10
	MaxAttack     = 1.5
	MaxArmor      = 16
	StartHealth   = 100
	StartAttack   = 1.0
	StartArmor    = 10
)

func (p *Player) initialize() {
	p.Health = StartHealth
	p.Attack = StartAttack
	p.Armor = StartArmor
}

func (p *Player) attack(defender *Player) {
	attackRoll := rand.Intn(MaxAttackRoll-1) + 1
	if defender.Dodging {
		secondAttackRoll := rand.Intn(MaxAttackRoll-1) + 1
		attackRoll = min(attackRoll, secondAttackRoll)
	}
	if attackRoll <= defender.Armor {
		fmt.Printf(messages.AttackMiss, p.Name)
		return
	}
	damageRoll := rand.Intn(MaxDamageRoll-1) + 1
	var bonusDamageRoll int
	if attackRoll == 20 {
		fmt.Printf(messages.Critical)
		bonusDamageRoll = rand.Intn(MaxDamageRoll-1) + 1
	} else {
		bonusDamageRoll = 0
	}
	damage := float32(damageRoll+bonusDamageRoll) * p.Attack
	defender.Health -= int(damage)
	fmt.Printf(messages.AttackHit, p.Name, int(damage), defender.Name)
}

func (p *Player) defend() {
	if p.Armor < MaxArmor {
		p.Armor += 1
	}
	fmt.Printf(messages.Defense, p.Name)
}

func (p *Player) charge() {
	if p.Attack < MaxAttack {
		p.Attack += 0.1
	}
	fmt.Printf(messages.Charge, p.Name)
}

func (p *Player) dodge() {
	p.Dodging = true
	fmt.Printf(messages.Dodge, p.Name)
}
