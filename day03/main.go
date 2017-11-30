package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	possibleByRow := 0
	possibleByCol := 0

	for scanner.Scan() {
		line := scanner.Text()
		a1, b1, c1 := splitLine(line)

		if !scanner.Scan() { break }
		line = scanner.Text()
		a2, b2, c2 := splitLine(line)

		if !scanner.Scan() { break }
		line = scanner.Text()
		a3, b3, c3 := splitLine(line)

		if check(a1, b1, c1) { possibleByRow++ }
		if check(a2, b2, c2) { possibleByRow++ }
		if check(a3, b3, c3) { possibleByRow++ }

		if check(a1, a2, a3) { possibleByCol++ }
		if check(b1, b2, b3) { possibleByCol++ }
		if check(c1, c2, c3) { possibleByCol++ }
	}

	fmt.Println("Possible by row:", possibleByRow)
	fmt.Println("Possible by col:", possibleByCol)
}

func splitLine(line string) (int, int, int) {
	a, line := grabInt(line)
	b, line := grabInt(line)
	c, line := grabInt(line)
	return a, b, c
}

func grabInt(line string) (int, string) {
	line = strings.TrimSpace(line)
	index := strings.Index(line, " ")
	if index == -1 { index = len(line) }
	i, err := strconv.Atoi(line[:index])
	if err != nil { panic(err) }
	return i, line[index:]
}

func check(x, y, z int) bool {
	return x + y > z && y + z > x && z + x > y
}
