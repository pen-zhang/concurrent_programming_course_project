
package main

import (
	"container/list"
	"fmt"
	"math"
	"sync"
	"time"
    "os"
    "strconv"
    "log"
    "strings"
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
	cond sync.Cond
}

func newNode(myParent *Node) Node {
	n := Node{parent: myParent}
	n.cStatus = IDLE
	n.locked = false
	n.result = 0
	n.firstValue = 0
	n.secondValue = 0
	n.cond = sync.Cond{}
	n.cond.L = &sync.Mutex{}
	return n
}

func (node *Node) precombine() bool {
	node.cond.L.Lock()
	defer node.cond.L.Unlock()
	for node.locked {
		node.cond.Wait()
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
		fmt.Println("unexpected Node state in precombine", node.cStatus)
		return false // error
	}
}

func (node *Node) combine(combined int) int {
	node.cond.L.Lock()
	defer node.cond.L.Unlock()
	for node.locked {
		node.cond.Wait()
	}
	node.locked = true
	node.firstValue = combined
	switch node.cStatus {
	case FIRST:
		return node.firstValue
	case SECOND:
		return node.firstValue + node.secondValue
	default:
		fmt.Println("unexpected Node state in combine", node.cStatus)
		return -1 // error

	}
}

func (node *Node) op(combined int) int {
	node.cond.L.Lock()
	defer node.cond.L.Unlock()
	switch node.cStatus {
	case ROOT:
		prior := node.result
		node.result += combined
		return prior
	case SECOND:
		node.secondValue = combined
		node.locked = false
		node.cond.Broadcast()  // wake up waiting threads
		for node.cStatus != RESULT {
			node.cond.Wait()
		}
		node.locked = false
		node.cond.Broadcast()
		node.cStatus = IDLE
		return node.result
	default:
		fmt.Println("unexpected Node state in op", node.cStatus)
		return -1 // error

	}
}

func (node *Node) distribute(prior int) {
	node.cond.L.Lock()
	defer node.cond.L.Unlock()
	switch node.cStatus {
	case FIRST:
		node.cStatus = IDLE
		node.locked = false
	case SECOND:
		node.result = prior + node.firstValue
		node.cStatus = RESULT
	default:
		fmt.Println("unexpected Node state in distribute", node.cStatus)
	}
	node.cond.Broadcast()
}

type CombinningTree struct {
	nodes []Node
	leaf  []*Node
}

func newTree(width int) CombinningTree {

	tree := CombinningTree{}
	tree.nodes = make([]Node, 2*width-1)
	tree.leaf = make([]*Node, width)
	length := 2*width - 1
	tree.nodes[0] = Node{
		parent:      nil,
		cStatus:     ROOT,
		locked:      false,
		result:      0,
		firstValue:  0,
		secondValue: 0,
		cond: sync.Cond{},
	}
	tree.nodes[0].cond.L = &sync.Mutex{}
	for i := 1; i < length; i++ {
		tree.nodes[i] = newNode(&tree.nodes[(i-1)/2])
	}
	for i := 0; i < width; i++ {
		tree.leaf[i] = &tree.nodes[length-i-1]
	}
	return tree
}

func (tree *CombinningTree) getAndIncrement(id int) int {
	stack := list.New()
	myLeaf := tree.leaf[id%len(tree.leaf)]
	node := myLeaf
	for node.precombine() {
		node = node.parent
	}
	stop := node
	combined := 1
	for node = myLeaf; node != stop; node = node.parent {
		combined = node.combine(combined)
		stack.PushBack(node)
	}
	prior := stop.op(combined)
	for stack.Len() > 0 {
		node = stack.Remove(stack.Back()).(*Node)
		node.distribute(prior)
	}
	return prior

}

var TH int
var NUM int

func main() {
	TH, _ = strconv.Atoi(os.Args[1])
	NUM, _ = strconv.Atoi(os.Args[2])
	width := int(math.Ceil(float64(TH / 2)))
	tree := newTree(width)
	fmt.Println(width, "-leaves Combining tree.")
	fmt.Println("Starting", TH ,"threads doing increments ...")
	var wg sync.WaitGroup

	start := time.Now()
	for index := 0; index < TH; index++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			s := time.Now()
			for i := 0; i < NUM; i++ {
				tree.getAndIncrement(id)
			}
			e := time.Now()
			fmt.Println(id, "done in", e.Sub(s))
		}(index)
	}
	wg.Wait()
	stop := time.Now()
	diff := stop.Sub(start)
	fmt.Println("Total:", tree.nodes[0].result)
	fmt.Println("Total time:", diff)
    
    f, err := os.OpenFile("go-Combining.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    _, err2 := f.WriteString("go,"+strconv.Itoa(TH)+","+strconv.Itoa(NUM)+","+strings.Split(diff.String(),"ms")[0]+"\n")
    if err2 != nil {
        log.Fatal(err2)
    }

    fmt.Println("done")
    

}
