package main

import (
	"fmt"
	"math/rand"
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
	winner bool
	exitMessage string
}

func (p player ) status() {
	fmt.Printf("\nYou stand in cavern %v.\n", chalk.Bold.TextStyle(strconv.Itoa(p.pos)))
	fmt.Printf("You have %v arrow(s) left.\n\n", chalk.Bold.TextStyle(strconv.Itoa(p.arrows)))
}

func (p player) move() string {
	moves := cave[p.pos];
	moveMenu := gocliselect.NewMenu("Move to?")
	for _, move := range moves {
		value := strconv.Itoa(move)
		moveMenu.AddItem(value, value)
	}

	moveMenu.AddItem("Back", "b")

	choice := moveMenu.Display()

	return choice
}

func (p player) shoot() string {
	moves := cave[p.pos]
	shootMenu := gocliselect.NewMenu("Shoot where?")
	for _, move := range moves {
		value := strconv.Itoa(move)
		shootMenu.AddItem(value, value)
	}

	shootMenu.AddItem("Back", "b")
	
	choice := shootMenu.Display()
	return choice
}

func (p player) sense() {
	for _, cavern := range cave[p.pos] {
		if cavern == wumpus {
			fmt.Println(chalk.Italic.TextStyle("You smell something terrible nearby."))
		}

		if cavern == bat1 || cavern == bat2 {
			fmt.Println(chalk.Italic.TextStyle("You hear a rustling."))
		}

		if cavern == pit1 || cavern == pit2 {
			fmt.Println(chalk.Italic.TextStyle("You feel a cold wind blowing from a nearby cavern."))
		}
	}
}

func menuOptions() string {
	menuOptions := []menuItem {
		{ text: "Move",  value: "move" },
		{ text: "Shoot", value: "shoot" },
		{ text: "Show Board", value: "game"},
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
	//TODO
}

func showBoard() {
	// TEMP DEBUG
	fmt.Printf("Wumpus pos %v\n", wumpus)
	fmt.Printf("Bat 1 pos: %v & Bat 2 pos: %v\n", bat1, bat2)
	fmt.Printf("Pit 1 pos: %v & Pit 2 pos: %v\n", pit1, pit2)
}

func main() {
	// game state setup
	p := player {
		pos: 1,
		alive: true,
		winner: false,
		arrows: 1,
	}

	wumpus = rand.Intn(len(cave)) //TODO: How to make these not be in room 1?
	bat1 = rand.Intn(len(cave))
	bat2 = rand.Intn(len(cave))
	pit1 = rand.Intn(len(cave))
	pit2 = rand.Intn(len(cave))

	// begin game
	fmt.Println("You enter a cave. Hunt the Wumpus...")
	for p.alive {
		p.status()
		// Player option select returned by menuOptions
		choice := menuOptions()

		switch choice {
		case "move":
			moveChoice := p.move()
			if moveChoice == "b" {
				// Reset to main menu
				continue
			}

			fmt.Println("\nYou move to cavern ", moveChoice)
			intMoveChoice, _ := strconv.Atoi(moveChoice)
			p.pos = intMoveChoice
			p.sense()
		case "shoot":
			shootChoice := p.shoot()
			if shootChoice == "b" {
				// Reset to main menu
				continue
			}

			fmt.Printf("\nYou shoot into cavern %v\n", shootChoice)
			intShootChoice, _ := strconv.Atoi(shootChoice)
			if intShootChoice == wumpus {
				p.winner = true
				p.alive = false
			}

			p.arrows -= 1
			if p.arrows == 0 {
				// TODO: (optional) Add "run away" state?
				fmt.Println(chalk.Italic.TextStyle("You shoot your last arrow."))
				p.alive = false
				p.exitMessage = "You ran out of arrows."
			}

			// TOOD: On miss, .75 chance for wumpus to wake and move
		case "game":
			showBoard()
		case "quit":
			p.exitMessage = "\nYou failed to hunt the wumpus."
			p.alive = false
		}
	}

	if p.winner {
		fmt.Println(chalk.Green.Color("Wumpus hunted, you win!"))
	} else {
		fmt.Printf(chalk.Red.Color("\nPlayer died ☠️. %v\n\n"), p.exitMessage)
	}
}
