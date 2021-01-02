package game

import (
	"container/ring"
	"tamago/pkg/actions"
	"tamago/pkg/food"
	"tamago/pkg/status"
	"tamago/pkg/tamagotchi"
	"tamago/pkg/utils"
	"time"

	_ "image/png" // need this for png support

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	// audio
	sampleRate = 22050

	// display
	screenWidth  = 330
	screenHeight = 330
	defaultTPS   = 60

	// timing
	feedingAnimationLength  = defaultTPS * 2
	sleepingAnimationLength = defaultTPS * 5
)

var (
	// animation images
	happyImage    *ebiten.Image
	deadImage     *ebiten.Image
	feedingImage  *ebiten.Image
	sleepingImage *ebiten.Image
	sadImage      *ebiten.Image
	sickImage     *ebiten.Image

	// action image
	actionFeedImage  *ebiten.Image
	actionCandyImage *ebiten.Image
	actionSleepImage *ebiten.Image
	actionsHealImage *ebiten.Image
)

// Game : ebiten game structure
type Game struct {
	count                  int
	Tamago                 *tamagotchi.Tamagotchi
	audioPlayer            *audio.Player
	audioContext           *audio.Context
	currentAnimation       status.Status
	currentAction          *ring.Ring
	currentActionID        int
	CurrentMusic           chan string
	deadMusicAlreadyPlayed bool
	muteMusic              bool
}

// Init : init method
func (g *Game) Init() {
	g.Tamago = tamagotchi.NewTamagotchi("Tama")

	g.Tamago.State = status.Happy

	// animation images
	happyImage, _, _ = ebitenutil.NewImageFromFile("assets/happy.png")
	deadImage, _, _ = ebitenutil.NewImageFromFile("assets/dead.png")
	feedingImage, _, _ = ebitenutil.NewImageFromFile("assets/miam.png")
	sleepingImage, _, _ = ebitenutil.NewImageFromFile("assets/dodo.png")
	sadImage, _, _ = ebitenutil.NewImageFromFile("assets/sad.png")
	sickImage, _, _ = ebitenutil.NewImageFromFile("assets/sick.png")

	// action bar images
	actionFeedImage, _, _ = ebitenutil.NewImageFromFile("assets/actions/feed.png")
	actionCandyImage, _, _ = ebitenutil.NewImageFromFile("assets/actions/candy.png")
	actionSleepImage, _, _ = ebitenutil.NewImageFromFile("assets/actions/sleep.png")
	actionsHealImage, _, _ = ebitenutil.NewImageFromFile("assets/actions/heal.png")

	g.currentAction = actions.NewTamagoActions()

	// audio
	g.CurrentMusic = make(chan string)
	g.audioContext = audio.NewContext(sampleRate)
	g.audioPlayer = utils.NewMusicFromFile("assets/music/theme.mp3", g.audioContext)
	g.deadMusicAlreadyPlayed = false
	g.muteMusic = false
}

func (g *Game) inifiniteThemeMusic() {
	for {
		if !g.audioPlayer.IsPlaying() && g.Tamago.State != status.Dead && !g.muteMusic {
			g.CurrentMusic <- "theme"
		}
		time.Sleep(2 * time.Second)
	}
}

// PlayMusic : play theme music
func (g *Game) PlayMusic() {
	for {
		go g.inifiniteThemeMusic()
		music := <-g.CurrentMusic
		g.audioPlayer.Close()
		if music == "theme" || music == "" {
			g.audioPlayer = utils.NewMusicFromFile("assets/music/theme.mp3", g.audioContext)
		} else if music == "dead" && !g.deadMusicAlreadyPlayed {
			g.audioPlayer = utils.NewMusicFromFile("assets/music/dead.mp3", g.audioContext)
			g.deadMusicAlreadyPlayed = true
		}
		g.audioPlayer.Play()
	}
}

// Update : ebiten update method
func (g *Game) Update() error {
	if g.Tamago.IsAlive() {
		// Mute the theme music
		if inpututil.IsKeyJustPressed(ebiten.KeyM) {
			g.muteMusic = !g.muteMusic
			if g.muteMusic {
				g.audioPlayer.Close()
			} else {
				g.CurrentMusic <- "theme"
			}
		}

		if g.currentAnimation != status.Sleeping {
			// Decrease action ID
			if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
				g.currentAction = g.currentAction.Prev()
			}

			// Increase action ID
			if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
				g.currentAction = g.currentAction.Next()
			}

			// Perform current action
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				if g.currentAction.Value == actions.Sleep {
					g.putTamagoInBed()
				} else if g.currentAction.Value == actions.Feed {
					g.feedTamago()
				} else if g.currentAction.Value == actions.Candy {
					g.giveACandy()
				} else if g.currentAction.Value == actions.Heal {
					g.healTamago()
				}
			}
		}

		if g.count > 0 {
			// Animation in progress, decrease the counter
			g.count--
		} else {
			// No animation in progress, let's update the animation
			if g.Tamago.State == status.Sleeping {
				g.putTamagoInBed()
			} else {
				// Default case
				// Update the animation with the current Tamago state
				g.currentAnimation = g.Tamago.State
			}
		}
	} else {
		if !g.deadMusicAlreadyPlayed {
			g.CurrentMusic <- "dead"
		}
	}

	return nil
}

// Draw : ebiten draw method
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	if g.Tamago.IsAlive() {
		if g.currentAnimation == status.Happy {
			screen.DrawImage(happyImage, op)
		} else if g.currentAnimation == status.Feeding {
			screen.DrawImage(feedingImage, op)
		} else if g.currentAnimation == status.Sleeping {
			screen.DrawImage(sleepingImage, op)
		} else if g.currentAnimation == status.Sad {
			screen.DrawImage(sadImage, op)
		} else if g.currentAnimation == status.Sick {
			screen.DrawImage(sickImage, op)
		}
		g.drawActionBar(screen)
	} else {
		screen.DrawImage(deadImage, op)
	}
}

// Layout : ebiten layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) feedTamago() {
	g.currentAnimation = status.Feeding
	g.count = feedingAnimationLength
	g.Tamago.Feed(food.Meat)
}

func (g *Game) giveACandy() {
	g.currentAnimation = status.Feeding
	g.count = feedingAnimationLength
	g.Tamago.Feed(food.Candy)
}

func (g *Game) putTamagoInBed() {
	g.currentAnimation = status.Sleeping
	g.count = sleepingAnimationLength
	g.Tamago.Bed()
}

func (g *Game) healTamago() {
	g.currentAnimation = status.Sleeping
	g.count = sleepingAnimationLength
	g.Tamago.Heal()
}

func (g *Game) drawActionBar(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(80.0, 250.0)

	if g.currentAction.Value == actions.Sleep {
		screen.DrawImage(actionSleepImage, op)
	} else if g.currentAction.Value == actions.Feed {
		screen.DrawImage(actionFeedImage, op)
	} else if g.currentAction.Value == actions.Candy {
		screen.DrawImage(actionCandyImage, op)
	} else if g.currentAction.Value == actions.Heal {
		screen.DrawImage(actionsHealImage, op)
	}
}
