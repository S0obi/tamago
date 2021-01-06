package game

import (
	"container/ring"
	"fmt"
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
	// SampleRate : audio sample rate
	SampleRate = 22050

	// ScreenWidth : screen width
	ScreenWidth = 330
	// ScreenHeight : screen height
	ScreenHeight = 330
	// DefaultTPS : default TPS
	DefaultTPS = 60

	// timing
	feedingAnimationLength  = DefaultTPS * 2
	sleepingAnimationLength = DefaultTPS * 5
	cleaningAnimationLength = DefaultTPS * 2
)

var (
	// game images
	introImage *ebiten.Image

	// animation images
	happyImage    *ebiten.Image
	deadImage     *ebiten.Image
	feedingImage  *ebiten.Image
	sleepingImage *ebiten.Image
	sadImage      *ebiten.Image
	sickImage     *ebiten.Image
	hungryImage   *ebiten.Image
	starvingImage *ebiten.Image
	cleaningImage *ebiten.Image

	// action images
	actionFeedImage   *ebiten.Image
	actionCandyImage  *ebiten.Image
	actionSleepImage  *ebiten.Image
	actionsHealImage  *ebiten.Image
	actionsCleanImage *ebiten.Image

	// element images
	poopImage *ebiten.Image
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
	stateBeforePaused      status.Status
}

// Init : init method
func (g *Game) Init() {
	g.Tamago = tamagotchi.NewTamagotchi("Tama")

	g.Tamago.State = status.Paused
	g.stateBeforePaused = status.Happy

	// game images
	introImage, _, _ = ebitenutil.NewImageFromFile("assets/intro.png")

	// animation images
	happyImage, _, _ = ebitenutil.NewImageFromFile("assets/happy.png")
	deadImage, _, _ = ebitenutil.NewImageFromFile("assets/dead.png")
	feedingImage, _, _ = ebitenutil.NewImageFromFile("assets/miam.png")
	sleepingImage, _, _ = ebitenutil.NewImageFromFile("assets/dodo.png")
	sadImage, _, _ = ebitenutil.NewImageFromFile("assets/sad.png")
	sickImage, _, _ = ebitenutil.NewImageFromFile("assets/sick.png")
	hungryImage, _, _ = ebitenutil.NewImageFromFile("assets/hungry.png")
	starvingImage, _, _ = ebitenutil.NewImageFromFile("assets/starving.png")
	cleaningImage, _, _ = ebitenutil.NewImageFromFile("assets/cleaning.png")

	// action bar images
	actionFeedImage, _, _ = ebitenutil.NewImageFromFile("assets/actions/feed.png")
	actionCandyImage, _, _ = ebitenutil.NewImageFromFile("assets/actions/candy.png")
	actionSleepImage, _, _ = ebitenutil.NewImageFromFile("assets/actions/sleep.png")
	actionsHealImage, _, _ = ebitenutil.NewImageFromFile("assets/actions/heal.png")
	actionsCleanImage, _, _ = ebitenutil.NewImageFromFile("assets/actions/clean.png")

	// element images
	poopImage, _, _ = ebitenutil.NewImageFromFile("assets/poop.png")

	g.currentAction = actions.NewTamagoActions()

	// audio
	g.CurrentMusic = make(chan string)
	g.audioContext = audio.NewContext(SampleRate)
	g.audioPlayer = utils.NewMusicFromFile("assets/music/theme.mp3", g.audioContext)
	g.deadMusicAlreadyPlayed = false
	g.muteMusic = false
}

func (g *Game) inifiniteThemeMusic() {
	for {
		if !g.audioPlayer.IsPlaying() && g.Tamago.IsAlive() && !g.muteMusic {
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
				g.audioPlayer.Pause()
			} else {
				g.audioPlayer.Play()
			}
		}

		// Pause the game
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if g.Tamago.State == status.Paused {
				g.Tamago.State = g.stateBeforePaused
			} else {
				g.stateBeforePaused = g.Tamago.State
				g.Tamago.State = status.Paused
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
				} else if g.currentAction.Value == actions.Clean {
					g.cleanTamago()
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
		// Play the death song once :(
		if !g.deadMusicAlreadyPlayed {
			g.CurrentMusic <- "dead"
		}
	}

	return nil
}

// Draw : ebiten draw method
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	if g.Tamago.State != status.Paused {
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
			} else if g.currentAnimation == status.Hungry {
				screen.DrawImage(hungryImage, op)
			} else if g.currentAnimation == status.Starving {
				screen.DrawImage(starvingImage, op)
			} else if g.currentAnimation == status.Cleaning {
				screen.DrawImage(cleaningImage, op)
			}
			if g.Tamago.Dirty {
				g.drawDirtyElements(screen)
			}
			g.drawActionBar(screen)
		} else {
			screen.DrawImage(deadImage, op)
		}
		// Display the Tamago life in top left
		ebitenutil.DebugPrint(screen, fmt.Sprintf("%dHP", g.Tamago.Life))
	} else {
		// Tamago is paused, let's draw the intro image
		screen.DrawImage(introImage, op)
	}

}

// Layout : ebiten layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
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

func (g *Game) cleanTamago() {
	g.currentAnimation = status.Cleaning
	g.count = cleaningAnimationLength
	g.Tamago.Clean()
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
	} else if g.currentAction.Value == actions.Clean {
		screen.DrawImage(actionsCleanImage, op)
	}
}

func (g *Game) drawDirtyElements(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(20.0, 190.0)

	screen.DrawImage(poopImage, op)
}
