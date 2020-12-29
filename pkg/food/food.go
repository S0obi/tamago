package food

// Food : Tamagotchi food
type Food int

const (
	// Meat : very nourishing
	Meat Food = iota
	// Candy : not very nourishing, but increase hapiness
	Candy
)

func (d Food) String() string {
	return [...]string{"Meat", "Candy"}[d]
}
