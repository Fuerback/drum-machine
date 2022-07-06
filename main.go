package main

import (
	"bufio"
	"fmt"
	"strings"
)

type Pattern struct {
	instrumentNames []string
	track           [][16]bool
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

type drumMachine struct{}

func NewDrumMachine() Sequencer {
	return &drumMachine{}
}

func (d *drumMachine) Parse(pattern string) (Pattern, error) {
	//track := make(map[string][]string)
	scanner := bufio.NewScanner(strings.NewReader(pattern))
	for scanner.Scan() {
		before, after, found := strings.Cut(scanner.Text(), "|") // get instrument name
		if !found {
			continue
		}
		instrumentName := strings.TrimSpace(before) // remove white spaces from instrument name
		sequence := strings.Trim(after, "|")        // remove last | from sequence
		//sequenceList := strings.Split(sequence, "|")
		//if len(sequenceList) != 4 { // verify if sequence has four steps
		//	continue
		//}

		fmt.Println(instrumentName, sequence)
	}
	return Pattern{}, nil
}

//|hi-hat,kick|-|hi-hat|-|hi-hat,snare|-|hi-hat|-|hi-hat,kick|-|hi-hat|-|hi-hat,snare|-|hi-hat|-|
func (d *drumMachine) Render(pattern Pattern) (string, error) {
	//fmt.Println(len(pattern.track), len(pattern.track[0]))

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
			if pattern.track[j][i] {
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

	return play, nil
}

func (d *drumMachine) Play(bpm int32) error {
	return nil
}

func main() {
	//pattern := "hi-hat |x-x-|x-x-|x-x-|x-x-|\nsnare  |----|x---|----|x---|\nkick   |x---|----|x---|----|"
	drumMachine := NewDrumMachine()
	//drumMachine.Parse(pattern)

	instrumentNames := []string{"hi-hat", "snare", "kick"}
	track := [][16]bool{
		{true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false},
		{false, false, false, false, true, false, false, false, false, false, false, false, true, false, false, false},
		{true, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false},
	}
	p := Pattern{track: track, instrumentNames: instrumentNames}
	drumMachine.Render(p)

}
