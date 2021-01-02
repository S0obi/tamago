package status

// Status : Tamagotchi status
type Status int

const (
	// Dead : tamagoshi is dead
	Dead Status = iota
	// Happy : tamagoshi is happy
	Happy
	// Feeding : tamagoshi is being feed
	Feeding
	// Sleeping : tamagoshi is sleeping
	Sleeping
	// Sad : tamagoshi is sad
	Sad
)

func (s Status) String() string {
	return [...]string{"Dead", "Happy", "Feeding", "Sleeping", "Sad"}[s]
}
