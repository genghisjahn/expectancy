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

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var chance int = 100
var wg sync.WaitGroup
var population []*tribble

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var maxpopulation = flag.Int("maxpop", 100, "max population, default is 10")

type tribble struct {
	Name string
	Age  int
	Dead bool
	Incr float64
}

func newTribble(incr float64) tribble {
	return tribble{Incr: incr}
}

func (t *tribble) Tick() {
	var r float64
	r = rand.Float64() * float64(chance)
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

	worldLoop()

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
func worldLoop() {
	ts := time.Now()
	rand.Seed(time.Now().UnixNano())
	mp := *maxpopulation
	for i := 0; i < mp; i++ {
		t := newTribble(1)
		population = append(population, &t)
	}
	for _, v := range population {
		wg.Add(1)
		go v.Live()
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
	p := message.NewPrinter(language.English)
	s := p.Sprintf("%d", mp)
	fmt.Println("Population Size:", s)
	fmt.Println("Max Age:", max)
	fmt.Println("Average:", total/float64(len(population)))

	duration := time.Since(ts)
	fmt.Println("Duration: " + fmt.Sprintf("%v", duration.Seconds()) + " seconds")
}
