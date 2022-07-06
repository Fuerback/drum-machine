# drum machine

This exercise assumes you are somewhat familiar with drum machines.
If you aren't please read http://en.wikipedia.org/wiki/Drum_machine.

Your challenge is to implement the *sequencer* portion of a drum machine.
A sequencer builds a pattern over time by adding individual sounds at certain points.

For example, here is a visualization of a [4-on-the-floor](https://en.wikipedia.org/wiki/Four_on_the_floor_(music)) rhythm pattern.

![A 4-on-the-floor drum pattern](https://upload.wikimedia.org/wikipedia/commons/thumb/c/c5/Four_to_the_floor_Roland_TR-707.jpg/330px-Four_to_the_floor_Roland_TR-707.jpg)

A sequencer would use this representation to play the pattern at a particular tempo.

## Instructions

Implement the following API in a programming language of your choice.

> The API is specified in Go, but feel free to use a different language.

```golang
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
```
