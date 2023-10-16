package main

// The Dining Philosophers problem is well known in computer science circles.
// Five philosophers, numbered from 0 through 4, live in a house where the
// table is laid for them; each philosopher has their own place at the table.
// Their only difficulty – besides those of philosophy – is that the dish
// served is an unusual kind of spaghetti which has to be eaten with
// two forks. There are two forks next to each plate, so that presents no
// difficulty. As a consequence, however, this means that no two neighbours
// may be eating simultaneously, since there are five philosophers and five forks.
//
// This is a simple implementation of Dijkstra's solution to the "Dining
// Philosophers" dilemma.

import (
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
	Name      string
	RightFork int
	LeftFork  int
}

// list of philosophers
var philosophers = []Philosopher{
	{Name: "Aristoteles", LeftFork: 4, RightFork: 0},
	{Name: "Rumi", LeftFork: 0, RightFork: 1},
	{Name: "Ibn Sina", LeftFork: 1, RightFork: 2},
	{Name: "Plato", LeftFork: 2, RightFork: 3},
	{Name: "Pascal", LeftFork: 3, RightFork: 4},
}

var hunger = 3 // how many times does a person eat?
var eatTime = 1 * time.Second
var thinkTime = 3 * time.Second
var sleepTime = 1 * time.Second

// store philosopher who has finished eating
var finishedPhilosophers = []string{}
var finishedLock = &sync.Mutex{}

func main() {
	// print out a welcome message
	fmt.Println("-------------------------------")
	fmt.Println("# Dining Philosophers Problem #")
	fmt.Println("-------------------------------")
	fmt.Println("> The table is empty")

	// start the meal
	dine()

	// print out finished message
	fmt.Println("> The table is empty")

	// print out the philosophers in orders
	fmt.Println("")
	fmt.Println("")
	fmt.Println("-------------------------------")
	fmt.Println("#     Philosophers Orders     #")
	fmt.Println("-------------------------------")
	for i, philosopherName := range finishedPhilosophers {
		fmt.Printf("%d. %s\n", i+1, philosopherName)
	}
}

func dine() {
	// eatTime = 0 * time.Second
	// thinkTime = 0 * time.Second
	// sleepTime = 0 * time.Second

	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seatedWg := &sync.WaitGroup{}
	seatedWg.Add(len(philosophers))

	// forks is a map of all 5 forks
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// start the meal
	for i := 0; i < len(philosophers); i++ {
		// fire off a goroutine for current philosopher
		go diningProblem(philosophers[i], wg, forks, seatedWg)
	}

	wg.Wait()
}

// diningProblem is the function fired off as a goroutine for each of our philosophers. It takes one
// philosopher, our WaitGroup to determine when everyone is done, a map containing the mutexes for every
// fork on the table, and a WaitGroup used to pause execution of every instance of this goroutine
// until everyone is seated at the table.
func diningProblem(p Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seatedWg *sync.WaitGroup) {
	defer wg.Done()

	// seat the philosopher at the table
	fmt.Printf("> %s is seated at the table.\n", p.Name)
	seatedWg.Done()

	seatedWg.Wait()

	// eat three times
	for i := hunger; i > 0; i-- {
		// get lock on both forks
		if p.LeftFork > p.RightFork {
			forks[p.RightFork].Lock()
			fmt.Printf(">\t %s takes the right fork.\n", p.Name)
			forks[p.LeftFork].Lock()
			fmt.Printf(">\t %s takes the left fork.\n", p.Name)
		} else {
			forks[p.LeftFork].Lock()
			fmt.Printf(">\t %s takes the left fork.\n", p.Name)
			forks[p.RightFork].Lock()
			fmt.Printf(">\t %s takes the right fork.\n", p.Name)
		}

		fmt.Printf(">\t %s has both forks and is eating.\n", p.Name)
		time.Sleep(eatTime)

		fmt.Printf(">\t %s is thinking.\n", p.Name)
		time.Sleep(thinkTime)

		forks[p.LeftFork].Unlock()
		forks[p.RightFork].Unlock()
		fmt.Printf(">\t %s put down the forks.\n", p.Name)
	}

	fmt.Printf("> %s is satisfied.\n", p.Name)
	time.Sleep(sleepTime)
	fmt.Printf("> %s left the table.\n", p.Name)

	finishedLock.Lock()
	finishedPhilosophers = append(finishedPhilosophers, p.Name)
	finishedLock.Unlock()
}
