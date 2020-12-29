package tamagotchi

import (
	"fmt"
	"math/rand"
	"tamago/pkg/food"
	"time"
)

// Tamagotchi : Tamagotchi representation
type Tamagotchi struct {
	Name     string
	Hunger   int
	Fatigue  int
	Hapiness int
}

// NewTamagotchi : Constructor of Tamagotchi struct
func NewTamagotchi(name string) *Tamagotchi {
	return &Tamagotchi{Name: name, Hunger: 80, Fatigue: 0, Hapiness: 100}
}

// PrintStatus : Print the Tamagotchi status
func (tamago *Tamagotchi) PrintStatus() {
	fmt.Printf("[%s] hunger: %d, fatigue: %d, hapiness: %d\n", tamago.Name, tamago.Hunger, tamago.Fatigue, tamago.Hapiness)
}

// Live : Main Tamagotchi life loop
func (tamago *Tamagotchi) Live() {
	for tamago.IsAlive() {
		tamago.Hunger += increase(tamago.Hunger, 1)
		tamago.Fatigue += increase(tamago.Fatigue, 1)
		if rand.Intn(4) == 2 {
			tamago.Hapiness += increase(tamago.Hapiness, 1)
		}
		time.Sleep(1 * time.Second)
	}
}

// IsAlive : return a boolean to know if Tamagotchi is still alive
func (tamago *Tamagotchi) IsAlive() bool {
	if tamago.Hunger >= 100 {
		return false
	} else {
		return true
	}
}

// Feed : give food to Tamagotchi
func (tamago *Tamagotchi) Feed(yummy food.Food) {
	if yummy == food.Candy {
		tamago.Hunger -= decrease(tamago.Hunger, 25)
		tamago.Hapiness += increase(tamago.Hapiness, 10)
	} else if yummy == food.Meat {
		tamago.Hunger -= decrease(tamago.Hunger, 50)
	}
}

// Bed : tamagotchi will go to bed
func (tamago *Tamagotchi) Bed() {
	tamago.Fatigue = 0
}

func decrease(base int, value int) int {
	if base-value <= 0 {
		return base
	} else {
		return value
	}
}

func increase(base int, value int) int {
	if base+value >= 100 {
		return 0
	} else {
		return value
	}
}
