package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var discs []disc

	rule := regexp.MustCompile(`Disc #\d+ has (\d+) positions; at time=0, it is at position (\d+).`)

	for scanner.Scan() {
		line := scanner.Text()
		data := rule.FindStringSubmatch(line)
		d := disc{ toInt(data[1]), toInt(data[2]) }
		discs = append(discs, d)
	}

	fmt.Println("--- Part One ---")
	discsOne := make([]disc, len(discs))
	copy(discsOne, discs)
	inverseTimeSkew(discsOne)
	fmt.Println(findFirstTime(discsOne))

	fmt.Println("--- Part Two ---")
	discsTwo := make([]disc, len(discs), len(discs) + 1)
	copy(discsTwo, discs)
	discsTwo = append(discsTwo, disc{ 11, 0 })
	inverseTimeSkew(discsTwo)
	fmt.Println(findFirstTime(discsTwo))
}

type disc struct {
	mod, pos int
}

func inverseTimeSkew(discs []disc) {
	for i := 0; i < len(discs); i++ {
		d := &discs[i]
		d.pos = (d.pos + i + 1) % d.mod
	}
}

func findFirstTime(discs []disc) int {
	for time := 0; ; time++ {
		fits := true
		for _, d := range discs {
			if d.pos != 0 {
				fits = false
				break
			}
		}

		if fits {
			return time
		}

		for i := 0; i < len(discs); i++ {
			d := &discs[i]
			d.pos++
			if d.pos == d.mod { d.pos = 0 }
		}
	}
}

func toInt(v string) int {
	i, e := strconv.Atoi(v)
	if e != nil { panic(e) }
	return i
}
