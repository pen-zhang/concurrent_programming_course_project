
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
    "os"
    "strconv"
    "log"
    "strings"
)

type Balancer struct {
	toggle bool
	lock   sync.Mutex
}

func (b *Balancer) balancerTraverse() int {
	b.lock.Lock()
	defer b.lock.Unlock()
	defer func() { b.toggle = !b.toggle }()
	if b.toggle {
		return 0
	} else {
		return 1
	}
}

type Merger struct {
	width int
	layer []Balancer
	half  []Merger
}

func newMerger(myWidth int) Merger {
	m := Merger{}
	m.width = myWidth
	m.layer = make([]Balancer, myWidth/2)
	for i := 0; i < m.width/2; i++ {
		m.layer[i] = Balancer{toggle: true}
	}
	if m.width > 2 {
		m.half = []Merger{newMerger(myWidth / 2), newMerger(myWidth / 2)}
	}
	return m
}

func (m *Merger) mergerTraverse(input int) int {
	output := 0
	if m.width <= 2 {
		return m.layer[0].balancerTraverse()
	}
	if input < m.width/2 {
		output = m.half[input%2].mergerTraverse(input / 2)
	} else {
		output = m.half[1-(input%2)].mergerTraverse(input / 2)
	}
	return (2 * output) + m.layer[output].balancerTraverse()

}

type Bitonic struct {
	half   []Bitonic
	merger Merger
	width  int
}

func newBitonic(myWidth int) Bitonic {
	bit := Bitonic{}
	bit.width = myWidth
	bit.merger = newMerger(myWidth)
	if myWidth > 2 {
		bit.half = []Bitonic{newBitonic(myWidth / 2), newBitonic(myWidth / 2)}
	}
	return bit
}

func (bit *Bitonic) bitonicTraverse(input int) int {
	output := 0
	subnet := input / (bit.width / 2)
	if bit.width > 2 {
		output = bit.half[subnet].bitonicTraverse(input / 2)
	}
	if input >= bit.width/2 {
		return bit.merger.mergerTraverse(bit.width/2 + output)
	} else {
		return bit.merger.mergerTraverse(output)
	}
}

func main() {
	width, _ := strconv.Atoi(os.Args[1])
	counters := make([]int, width)
	bitonic := newBitonic(width)
	tokenCount, _ := strconv.Atoi(os.Args[2])
	tokens := make([]int, tokenCount)
	var wg sync.WaitGroup

	for i := 0; i < tokenCount; i++ {
		tokens[i] = rand.Intn(width)
	}

	start := time.Now()
	for idx := 0; idx < tokenCount; idx++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			counters[bitonic.bitonicTraverse(tokens[i])] += 1
		}(idx)
	}

	wg.Wait()
	stop := time.Now()
	diff := stop.Sub(start)

	for i := 0; i < width; i++ {
		fmt.Println("output: ", i, ", Count: ", counters[i])
	}
	fmt.Println("Total time to traverse the network: ", diff)
    
    f, err := os.OpenFile("go-counting.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    _, err2 := f.WriteString("go,"+strconv.Itoa(width)+","+strconv.Itoa(tokenCount)+","+strings.Split(diff.String(),"ms")[0]+"\n")
    if err2 != nil {
        log.Fatal(err2)
    }

    fmt.Println("done")
}
