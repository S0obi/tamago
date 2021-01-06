package actions

import "container/ring"

// Action : Tamagotchi action
type Action int

const (
	// Sleep : Put Tamagotchi in bed
	Sleep = "sleep"
	// Candy : Give a candy
	Candy = "candy"
	// Feed : Feed Tamagotchi
	Feed = "feed"
	// Heal : Heal Tamagotchi
	Heal = "heal"
	// Clean : Clean poop
	Clean = "clean"
)

var (
	actions = [...]string{Sleep, Candy, Feed, Heal, Clean}
)

// NewTamagoActions : Constructor of ActionList
func NewTamagoActions() *ring.Ring {
	r := ring.New(len(actions))
	for i := 0; i < r.Len(); i++ {
		r.Value = actions[i]
		r = r.Next()
	}
	return r
}
