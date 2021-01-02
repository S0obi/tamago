package tamagotchi

import (
	"fmt"
	"math/rand"
	"tamago/pkg/food"
	"tamago/pkg/status"
	"time"
)

// Tamagotchi : Tamagotchi representation
type Tamagotchi struct {
	Name     string
	Hunger   int
	Fatigue  int
	Hapiness int
	State    status.Status
}

const (
	// HapinessThreshold : Min level of hapiness to be happy
	HapinessThreshold = 25
	// FatigueThreshold : Max level of fatigue before sleeping
	FatigueThreshold = 90
)

// NewTamagotchi : Constructor of Tamagotchi struct
func NewTamagotchi(name string) *Tamagotchi {
	return &Tamagotchi{Name: name, Hunger: 50, Fatigue: 50, Hapiness: 50}
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

		if rand.Intn(10) == 2 {
			tamago.Hapiness += increase(tamago.Hapiness, 10)
		} else {
			tamago.Hapiness -= decrease(tamago.Hapiness, 1)
		}

		if tamago.Fatigue > FatigueThreshold {
			tamago.State = status.Sleeping
		} else if tamago.Hapiness <= HapinessThreshold {
			tamago.State = status.Sad
		} else {
			tamago.State = status.Happy
		}

		time.Sleep(1 * time.Second)
	}
	tamago.State = status.Dead
}

// IsAlive : return a boolean to know if Tamagotchi is still alive
func (tamago *Tamagotchi) IsAlive() bool {
	if tamago.Hunger >= 100 {
		return false
	}
	return true
}

// Feed : give food to Tamagotchi
func (tamago *Tamagotchi) Feed(yummy food.Food) {
	if yummy == food.Candy {
		tamago.Hunger -= decrease(tamago.Hunger, 15)
		tamago.Hapiness += increase(tamago.Hapiness, 25)
		tamago.Fatigue += increase(tamago.Fatigue, 5)
	} else if yummy == food.Meat {
		tamago.Hunger -= decrease(tamago.Hunger, 50)
		tamago.Fatigue += increase(tamago.Fatigue, 10)
	}
}

// Bed : tamagotchi will go to bed
func (tamago *Tamagotchi) Bed() {
	tamago.Fatigue = 0
}

func decrease(base int, value int) int {
	if base-value <= 0 {
		return base
	}
	return value
}

func increase(base int, value int) int {
	if base+value > 100 {
		return 0
	}
	return value
}
