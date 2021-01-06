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
	Life     int
	Dirty    bool
	State    status.Status
}

const (
	// HapinessThreshold : Min level of hapiness to be happy
	HapinessThreshold = 25
	// FatigueThreshold : Max level of fatigue before sleeping
	FatigueThreshold = 90
	// HungerThreshold : Min level of hungryness
	HungerThreshold = 50
	// StarvingThreshold : Min level of hungryness before starving
	StarvingThreshold = 75
	// GameTick : Tamagotchi game tick in seconds
	GameTick = 2
)

// NewTamagotchi : Constructor of Tamagotchi struct
func NewTamagotchi(name string) *Tamagotchi {
	return &Tamagotchi{Name: name, Hunger: 25, Fatigue: 50, Hapiness: 50, Life: 100, Dirty: false}
}

// PrintStatus : Print the Tamagotchi status
func (tamago *Tamagotchi) PrintStatus() {
	fmt.Printf("[%s:%d:%s] hunger: %d, fatigue: %d, hapiness: %d, dirty: %t\n", tamago.Name, tamago.Life, tamago.State.String(), tamago.Hunger, tamago.Fatigue, tamago.Hapiness, tamago.Dirty)
}

// Live : Main Tamagotchi life loop
func (tamago *Tamagotchi) Live() {
	for tamago.IsAlive() {
		if tamago.State != status.Paused {
			tamago.Hunger += increase(tamago.Hunger, 1)
			tamago.Fatigue += increase(tamago.Fatigue, 1)
			tamago.Hapiness = tamago.drawHapiness()

			// Tamagotchi will loose hapiness if he is dirty
			if tamago.Dirty {
				tamago.Hapiness -= decrease(tamago.Hapiness, 10)
			} else {
				tamago.Dirty = tamago.drawDirty()
			}

			// Tamagotchi will loose life points if he is starving
			if tamago.Hunger > StarvingThreshold {
				tamago.Life -= decrease(tamago.Life, 1)
			} else if tamago.Hunger >= 100 {
				tamago.Life -= decrease(tamago.Life, 5)
			}

			if tamago.State != status.Sick {
				if tamago.Fatigue > FatigueThreshold {
					tamago.State = status.Sleeping
				} else if tamago.drawSickness() {
					tamago.State = status.Sick
				} else if tamago.Hunger > StarvingThreshold {
					tamago.State = status.Starving
				} else if tamago.Hunger > HungerThreshold {
					tamago.State = status.Hungry
				} else if tamago.Hapiness <= HapinessThreshold {
					tamago.State = status.Sad
				} else {
					tamago.State = status.Happy
				}
			} else {
				tamago.Life -= decrease(tamago.Life, 5)
			}
		}

		time.Sleep(GameTick * time.Second)
	}
	tamago.State = status.Dead
}

// IsAlive : return a boolean to know if Tamagotchi is still alive
func (tamago *Tamagotchi) IsAlive() bool {
	if tamago.Life <= 0 {
		return false
	}
	return true
}

// Feed : give food to Tamagotchi
func (tamago *Tamagotchi) Feed(yummy food.Food) {
	if yummy == food.Candy {
		tamago.Hunger -= decrease(tamago.Hunger, 5)
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

// Heal : tamagotchi will be healed
func (tamago *Tamagotchi) Heal() {
	tamago.State = status.Happy
	if rand.Intn(4) == 2 {
		tamago.Life += increase(tamago.Life, 5)
	}
	tamago.Fatigue -= decrease(tamago.Fatigue, 25)
}

// Clean : tamagotchi will be cleaned
func (tamago *Tamagotchi) Clean() {
	tamago.Dirty = false
	tamago.Hapiness += increase(tamago.Hapiness, 10)
}

func (tamago *Tamagotchi) drawHapiness() int {
	if rand.Intn(10) == 2 {
		return tamago.Hapiness + increase(tamago.Hapiness, 10)
	}
	return tamago.Hapiness - decrease(tamago.Hapiness, 1)
}

func (tamago *Tamagotchi) drawSickness() bool {
	probabilityRange := 100
	if tamago.Hunger > HungerThreshold {
		probabilityRange -= 10
	}

	if tamago.Hapiness < 50 {
		probabilityRange -= 10
	}

	if tamago.Fatigue > 50 {
		probabilityRange -= 20
	}

	if tamago.Dirty {
		probabilityRange -= 25
	}

	if rand.Intn(probabilityRange) == 2 {
		return true
	}
	return false
}

func (tamago *Tamagotchi) drawDirty() bool {
	if tamago.Hunger < 10 {
		if rand.Intn(20) == 2 {
			return true
		}
	} else {
		if rand.Intn(40) == 2 {
			return true
		}
	}
	return false
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
