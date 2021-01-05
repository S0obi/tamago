package status

// Status : Tamagotchi status
type Status int

const (
	// Dead : tamagoshi is dead
	Dead Status = iota
	// Paused : tamagotchi is paused
	Paused
	// Happy : tamagoshi is happy
	Happy
	// Feeding : tamagoshi is being feed
	Feeding
	// Sleeping : tamagoshi is sleeping
	Sleeping
	// Sad : tamagoshi is sad
	Sad
	// Sick : tamagotchi is sick
	Sick
	// Hungry : tamagotchi is hungry
	Hungry
	// Starving : tamagotchi is starving
	Starving
)

func (s Status) String() string {
	return [...]string{"Dead", "Paused", "Happy", "Feeding", "Sleeping", "Sad", "Sick", "Hungry", "Starving"}[s]
}
