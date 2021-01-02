package game

import (
	"image"
	_ "image/png"
	"log"
	"os"
	"tamago/pkg/food"
	"tamago/pkg/tamagotchi"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240
	defaultTPS   = 60
)

var (
	tamago        *tamagotchi.Tamagotchi
	runnerImage   *ebiten.Image
	happyImage    *ebiten.Image
	deadImage     *ebiten.Image
	feedingImage  *ebiten.Image
	sleepingImage *ebiten.Image
)

// Game : ebiten game structure
type Game struct {
	count  int
	state  string
	Tamago *tamagotchi.Tamagotchi
}

func newImageFromFile(path string) *ebiten.Image {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}

// Init : init method
func (g *Game) Init() {
	g.Tamago = tamagotchi.NewTamagotchi("Tama")

	go g.Tamago.Live()

	happyImage = newImageFromFile("assets/happy.png")
	deadImage = newImageFromFile("assets/dead.png")
	feedingImage = newImageFromFile("assets/miam.png")
	sleepingImage = newImageFromFile("assets/dodo.png")
}

// Update : ebiten update method
func (g *Game) Update() error {
	if g.Tamago.IsAlive() {

		// feeding
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.count = defaultTPS * 2
			g.state = "feeding"
			g.Tamago.Feed(food.Meat)
		}

		// sleeping
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			g.count = defaultTPS * 5
			g.state = "sleeping"
			g.Tamago.Bed()
		}

		if g.count > 0 {
			g.count--
		} else {
			g.state = "happy"
		}

	} else {
		g.state = "dead"
	}
	return nil
}

// Draw : ebiten draw method
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	if g.Tamago.IsAlive() {
		if g.state == "dead" {
			screen.DrawImage(deadImage, op)
		} else if g.state == "happy" {
			screen.DrawImage(happyImage, op)
		} else if g.state == "feeding" {
			screen.DrawImage(feedingImage, op)
		} else if g.state == "sleeping" {
			screen.DrawImage(sleepingImage, op)
		}
	} else {
		screen.DrawImage(deadImage, op)
	}
}

// Layout : ebiten layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
