package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Pattern struct {
	instrumentNames []string
	track           [][]bool
}

type Sequencer interface {
	// Parse parses a string representation of a pattern into a structured Pattern.
	//
	// The input string pattern is of the form:
	// hi-hat |x-x-|x-x-|x-x-|x-x-|
	// snare  |----|x---|----|x---|
	// kick   |x---|----|x---|----|
	//
	// Each line of a pattern file represents a track.
	// There is no limit to the number of tracks in a pattern.
	// A track contains an instrument name and a 16-step sequence.
	// The instrument name is an identifier and should only appear once per pattern.
	// Each sequence represents a single measure in 4/4 time divided into 16th
	// note steps (`x` for play and `-` for silent).
	Parse(pattern string) (Pattern, error)

	// Render returns a string that represents a single "play" of the pattern.
	//
	// The returned string is of the form:
	// |hi-hat,kick|-|hi-hat|-|hi-hat,snare|-|hi-hat|-|hi-hat,kick|-|hi-hat|-|hi-hat,snare|-|hi-hat|-|
	//
	// In other words, it displays the instrument(s) that should be played in
	// each of the pattern's 16 steps.
	Render(pattern Pattern) (string, error)

	// Play prints the output of Render at the specified tempo (aka, beats per minute).
	// For example, with a bpm of 60, Play will print one step per second.
	Play(bpm int32) error
}

type drumMachine struct {
	render string
}

func NewDrumMachine() Sequencer {
	return &drumMachine{}
}

// expected return
//instrumentNames := []string{"hi-hat", "snare", "kick"}
//track := [][]bool{
//	{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false},
//	{false, false, false, false, true, false, false, false, false, false, false, false, true, false, false, false},
//	{true, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false},
//}
func (d *drumMachine) Parse(pattern string) (Pattern, error) {
	var instrumentNames []string
	track := make([][]bool, 0)

	scanner := bufio.NewScanner(strings.NewReader(pattern)) // reading line by line
	for scanner.Scan() {
		instrumentName, sequence, err := getInstrumentAndSequence(scanner.Text())
		if err != nil {
			return Pattern{}, err
		}

		if contains(instrumentNames, instrumentName) { // not the fastest way, but we should check the duplicated instrument names
			return Pattern{}, errors.New("duplicated instrument name")
		}

		instrumentNames = append(instrumentNames, instrumentName)
		var trackRow []bool
		for _, v := range sequence { // read sequence and saving the booleans in a list
			trackRow = append(trackRow, getBooleanPlay(string(v)))
		}

		track = append(track, trackRow) // saving the track row on the 2d slice
	}

	return Pattern{instrumentNames: instrumentNames, track: track}, nil
}

func getInstrumentAndSequence(line string) (string, string, error) {
	name, seq, found := strings.Cut(line, "|") // get instrument name and track
	if !found {
		return "", "", errors.New("incorrect format")
	}
	instrumentName := strings.TrimSpace(name)    // remove white spaces from instrument name
	sequence := strings.ReplaceAll(seq, "|", "") // remove all | from sequence
	return instrumentName, sequence, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getBooleanPlay(beat string) bool {
	if beat != "-" {
		return true
	}
	return false
}

// expected return |hi-hat,kick|-|hi-hat|-|hi-hat,snare|-|hi-hat|-|hi-hat,kick|-|hi-hat|-|hi-hat,snare|-|hi-hat|-|
func (d *drumMachine) Render(pattern Pattern) (string, error) {
	//fmt.Println(pattern.track[0][0])
	//fmt.Println(pattern.track[1][0])
	//fmt.Println(pattern.track[2][0])
	//fmt.Println(pattern.track[0][1])
	//fmt.Println(pattern.track[1][1])
	//fmt.Println(pattern.track[2][1])

	var columnPlay []string
	var play string
	divisor := "|"

	for i := 0; i < len(pattern.track[0]); i++ { // iterate over column
		for j := 0; j < len(pattern.track); j++ { // iterate over rows
			if pattern.track[j][i] { // if it has a beat, store the instrument
				columnPlay = append(columnPlay, pattern.instrumentNames[j])
			}
		}

		if len(columnPlay) > 0 {
			play += strings.Join(columnPlay, ",") // join comma
			play += divisor                       // add divisor
			columnPlay = nil                      // empty slice
			continue                              // jump to the next column
		}

		play += "-" + divisor
	}

	if len(play) > 0 { // if there is a play, add divisor prefix
		play = divisor + play
	}

	d.render = play // add new render on drum machine

	return play, nil
}

func (d *drumMachine) Play(bpm int32) error {
	render := strings.Trim(d.render, "|") // remove prefix and sufix from render - avoiding spaces in slice
	beats := strings.Split(render, "|")   // split render in beats

	beatsPerSecond := float32(bpm) / 60 // get number of beats per second
	fmt.Println(beatsPerSecond, "beats per second")
	milisecondsToBeat := 1000 / beatsPerSecond // get the time in milisecond to wait until next beat
	fmt.Println(milisecondsToBeat, "miliseconds")

	for i, beat := range beats {
		if i == 0 {
			fmt.Print("|" + beat + "|")
			time.Sleep(time.Millisecond * time.Duration(milisecondsToBeat))
			continue
		}
		fmt.Print(beat + "|")
		time.Sleep(time.Millisecond * time.Duration(milisecondsToBeat))
	}

	return nil
}

func main() {
	// hi-hat |x-x-|x-x-|x-x-|x-x-|
	// snare  |----|x---|----|x---|
	// kick   |x---|----|x---|----|
	pattern := "hi-hat |x-x-|x-x-|x-x-|x-x-|\nsnare  |----|x---|----|x---|\nkick   |x---|----|x---|----|"
	drumMachine := NewDrumMachine()
	p, err := drumMachine.Parse(pattern)
	if err != nil {
		fmt.Println("error on parse pattern:", err)
	}

	//instrumentNames := []string{"hi-hat", "snare", "kick"}
	//track := [][]bool{
	//	{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false},
	//	{false, false, false, false, true, false, false, false, false, false, false, false, true, false, false, false},
	//	{true, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false},
	//}
	//p := Pattern{track: track, instrumentNames: instrumentNames}
	_, err = drumMachine.Render(p)
	if err != nil {
		fmt.Println("error on Render pattern:", err)
	}

	drumMachine.Play(30)

}
