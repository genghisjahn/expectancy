package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var chance int = 100
var wg sync.WaitGroup

type tribble struct {
	Name string
	Age  int
	Dead bool
	Incr int
}

func newTribble(name string, incr int) tribble {
	return tribble{Name: name, Incr: incr}
}

func (t *tribble) Tick() {
	r := rand.Intn(100)
	if r <= t.Incr {
		t.Dead = false
		return
	}
	t.Incr++
	t.Age++
}

func (t *tribble) Live() {
	for {
		if t.Dead {
			break
		}
		t.Tick()
	}
	fmt.Printf("%v died at age %v", t.Name, t.Age)
	wg.Done()
}

func main() {

	wg.Add(1)
	abe := newTribble("Abe", 1)
	abe.Live()
	wg.Wait()
}
