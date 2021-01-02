package game

import (
	"tamago/pkg/food"
	"tamago/pkg/status"
	"tamago/pkg/tamagotchi"
	"tamago/pkg/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	// display
	screenWidth  = 320
	screenHeight = 240
	defaultTPS   = 60

	// timing
	feedingAnimationLength  = defaultTPS * 2
	sleepingAnimationLength = defaultTPS * 5
)

var (
	tamago        *tamagotchi.Tamagotchi
	runnerImage   *ebiten.Image
	happyImage    *ebiten.Image
	deadImage     *ebiten.Image
	feedingImage  *ebiten.Image
	sleepingImage *ebiten.Image
	sadImage      *ebiten.Image
)

// Game : ebiten game structure
type Game struct {
	count       int
	state       status.Status
	Tamago      *tamagotchi.Tamagotchi
	AudioPlayer *audio.Player
}

// Init : init method
func (g *Game) Init() {
	g.Tamago = tamagotchi.NewTamagotchi("Tama")

	// images
	happyImage = utils.NewImageFromFile("assets/happy.png")
	deadImage = utils.NewImageFromFile("assets/dead.png")
	feedingImage = utils.NewImageFromFile("assets/miam.png")
	sleepingImage = utils.NewImageFromFile("assets/dodo.png")
	sadImage = utils.NewImageFromFile("assets/sad.png")

	// audio
	g.AudioPlayer = utils.NewMusicFromFile("assets/music/theme.mp3")
}

// Update : ebiten update method
func (g *Game) Update() error {
	if g.Tamago.IsAlive() {
		if g.count > 0 {
			// Decrement the animation time
			g.count--
		} else {
			// default
			if g.Tamago.Hapiness > tamagotchi.HapinessThreshold {
				g.state = status.Happy
			} else {
				g.state = status.Sad
			}

			// feeding
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.feedTamago()
			}

			// sleeping
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
				g.putTamagoInBed()
			}
		}
	} else {
		g.state = status.Dead
	}

	return nil
}

// Draw : ebiten draw method
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	if g.Tamago.IsAlive() {
		if g.state == status.Dead {
			screen.DrawImage(deadImage, op)
		} else if g.state == status.Happy {
			screen.DrawImage(happyImage, op)
		} else if g.state == status.Feeding {
			screen.DrawImage(feedingImage, op)
		} else if g.state == status.Sleeping {
			screen.DrawImage(sleepingImage, op)
		} else if g.state == status.Sad {
			screen.DrawImage(sadImage, op)
		}
	} else {
		screen.DrawImage(deadImage, op)
	}
}

// Layout : ebiten layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) feedTamago() {
	g.count = feedingAnimationLength
	g.state = status.Feeding
	g.Tamago.Feed(food.Meat)
}

func (g *Game) putTamagoInBed() {
	g.count = sleepingAnimationLength
	g.state = status.Sleeping
	g.Tamago.Bed()
}
