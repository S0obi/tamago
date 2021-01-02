package utils

import (
	"image"
	_ "image/png" // need this for png support
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const (
	// audio
	sampleRate = 22050
)

// NewImageFromFile : Get image from a PNG file
func NewImageFromFile(path string) *ebiten.Image {
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

// NewMusicFromFile : Get Music player from a MP3 file
func NewMusicFromFile(path string) *audio.Player {
	audioContext := audio.NewContext(sampleRate)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	s, err := mp3.Decode(audioContext, file)
	audioPlayer, err := audio.NewPlayer(audioContext, s)
	audioPlayer.SetVolume(0.10)
	return audioPlayer
}
