package main

import (
	"fmt"
	"tamago/pkg/food"
	"tamago/pkg/tamagotchi"

	"github.com/manifoldco/promptui"
)

func main() {
	tamago := tamagotchi.NewTamagotchi("Tama")

	go tamago.Live()

	prompt := promptui.Select{
		Label: "Select an action",
		Items: []string{"Feed (meat)", "Feed (candy)", "Bed"},
	}

	for tamago.IsAlive() {
		_, action, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if action == "Feed (meat)" {
			tamago.Feed(food.Meat)
		} else if action == "Feed (candy)" {
			tamago.Feed(food.Candy)
		} else if action == "Bed" {
			tamago.Bed()
		}

		tamago.PrintStatus()
	}

	fmt.Println("So sad ... Your Tamagotchi is dead!")
}
