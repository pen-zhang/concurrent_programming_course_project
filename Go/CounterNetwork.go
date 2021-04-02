package main

import "math/rand"
import "time"
import "fmt"
import "sync"

type Balancer struct {
	toggle bool
	lock   sync.Mutex
}

func balancerTraverse(b *Balancer) int {
	b.lock.Lock()
	fmt.Println(b.toggle)
	defer func() { b.toggle = !b.toggle }()
	defer b.lock.Unlock()
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

func newMerger(myWidth int) *Merger {
	m := Merger{}
	m.width = myWidth
	m.layer = make([]Balancer, myWidth/2)
	for i := 0; i < m.width/2; i++ {
		m.layer[i] = Balancer{toggle: true}
	}
	if m.width > 2 {
		m.half = []Merger{*newMerger(myWidth / 2), *newMerger(myWidth / 2)}
	}
	return &m
}

func mergerTraverse(input int, m *Merger) int {
	output := 0
	if m.width <= 2 {
		return balancerTraverse(&m.layer[0])
	}
	if input < m.width/2 {
		output = mergerTraverse(input/2, &m.half[input%2])
	} else {
		output = mergerTraverse(input/2, &m.half[1-(input%2)])
	}
	return (2 * output) + balancerTraverse(&m.layer[output])

}

type Bitonic struct {
	half   []Bitonic
	merger Merger
	width  int
}

func newBitonic(myWidth int) *Bitonic {
	bit := Bitonic{}
	bit.width = myWidth
	bit.merger = *newMerger(myWidth)
	if myWidth > 2 {
		bit.half = []Bitonic{*newBitonic(myWidth / 2), *newBitonic(myWidth / 2)}
	}
	return &bit
}

func bitonicTraverse(input int, bit *Bitonic) int {
	output := 0
	subnet := input / (bit.width / 2)
	if bit.width > 2 {
		output = bitonicTraverse(input/2, &bit.half[subnet])
	}
	if input >= bit.width/2 {
		return mergerTraverse(bit.width/2+output, &bit.merger)
	} else {
		return mergerTraverse(output, &bit.merger)
	}
}

func main() {
	width := 4
	counters := make([]int, width)
	bitonic := newBitonic(width)
	tokenCount := 10
	tokens := make([]int, tokenCount)

	for i := 0; i < tokenCount; i++ {
		tokens[i] = rand.Intn(width)
	}
	var wg sync.WaitGroup

	start := time.Now()
	for idx := 0; idx < tokenCount; idx++ {
		wg.Add(1)
		go func(i int) {
			counters[bitonicTraverse(tokens[i], bitonic)] += 1
		}(idx)
	}

	time.Sleep(10000)
	// wg.Wait()
	stop := time.Now()
	diff := stop.Sub(start)

	for i := 0; i < width; i++ {
		fmt.Println("output: ", i, ", Count: ", counters[i])
	}
	fmt.Println("Total time to traverse the network: ", diff)
}
