package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var chance int = 100
var wg sync.WaitGroup
var population []*tribble

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

type tribble struct {
	Name string
	Age  int
	Dead bool
	Incr int
}

var maxpopulation = 1000000

func newTribble(name string, incr int) tribble {
	return tribble{Name: name, Incr: incr}
}

func (t *tribble) Tick() {
	var r = 0
	defer func() {
		if err := recover(); err != nil {
			t.Tick()
		}
	}()
	r = r1.Intn(100) + 1
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
	//fmt.Printf("%v died at age %v\n", t.Name, t.Age)
	wg.Done()
}

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	//This is the program

	for i := 0; i < maxpopulation; i++ {
		wg.Add(1)
		t := newTribble(generateStupidName(), 1)
		population = append(population, &t)
		go t.Live()
	}
	wg.Wait()

	var total float64
	var max int
	for _, v := range population {
		if v.Age > max {
			max = v.Age
		}
		total += float64(v.Age)
	}
	fmt.Println("Max Age:", max)
	fmt.Println("Average:", total/float64(len(population)))

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
