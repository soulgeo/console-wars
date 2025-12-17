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
	Heals   int
	Dodging bool
}

const (
	MaxAttackRoll = 20
	MaxDamageRoll = 10
	MaxHealRoll   = 6
	MaxAttack     = 1.5
	MaxArmor      = 16
	StartHealth   = 30
	StartAttack   = 1.0
	StartArmor    = 10
	StartHeals    = 3
)

func (p *Player) initialize() {
	p.Health = StartHealth
	p.Attack = StartAttack
	p.Armor = StartArmor
	p.Heals = StartHeals
}

func (p *Player) attack(defender *Player, logChan chan string) {
	attackRoll := rand.Intn(MaxAttackRoll) + 1
	if defender.Dodging {
		secondAttackRoll := rand.Intn(MaxAttackRoll) + 1
		attackRoll = min(attackRoll, secondAttackRoll)
		defender.Dodging = false
	}
	if attackRoll <= defender.Armor {
		logChan <- fmt.Sprintf(messages.AttackMiss, p.Name)
		return
	}
	damageRoll := rand.Intn(MaxDamageRoll) + 1
	var bonusDamageRoll int
	if attackRoll == 20 {
		logChan <- fmt.Sprintf(messages.Critical)
		bonusDamageRoll = rand.Intn(MaxDamageRoll) + 1
	} else {
		bonusDamageRoll = 0
	}
	damage := float32(damageRoll+bonusDamageRoll) * p.Attack
	defender.Health -= int(damage)
	logChan <- fmt.Sprintf(messages.AttackHit, p.Name, int(damage), defender.Name)
}

func (p *Player) defend(logChan chan string) {
	if p.Armor < MaxArmor {
		p.Armor += 1
	}
	logChan <- fmt.Sprintf(messages.Defense, p.Name)
}

func (p *Player) charge(logChan chan string) {
	if p.Attack < MaxAttack {
		p.Attack += 0.1
	}
	logChan <- fmt.Sprintf(messages.Charge, p.Name)
}

func (p *Player) dodge(logChan chan string) {
	p.Dodging = true
	logChan <- fmt.Sprintf(messages.Dodge, p.Name)
}

func (p *Player) heal(logChan chan string) {
	if p.Heals == 0 {
		logChan <- fmt.Sprintf(messages.HealFail, p.Name)
		return
	}
	healRoll := rand.Intn(MaxHealRoll) + 1
	p.Health += healRoll
	p.Heals -= 1
	logChan <- fmt.Sprintf(messages.Heal, p.Name, healRoll, p.Heals)
}
