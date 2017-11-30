package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"encoding/csv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	reader := csv.NewReader(file)
	record, err := reader.Read()
	if err != nil { panic(err) }

	dir := 0 // 0 up, 1 right, 2 down, 3 left
	pos := point{}
	visited := make(map[point]bool)
	var first *point

	visited[pos] = true

	for _, instr := range record {
		instr = strings.TrimSpace(instr)

		if instr[:1] == "R" {
			dir = (dir + 1) % 4;
		} else {
			dir = (dir + 3) % 4;
		}

		value, err := strconv.Atoi(instr[1:])
		if err != nil { panic(err) }

		for i := 0; i < value; i++ {
			pos.move(dir)

			if first == nil && visited[pos] {
				copy := pos
				first = &copy
			}

			visited[pos] = true
		}
	}

	fmt.Println("Final Location: ", pos.distance())
	if (first != nil) {
		fmt.Println("Actual Location:", first.distance())
	}
}

type point struct {
	x, y int
}

func (p *point) move(direction int) {
	switch (direction) {
		case 0: p.y += 1
		case 1: p.x += 1
		case 2: p.y -= 1
		case 3: p.x -= 1
	}
}

func (p *point) distance() int {
	return abs(p.x) + abs(p.y)
}

func abs(v int) int {
	if v < 0 { return -v }
	return v
}
