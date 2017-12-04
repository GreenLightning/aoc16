package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"container/heap"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	nodeRule := regexp.MustCompile(`/dev/grid/node-x(\d+)-y(\d+)\s+(\d+)T\s+(\d+)T\s+(\d+)T\s+(\d+)%`)

	var nodes []node

	maxX, maxY := -1, -1

	for scanner.Scan() {
		line := scanner.Text()
		if result := nodeRule.FindStringSubmatch(line); result != nil {
			x, y := toInt(result[1]), toInt(result[2])
			size, used, avail := toInt(result[3]), toInt(result[4]), toInt(result[5])
			nodes = append(nodes, node{ x, y, size, used, avail })
			if x > maxX { maxX = x }
			if y > maxY { maxY = y }
		}
	}

	fmt.Println("--- Part One ---")
	viablePairs := 0
	for i, a := range nodes {
		for j, b := range nodes {
			if i != j && a.used != 0 && a.used <= b.avail {
				viablePairs++
			}
		}
	}
	fmt.Println(viablePairs)

	fmt.Println("--- Part Two ---")
	width, height := maxX + 1, maxY + 1
	grid := make([][]bool, height)
	for y := 0; y < height; y++ {
		row := make([]bool, width)
		for x := 0; x < width; x++ {
			row[x] = true
		}
		grid[y] = row
	}
	emptyX, emptyY := -1, -1
	for _, node := range nodes {
		if node.used == 0 { emptyX, emptyY = node.x, node.y }
		for _, test := range nodes {
			if node.used > test.size {
				grid[node.y][node.x] = false
				break
			}
		}
	}
	f := newFinder(grid, width, height, state{ emptyX, emptyY, maxX, 0 })
	if length, ok := f.findShortestPath(); ok {
		fmt.Println(length)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x == 0 && y == 0 {
				fmt.Print("T")
			} else if x == maxX && y == 0 {
				fmt.Print("D")
			} else if x == emptyX && y == emptyY {
				fmt.Print("O")
			} else if grid[y][x] {
				fmt.Print(".")
			} else {
				fmt.Print("X")
			}
		}
		fmt.Println()
	}
}

type node struct {
	x, y int
	size, used, avail int
}

type state struct {
	emptyX, emptyY int
	dataX, dataY int
}

type stateInfo struct {
	current state
	done bool
	steps int
	index int
}

type finder struct {
	grid [][]bool
	width, height int
	queue stateInfoQueue
	infos map[state]*stateInfo
}

func newFinder(grid [][]bool, width, height int, start state) *finder {
	f := &finder{}
	f.grid = grid
	f.width, f.height = width, height
	f.infos = make(map[state]*stateInfo)
	info := &stateInfo{ start, false, 0, -1 }
	f.infos[start] = info
	f.queue.doPush(info)
	return f
}

func (f *finder) findShortestPath() (int, bool) {
	for !f.queue.isEmpty() {
		info := f.queue.doPop()
		info.done = true
		if info.current.dataX == 0 && info.current.dataY == 0 {
			return info.steps, true
		}
		x, y := info.current.emptyX, info.current.emptyY
		f.update(info.current, x+1, y  , info.steps+1)
		f.update(info.current, x-1, y  , info.steps+1)
		f.update(info.current, x  , y+1, info.steps+1)
		f.update(info.current, x  , y-1, info.steps+1)
	}
	return 0, false
}

func (f *finder) update(next state, x, y int, steps int) {
	if x < 0 || y < 0 || x >= f.width || y >= f.height || !f.grid[y][x] {
		return
	}
	if next.dataX == x && next.dataY == y {
		next.dataX, next.dataY = next.emptyX, next.emptyY
	}
	next.emptyX, next.emptyY = x, y
	if info, ok := f.infos[next]; ok {
		if steps < info.steps {
			f.queue.doUpdate(info, steps)
		}
	} else {
		info := &stateInfo{ next, false, steps, -1 }
		f.infos[next] = info
		f.queue.doPush(info)
	}
}

type stateInfoQueue []*stateInfo

func (q *stateInfoQueue) isEmpty() bool {
	return len(*q) == 0
}

func (q *stateInfoQueue) doPush(info *stateInfo) {
	heap.Push(q, info)
}

func (q *stateInfoQueue) doPop() *stateInfo {
	return heap.Pop(q).(*stateInfo)
}

func (q *stateInfoQueue) doUpdate(info *stateInfo, steps int) {
	info.steps = steps
	heap.Fix(q, info.index)
}

func (q stateInfoQueue) Len() int { return len(q) }

func (q stateInfoQueue) Less(i, j int) bool {
	return q[i].steps < q[j].steps
}

func (q stateInfoQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *stateInfoQueue) Push(x interface{}) {
	info := x.(*stateInfo)
	info.index = len(*q)
	*q = append(*q, info)
}

func (q *stateInfoQueue) Pop() interface{} {
	n := len(*q)
	info := (*q)[n-1]
	info.index = -1
	*q = (*q)[0:n-1]
	return info
}

func toInt(v string) int {
	i, e := strconv.Atoi(v)
	if e != nil { panic(e) }
	return i
}
