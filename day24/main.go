package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"container/heap"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lines := make([]string, 0)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	width, height := len(lines[0]), len(lines)

	nodes := make([][]*node, height)
	numbered := make([]*node, 10)

	for y := 0; y < height; y++ {
		row := make([]*node, width)
		for x := 0; x < width; x++ {
			if lines[y][x] != '#' {
				node := &node{}
				node.id = strings.IndexByte("0123456789", lines[y][x])
				if node.id >= 0 { numbered[node.id] = node }
				node.x, node.y = x, y
				row[x] = node
			}
		}
		nodes[y] = row
	}

	{ // reslice numbered
		count := len(numbered)
		for count > 0 && numbered[count-1] == nil { count-- }
		numbered = numbered[0:count]
	}

	{ // connect nodes
		for y := 1; y < height-1; y++ {
			for x := 1; x < width-1; x++ {
				if node := nodes[y][x]; node != nil {
					node.upNode    = nodes[y-1][x]
					node.downNode  = nodes[y+1][x]
					node.leftNode  = nodes[y][x-1]
					node.rightNode = nodes[y][x+1]
					if node.upNode    != nil { node.upDist    = 1 }
					if node.downNode  != nil { node.downDist  = 1 }
					if node.leftNode  != nil { node.leftDist  = 1 }
					if node.rightNode != nil { node.rightDist = 1 }
				}
			}
		}
	}

	{ // remove dead ends
		changed := true
		for changed {
			changed = false
			for y := 1; y < height-1; y++ {
				for x := 1; x < width-1; x++ {
					if node := nodes[y][x]; node != nil && node.id < 0 {
						count := 0
						if node.upNode    != nil { count++ }
						if node.downNode  != nil { count++ }
						if node.leftNode  != nil { count++ }
						if node.rightNode != nil { count++ }
						if count <= 1 {
							if node.upNode    != nil { node.upNode.downNode = nil }
							if node.downNode  != nil { node.downNode.upNode = nil }
							if node.leftNode  != nil { node.leftNode.rightNode = nil }
							if node.rightNode != nil { node.rightNode.leftNode = nil }
							nodes[y][x] = nil
							changed = true
						}
					}
				}
			}
		}
	}

	{ // remove straight passages
		changed := true
		for changed {
			changed = false
			for y := 1; y < height-1; y++ {
				for x := 1; x < width-1; x++ {
					if node := nodes[y][x]; node != nil && node.id < 0 {
						if node.upNode != nil && node.downNode != nil && node.leftNode == nil && node.rightNode == nil {
							node.upNode.downNode  = node.downNode
							node.upNode.downDist += node.downDist
							node.downNode.upNode  = node.upNode
							node.downNode.upDist += node.upDist
							nodes[y][x] = nil
							changed = true
						} else if node.upNode == nil && node.downNode == nil && node.leftNode != nil && node.rightNode != nil {
							node.leftNode.rightNode  = node.rightNode
							node.leftNode.rightDist += node.rightDist
							node.rightNode.leftNode  = node.leftNode
							node.rightNode.leftDist += node.leftDist
							nodes[y][x] = nil
							changed = true
						}
					}
				}
			}
		}
	}

	shortestPath := make([][]int, len(numbered))
	for i := 0; i < len(shortestPath); i++ {
		f := newFinder(numbered[i], len(numbered))
		shortestPath[i] = f.findShortestPaths()
	}

	available := make([]int, len(numbered)-1)
	for i := 0; i < len(available); i++ {
		available[i] = i + 1
	}

	fmt.Println("--- Part One ---")
	fmt.Println(findBestPermutation(shortestPath, 0, available, false))

	fmt.Println("--- Part Two ---")
	fmt.Println(findBestPermutation(shortestPath, 0, available, true))
}

type node struct {
	id int
	x, y int
	upNode, downNode, leftNode, rightNode *node
	upDist, downDist, leftDist, rightDist int
}

func findBestPermutation(shortestPath [][]int, current int, available []int, goBack bool) int {
	if len(available) == 0 {
		length := 0
		if goBack {
			length = shortestPath[current][0]
		}
		return length
	} else {
		best, found := 0, false
		for i := 0; i < len(available); i++ {
			remaining := make([]int, len(available)-1)
			for j, k := 0, 0; j < len(available); j++ {
				if j != i {
					remaining[k] = available[j]
					k++
				}
			}
			length := shortestPath[current][available[i]] + findBestPermutation(shortestPath, available[i], remaining, goBack)
			if !found || length < best {
				best, found = length, true
			}
		}
		return best
	}
}

type nodeInfo struct {
	node *node
	distance int
	index int
}

type finder struct {
	result []int
	queue nodeInfoQueue
	infos map[*node]*nodeInfo
}

func newFinder(start *node, length int) *finder {
	f := &finder{}
	f.result = make([]int, length)
	f.infos = make(map[*node]*nodeInfo)
	info := &nodeInfo{ start, 0, -1 }
	f.infos[start] = info
	f.queue.doPush(info)
	return f
}

func (f *finder) findShortestPaths() []int {
	for !f.queue.isEmpty() {
		info := f.queue.doPop()
		if info.node.id >= 0 {
			f.result[info.node.id] = info.distance
		}
		f.update(info.node.upNode   , info.distance + info.node.upDist   )
		f.update(info.node.downNode , info.distance + info.node.downDist )
		f.update(info.node.leftNode , info.distance + info.node.leftDist )
		f.update(info.node.rightNode, info.distance + info.node.rightDist)
	}
	return f.result
}

func (f *finder) update(next *node, distance int) {
	if next == nil { return }
	if info, ok := f.infos[next]; ok {
		if distance < info.distance {
			f.queue.doUpdate(info, distance)
		}
	} else {
		info := &nodeInfo{ next, distance, -1 }
		f.infos[next] = info
		f.queue.doPush(info)
	}
}

type nodeInfoQueue []*nodeInfo

func (q *nodeInfoQueue) isEmpty() bool {
	return len(*q) == 0
}

func (q *nodeInfoQueue) doPush(info *nodeInfo) {
	heap.Push(q, info)
}

func (q *nodeInfoQueue) doPop() *nodeInfo {
	return heap.Pop(q).(*nodeInfo)
}

func (q *nodeInfoQueue) doUpdate(info *nodeInfo, distance int) {
	info.distance = distance
	heap.Fix(q, info.index)
}

func (q nodeInfoQueue) Len() int { return len(q) }

func (q nodeInfoQueue) Less(i, j int) bool {
	return q[i].distance < q[j].distance
}

func (q nodeInfoQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *nodeInfoQueue) Push(x interface{}) {
	info := x.(*nodeInfo)
	info.index = len(*q)
	*q = append(*q, info)
}

func (q *nodeInfoQueue) Pop() interface{} {
	n := len(*q)
	info := (*q)[n-1]
	info.index = -1
	*q = (*q)[0:n-1]
	return info
}
