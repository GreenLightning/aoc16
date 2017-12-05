package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"container/heap"
)

var input uint32

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	input = toUint32(scanner.Text())

	start  := point{  1,  1 }
	target := point{ 31, 39 }

	{
		fmt.Println("--- Part One ---")
		f := newFinder(start, target, 0, false)
		length, found := f.findShortestPath()
		if found {
			fmt.Println(length)
		}
	}

	{
		fmt.Println("--- Part Two ---")
		f := newFinder(start, point{}, 50, true)
		f.findShortestPath()
		fmt.Println(len(f.infos))
	}
}

func isWall(x, y int) bool {
	if x < 0 || y < 0 { return true }
	v := uint32(x*x + 3*x + 2*x*y + y + y*y) + input
	bits := uint32(0)
	for i := uint(0); i < 32; i++ {
		bits ^= (v >> i)
	}
	return bits & 1 == 1
}

type point struct {
	x, y int
}

type pointInfo struct {
	current point
	distance int
	index int
}

type finder struct {
	target point
	max int
	all bool
	queue pointInfoQueue
	infos map[point]*pointInfo
}

func newFinder(start point, target point, max int, all bool) *finder {
	f := &finder{}
	f.target = target
	f.max = max
	f.all = all
	f.infos = make(map[point]*pointInfo)
	info := &pointInfo{ start, 0, -1 }
	f.infos[start] = info
	f.queue.doPush(info)
	return f
}

func (f *finder) findShortestPath() (int, bool) {
	dist, found := 0, false
	for !f.queue.isEmpty() {
		info := f.queue.doPop()
		if info.current == f.target {
			dist, found = info.distance, true
			if !f.all { break }
		}
		x, y := info.current.x, info.current.y
		f.update(x+1, y  , info.distance+1)
		f.update(x-1, y  , info.distance+1)
		f.update(x  , y+1, info.distance+1)
		f.update(x  , y-1, info.distance+1)
	}
	return dist, found
}

func (f *finder) update(x, y int, distance int) {
	if isWall(x, y) || (f.all && distance > f.max) {
		return
	}
	next := point{ x, y }
	if info, ok := f.infos[next]; ok {
		if distance < info.distance {
			f.queue.doUpdate(info, distance)
		}
	} else {
		info := &pointInfo{ next, distance, -1 }
		f.infos[next] = info
		f.queue.doPush(info)
	}
}

type pointInfoQueue []*pointInfo

func (q *pointInfoQueue) isEmpty() bool {
	return len(*q) == 0
}

func (q *pointInfoQueue) doPush(info *pointInfo) {
	heap.Push(q, info)
}

func (q *pointInfoQueue) doPop() *pointInfo {
	return heap.Pop(q).(*pointInfo)
}

func (q *pointInfoQueue) doUpdate(info *pointInfo, distance int) {
	info.distance = distance
	heap.Fix(q, info.index)
}

func (q pointInfoQueue) Len() int { return len(q) }

func (q pointInfoQueue) Less(i, j int) bool {
	return q[i].distance < q[j].distance
}

func (q pointInfoQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *pointInfoQueue) Push(x interface{}) {
	info := x.(*pointInfo)
	info.index = len(*q)
	*q = append(*q, info)
}

func (q *pointInfoQueue) Pop() interface{} {
	n := len(*q)
	info := (*q)[n-1]
	info.index = -1
	*q = (*q)[0:n-1]
	return info
}

func toUint32(v string) uint32 {
	result, err := strconv.ParseUint(v, 10, 32)
	if err != nil { panic(err) }
	return uint32(result)
}
