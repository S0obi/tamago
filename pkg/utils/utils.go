package utils

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const (
	// audio
	sampleRate = 22050
)

// NewMusicFromFile : Get Music player from a MP3 file
func NewMusicFromFile(path string, audioContext *audio.Context) *audio.Player {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	s, err := mp3.Decode(audioContext, file)
	audioPlayer, err := audio.NewPlayer(audioContext, s)
	audioPlayer.SetVolume(0.05)
	return audioPlayer
}
