// Package messages
package messages

const (
	Connected          = "Connected to the server.\n"
	Waiting            = "Waiting for another player to join...\n"
	MatchFound         = "Match found.\n"
	Welcome            = "Welcome to CONSOLE WARS!\n"
	GameStart          = "Ready? FIGHT!\n"
	NewTurn            = "----- TURN %d -----\n"
	CurrentHealth      = "Current Health:\n%s: %d HP\n%s: %d HP\n"
	AwaitAction        = "What will you do? > "
	AttackHit          = "%s attacks and deals %d damage to %s.\n"
	AttackMiss         = "%s attacks but misses.\n"
	Critical           = "Critical hit! "
	Defense            = "%s increases their defenses.\n"
	Charge             = "%s powers up.\n"
	Dodge              = "%s attempts to dodge.\n"
	Heal               = "%s heals for %d health. %d heals remaining.\n"
	HealFail           = "%s attempts to heal but is out of heals.\n"
	Victory            = "%s is dead. %s wins!\n"
	Tie                = "It's a tie!\n"
	PlayerDisconnected = "player disconnected"
	GameTimeout        = "game timed out due to inactivity"
)
