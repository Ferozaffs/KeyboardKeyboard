package main

import (
	"log"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/generators"
	"github.com/gopxl/beep/v2/speaker"
)

var freqs = map[rune]float64{
	'z': 130.81,  //C3
	'x': 146.83,  //D3
	'c': 164.81,  //E3
	'v': 174.61,  //F3
	'b': 196.00,  //G3
	'n': 220.00,  //A3
	'm': 246.94,  //B3
	'a': 261.63,  //C4
	's': 293.66,  //D4
	'd': 329.63,  //E4
	'f': 349.23,  //F4
	'g': 392.00,  //G4
	'h': 440.00,  //A4
	'j': 493.88,  //B4
	'k': 523.25,  //C5
	'l': 587.33,  //D5
	'q': 659.25,  //E5
	'w': 698.46,  //F5
	'e': 783.99,  //G5
	'r': 880.00,  //A5
	't': 987.77,  //B5
	'y': 1046.50, //C6
	'u': 1174.66, //D6
	'i': 1318.51, //E6
	'o': 1396.91, //F6
	'p': 1567.98, //G6
}

var fade = 500

func main() {
	sr := beep.SampleRate(48000)
	speaker.Init(sr, 4800)

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyArrowUp {
			fade += 50
		}
		if key == keyboard.KeyArrowDown {
			fade = max(0.0, fade-50)
		}

		go playTone(sr, freqs[char], 0.5, time.Duration(fade)*time.Millisecond)
	}

}

func playTone(sr beep.SampleRate, freq float64, amplitude float64, duration time.Duration) {
	sine, err := generators.SineTone(sr, freq)
	if err != nil {
		panic(err)
	}

	fader := effects.Transition(
		beep.Take(sr.N(duration), sine),
		sr.N(duration),
		amplitude,
		0.0,
		effects.TransitionEqualPower)

	done := make(chan struct{})
	speaker.Play(beep.Seq(fader, beep.Callback(func() {
		done <- struct{}{}
	})))
	<-done
}
