package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	eatTime = 0 * time.Second
	thinkTime = 0 * time.Second
	sleepTime = 0 * time.Second

	for i := 0; i < 10; i++ {
		finishedPhilosophers = []string{}
		dine()
		if len(finishedPhilosophers) != 5 {
			t.Errorf("incorrect length of slice; expected 5 but got %d", len(finishedPhilosophers))
		}
	}
}

func Test_dineWithVaryingDelays(t *testing.T) {
	var theTests = []struct {
		name  string
		delay time.Duration
	}{
		{"zero delay", time.Second * 0},
		{"quarter second delay", time.Millisecond * 250},
		{"half second delay", time.Millisecond * 500},
	}

	for _, e := range theTests {
		// Set finishedPhilosophers to an empty slice of strings.
		finishedPhilosophers = []string{}

		// Set all sleep times.
		eatTime = e.delay
		sleepTime = e.delay
		thinkTime = e.delay

		// Run the dine function.
		dine()

		// Perform our test.
		if len(finishedPhilosophers) != 5 {
			t.Errorf("%s: incorrect length of slice finishedPhilosophers; expected 5 but got %d", e.name, len(finishedPhilosophers))
		}
	}
}