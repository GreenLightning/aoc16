package main

import (
	"bufio"
	"fmt"
	"os"
	"crypto/md5"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	input := scanner.Text()

	shortest := ""
	longest := ""

	current := make([]string, 0, 256)
	next := make([]string, 0, 256)

	current = append(current, "")

	for len(current) > 0 {
		for _, path := range current {
			x, y := walk(path)

			if x == 3 && y == 3 {
				if shortest == "" { shortest = path }
				longest = path
			} else {
				hash := md5.Sum([]byte(input + path))

				if y > 0 && (hash[0] >> 4) > 10 {
					next = append(next, path + "U")
				}
				if y < 3 && (hash[0] & 0xf) > 10 {
					next = append(next, path + "D")
				}
				if x > 0 && (hash[1] >> 4) > 10 {
					next = append(next, path + "L")
				}
				if x < 3 && (hash[1] & 0xf) > 10 {
					next = append(next, path + "R")
				}
			}
		}
		current, next = next, current[:0]
	}

	fmt.Println("--- Part One ---")
	fmt.Println(shortest)
	fmt.Println("--- Part Two ---")
	fmt.Println(len(longest))
}

func walk(path string) (int, int) {
	x, y := 0, 0
	for _, move := range path {
		switch move {
			case 'U': y--
			case 'D': y++
			case 'L': x--
			case 'R': x++
		}
	}
	return x, y
}
