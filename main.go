package main

import (
	"fmt"
	"math/rand"
	"os"
)

var cave = map[int][3]int{
    1: {2, 3, 4}, 2: {1, 5, 6}, 3: {1, 7, 8}, 4: {1, 9, 10}, 5: {2, 9, 11},
    6: {2, 7, 12}, 7: {3, 6, 13}, 8: {3, 10, 14}, 9: {4, 5, 15}, 10: {4, 8, 16},
    11: {5, 12, 17}, 12: {6, 11, 18}, 13: {7, 14, 18}, 14: {8, 13, 19},
    15: {9, 16, 17}, 16: {10, 15, 19}, 17: {11, 20, 15}, 18: {12, 13, 20},
    19: {14, 16, 20}, 20: {17, 18, 19},
}

var wumpus, player, bat1, bat2, pit1, pit2 int

var arrows = 5

func sense() {
// Nearby Wumpus: "You smell something terrible nearby."
// Nearby bat: "You hear a rustling."
// Nearby pit: "You feel a cold wind blowing from a nearby cavern."
}

func printMenu() {
	gameStatus()
	menuOptions := []string{
		"i) Show status",
		"m) Move",
		"s) Shoot",
		"q) Quit",
	}
	
	for _, item := range menuOptions {
		fmt.Println(item)	
	}
	fmt.Println("Select an option:")

	var input string
	fmt.Scanln(&input)
}

func gameStatus() {
	fmt.Printf("You stand in room %v.\n", player)
	fmt.Printf("You have %v arrows left.\n", arrows)
}

func main() {
	// game state variables
	playerAlive := true
	gameWon := false
	var endMessage string

	// object locations
	player = 1
	wumpus = rand.Intn(len(cave))

	// begin game
	for playerAlive {
		fmt.Println("You enter a cave. Hunt the Wumpus...")
		printMenu()

		playerAlive = false
	}

	if gameWon {
		fmt.Println("Wumpus hunted, you win!")
		os.Exit(1)
	}

	fmt.Printf("Player died. %v\n", endMessage)
}
