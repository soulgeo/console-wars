# 🎮 Console Wars

A turn-based PvP battle game you can play over TCP from the terminal, written in Go. Run the game server and connect two clients to start a game.
Multiple games can run at the same time independent of each other.

## Installation

To install and run the server and/or the client locally, you need to have Go installed (version 1.25.5 or later).

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/soulgeo/console-wars.git
    cd console-wars
    ```

    (Replace the URL with the actual repository URL if different.)

2.  **Build the executables:**

    ```bash
    go build -o bin/console-wars ./cmd/console-wars
    go build -o bin/console-wars-server ./cmd/console-wars-server
    ```

    This will create two executable files in the `bin/` directory: `console-wars` (the client) and `console-wars-server` (the server).

## Usage

### Run the Server

To start the game server:

```bash
./bin/console-wars-server
```

The server will listen for incoming client connections on port `:4567`.

### Run the Client

To start the game client:

```bash
./bin/console-wars
```

This will connect to the default server address, which is `localhost:4567`.

If the server is running on a different address or port, you can specify it using the `-server` flag:

```bash
./bin/console-wars -server=<server_address>:<port>
```

For example:

```bash
./bin/console-wars -server=my-remote-server.com:4567
```

The client will prompt you for your player name and then connect to the server to join a game.

## How to Play

The game starts with both players having 30 HP. If one player's HP goes to 0, the other player wins.

Each turn, each player can take one of 5 different actions, by typing its name in lowercase on the console when prompted.

### Actions

- 🗡️ **Attack:** Attempt to deal up to 10 damage to your opponent, reducing their HP by that amount. There's a chance you might miss,
  which is increased the more armor your opponent has. There's a 1/20 chance you will critically hit, which can have you deal
  up to 20 damage to your opponent.

- 🛡️ **Defend:** Increase your armor, which decreases the chance that your opponent's attacks will go through.

- 💪 **Charge:** Increase your attack, which increases the damage your attacks will deal.

- 💨 **Dodge:** If your opponent takes the Attack action this turn, they will have a significantly lower chance of hitting you.

- ✨ **Heal:** Heal for up to 6 HP. You can only do this 3 times during the game.

Once both players select their action, they will happen simultaneously.

If both players fall bellow 0 HP in the same round, the result is a tie.
