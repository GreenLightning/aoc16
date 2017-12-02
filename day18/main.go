package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	input := scanner.Text()

	row := "." + input + "."

	fmt.Println("--- Part One ---")
	fmt.Println(countSafeTilesInRows(row, 40))
	fmt.Println("--- Part Two ---")
	fmt.Println(countSafeTilesInRows(row, 400000))
}

func countSafeTilesInRows(row string, limit int) int {
	safeTiles := countSafeTiles(row)
	for i := 0; i < limit-1; i++ {
		row = update(row)
		safeTiles += countSafeTiles(row)
	}
	return safeTiles
}

func countSafeTiles(row string) int {
	count := 0
	for i := 1; i < len(row)-1; i++ {
		if row[i] == '.' {
			count++
		}
	}
	return count
}

func update(row string) string {
	res := []byte(strings.Repeat(".", len(row)))
	for i := 1; i < len(row)-1; i++ {
		if row[i-1] != row[i+1] {
			res[i] = '^'
		}
	}
	return string(res)
}
