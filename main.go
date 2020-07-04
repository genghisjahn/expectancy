package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var chance int = 120
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
	var total float64
	rand.Seed(time.Now().UnixNano())
	mp := *maxpopulation
	for i := 0; i < mp; i++ {
		t := newTribble(1)
		population = append(population, &t)
	}
	for _, v := range population {
		v.Live()
		total += float64(v.Age)
	}

	sort.Slice(population, func(i, j int) bool {
		return population[i].Age < population[j].Age
	})

	lenpop := len(population)
	median := 0.0
	if lenpop%2 == 0 {
		f1 := float64(population[(lenpop-1)/2].Age)
		f2 := float64(population[((lenpop-1)/2)+1].Age)
		median = (f1 + f2) / 2.0
	} else {
		median = float64(population[(lenpop)/2].Age)
	}
	p := message.NewPrinter(language.English)
	s := p.Sprintf("%d", mp)
	fmt.Println("Population Size:", s)
	fmt.Println("Max Age:", population[len(population)-1].Age)
	fmt.Printf("Average: %.3f\n", total/float64(len(population)))
	fmt.Println("Median:", median)
	duration := time.Since(ts)
	fmt.Println("Duration: " + fmt.Sprintf("%d", duration.Milliseconds()) + " ms")
}
