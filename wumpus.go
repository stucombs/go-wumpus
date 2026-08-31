package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/nexidian/gocliselect"
	"github.com/ttacon/chalk"
)

type menuItem struct {
	text string
	value string
}

var cave = map[int][3]int{
    1: {2, 3, 4}, 2: {1, 5, 6}, 3: {1, 7, 8}, 4: {1, 9, 10}, 5: {2, 9, 11},
    6: {2, 7, 12}, 7: {3, 6, 13}, 8: {3, 10, 14}, 9: {4, 5, 15}, 10: {4, 8, 16},
    11: {5, 12, 17}, 12: {6, 11, 18}, 13: {7, 14, 18}, 14: {8, 13, 19},
    15: {9, 16, 17}, 16: {10, 15, 19}, 17: {11, 20, 15}, 18: {12, 13, 20},
    19: {14, 16, 20}, 20: {17, 18, 19},
}

var wumpus, bat1, bat2, pit1, pit2 int

type player struct {
	pos int
	alive bool
	arrows int
	isShooting bool
	isMoving bool
	winner bool
	exitMessage string
}

func (p player) status() {
	fmt.Printf("\nYou stand in cavern %v.\n", chalk.Bold.TextStyle(strconv.Itoa(p.pos)))
	fmt.Printf("You have %v arrow(s) left.\n\n", chalk.Bold.TextStyle(strconv.Itoa(p.arrows)))
}

func (p *player) move() {
	for p.isMoving {
		moves := cave[p.pos];
		moveMenu := gocliselect.NewMenu("Move to?")
		for _, move := range moves {
			value := strconv.Itoa(move)
			moveMenu.AddItem(value, value)
		}

		moveMenu.AddItem("Back", "b")

		choice := moveMenu.Display()

		if choice == "b" {
			p.isMoving = false
		} else {
			p.pos, _ = strconv.Atoi(choice)
			p.isMoving = false
		}	
	}
}

func (p *player) shoot() {
	for p.isShooting {
		moves := cave[p.pos]
		shootMenu := gocliselect.NewMenu("Shoot where?")
		for _, move := range moves {
			value := strconv.Itoa(move)
			shootMenu.AddItem(value, value)
		}

		shootMenu.AddItem("Back", "b")
		choice := shootMenu.Display()

		if choice == "b" {
			p.isShooting = false
		} else {
			fmt.Printf("You shoot into cavern %v\n", choice)
			shootInt, _ := strconv.Atoi(choice)
			if shootInt == wumpus {
				// Player wins
				p.winner = true
				p.alive = false
			} else {
				p.wakeWumpus()
				p.arrows--
				if p.arrows == 0 {
					fmt.Println(chalk.Italic.TextStyle("You shoot your last arrow."))
					p.alive = false
					p.exitMessage = "You ran out of arrows."
					
				}
			}
			p.isShooting = false
		}
	}
}

func (p player) sense() {
	for _, cavern := range cave[p.pos] {
		if cavern == wumpus {
			fmt.Println(chalk.Italic.TextStyle("You smell something terrible nearby...\n"))
		}

		if cavern == bat1 || cavern == bat2 {
			fmt.Println(chalk.Italic.TextStyle("You hear a rustling...\n"))
		}

		if cavern == pit1 || cavern == pit2 {
			fmt.Println(chalk.Italic.TextStyle("You feel a cold wind blowing from a nearby cavern...\n"))
		}
	}
}

func (p *player) isDead() bool {
// Determine if the player hit a game ending condition
	if p.pos == wumpus {
		p.exitMessage = "You enter a cavern with the Wumpus. You are eaten alive!"
		p.alive = false

		return true
	}

	if p.pos == pit1 || p.pos == pit2 {
		p.exitMessage = "You enter a cavern with a large, endless pit. You fall to your death!"
		p.alive = false

		return true
	}

	return false
}


func (p *player) isDanger() {
// Determine if the player entered a room with a bat
	if p.pos == bat1 || p.pos == bat2 {
		fmt.Println("\nYou enter a cavern and large bat picks you up and throws you into a new cavern!")
		p.randomRoom()
	}
}

func (p *player) randomRoom() {
// Place player in random, empty room
	isOccupied := true
	for isOccupied {
		p.pos = rand.Intn(19) + 1
		if p.pos == wumpus || p.pos == bat1 || p.pos == bat2 || p.pos == pit1 || p.pos == pit2 {
			continue
		} else {
			isOccupied = false
		}
	}
}

func (p *player) wakeWumpus() {
	wake := rand.Intn(4)
	if wake > 0 {
		wumpus = cave[wumpus][rand.Intn(3)]
		if wumpus == p.pos {
			p.exitMessage = "The Wumpus is startled and moves into your cave. You are eaten alive!"
			p.alive = false
		}
	}
}

func menuOptions() string {
	menuOptions := []menuItem {
		{ text: "Move",  value: "move" },
		{ text: "Shoot", value: "shoot" },
		{ text: "Quit",  value: "quit" },
	}

	menu := gocliselect.NewMenu("Select an option")
	
	for _, item := range menuOptions {
		menu.AddItem(item.text, item.value)	
	}

	choice := menu.Display()
	return choice
}

func clearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	// game state setup
	p := player {
		pos: 1,
		alive: true,
		winner: false,
		arrows: 5,
	}

	wumpus = rand.Intn(19) + 2 // caverns 2 - 20
	bat1 = rand.Intn(19) + 2
	bat2 = rand.Intn(19) + 2
	pit1 = rand.Intn(19) + 2
	pit2 = rand.Intn(19) + 2

	// begin game
	fmt.Println("You enter a cave. Hunt the Wumpus...")
	for p.alive {
		if p.isDead() {
			continue
		}
		p.isDanger()
		p.status()
		p.sense()
		// Player option select returned by menuOptions
		choice := menuOptions()

		switch choice {
		case "move":
			p.isMoving = true
			p.move()
			clearScreen()
		case "shoot":
			p.isShooting = true
			p.shoot()
			clearScreen()

		case "quit":
			clearScreen()
			p.exitMessage = "\nYou failed to hunt the wumpus."
			p.alive = false
		}
	}

	if p.winner {
		fmt.Println(chalk.Green.Color("\nWumpus hunted, you win!"))
	} else {
		fmt.Printf(chalk.Red.Color("\nPlayer died. %v\n\n"), p.exitMessage)
	}
}
