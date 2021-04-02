package main

import (
	"container/list"
	"fmt"
	"math"
	"sync"
	"time"
)

type CStatus int

const (
	IDLE CStatus = iota
	FIRST
	SECOND
	RESULT
	ROOT
)

type Node struct {
	parent      *Node
	cStatus     CStatus
	locked      bool
	result      int
	firstValue  int
	secondValue int
	release     chan bool
}

func newNode(myParent Node) Node {
	n := Node{parent: &myParent}
	n.cStatus = IDLE
	n.locked = false
	n.result = 0
	n.firstValue = 0
	n.secondValue = 0
	n.release = make(chan bool)
	return n
}

func precombine(node Node) bool {
	if node.locked {
		<-node.release
	}
	switch node.cStatus {
	case IDLE:
		node.cStatus = FIRST
		return true
	case FIRST:
		node.locked = true
		node.cStatus = SECOND
		return false
	case ROOT:
		return false
	default:
		println("unexpected Node state in precombine")
		return false // error
	}
}

func combine(combined int, node Node) int {
	if node.locked {
		<-node.release
	}
	node.locked = true
	node.firstValue = combined
	switch node.cStatus {
	case FIRST:
		return node.firstValue
	case SECOND:
		return node.firstValue + node.secondValue
	default:
		println("unexpected Node state in combine")
		return -1 // error

	}
}

func op(combined int, node Node) int {
	switch node.cStatus {
	case ROOT:
		prior := node.result
		node.result += combined
		return prior
	case SECOND:
		node.secondValue = combined
		node.locked = false
		node.release <- true // wake up waiting threads
		if node.cStatus != RESULT {
			<-node.release
		}
		node.locked = false
		node.release <- true
		node.cStatus = IDLE
		return node.result
	default:
		println("unexpected Node state in op")
		return -1 // error

	}
}

func distribute(prior int, node Node) {
	switch node.cStatus {
	case FIRST:
		node.cStatus = IDLE
		node.locked = false
		break
	case SECOND:
		node.result = prior + node.firstValue
		node.cStatus = RESULT
		break
	default:
		println("unexpected Node state in distribute")
	}
	node.release <- true
}

type CombinningTree struct {
	nodes []Node
	leaf  []Node
}

func newTree(width int) CombinningTree {

	tree := CombinningTree{}
	tree.nodes = make([]Node, 2*width-1)
	tree.leaf = make([]Node, width)

	length := 2*width - 1
	tree.nodes[0] = newNode(Node{parent: nil, cStatus: ROOT})
	for i := 1; i < length; i++ {
		tree.nodes[i] = newNode(tree.nodes[(i-1)/2])
	}
	for i := 1; i < width; i++ {
		tree.leaf[i] = tree.nodes[length-i-1]
	}
	return tree
}

func getAndIncrement(tree CombinningTree, id int) int {
	stack := list.New()
	myLeaf := tree.leaf[id%len(tree.leaf)]
	println(id, len(tree.leaf))
	node := myLeaf
	for node.parent != nil && precombine(node) {
		node = *node.parent
	}
	stop := node
	combined := 1
	for node = myLeaf; node != stop; node = *node.parent {
		combined = combine(combined, node)
		stack.PushBack(node)
	}

	prior := op(combined, stop)
	for stack.Len() > 0 {

		node = stack.Remove(stack.Back()).(Node)
		distribute(prior, node)
	}
	return prior

}

var TH int
var NUM int

func main() {
	TH = 12
	NUM = 100
	width := int(math.Ceil(float64(TH / 2)))
	tree := newTree(width)
	print(width)
	print("-leaves Combining tree.\n")
	print("Starting ")
	print(TH)
	print(" threads doing increments ...\n")
	var wg sync.WaitGroup

	start := time.Now()
	for index := 0; index < TH; index++ {
		wg.Add(1)
		go func(id int) {
			for i := 0; i < NUM; i++ {
				getAndIncrement(tree, id)
			}
		}(index)
	}

	wg.Wait()
	stop := time.Now()
	diff := stop.Sub(start)
	print("Total: ")
	print(tree.nodes[0].result)
	fmt.Println("\nTotal time: ", diff)

}
