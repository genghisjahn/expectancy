package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var chance int = 100
var wg sync.WaitGroup
var population []tribble

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
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	r := r1.Intn(100) + 1
	if r <= t.Incr {
		t.Dead = true
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
	fmt.Printf("%v died at age %v\n", t.Name, t.Age)
	wg.Done()
}

func main() {
	for i := 0; i < 1; i++ {
		wg.Add(1)
		t := newTribble(generateStupidName(), 1)
		t.Live()
	}

	wg.Wait()
}
