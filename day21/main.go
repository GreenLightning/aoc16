package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lines := make([]string, 0, 100)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	swapPositionRule := regexp.MustCompile(`swap position (\d+) with position (\d+)`)
	swapLetterRule := regexp.MustCompile(`swap letter (\w) with letter (\w)`)
	rotateAllRule := regexp.MustCompile(`rotate (left|right) (\d+) steps?`)
	rotatePositionRule := regexp.MustCompile(`rotate based on position of letter (\w)`)
	reverseRule := regexp.MustCompile(`reverse positions (\d+) through (\d+)`)
	moveRule := regexp.MustCompile(`move position (\d+) to position (\d+)`)

	{
		fmt.Println("--- Part One ---")

		pw := []byte("abcdefgh")

		for _, line := range lines {
			if result := swapPositionRule.FindStringSubmatch(line); result != nil {
				x, y := toInt(result[1]), toInt(result[2])
				pw[x], pw[y] = pw[y], pw[x]
			} else if result := swapLetterRule.FindStringSubmatch(line); result != nil {
				x, y := result[1][0], result[2][0]
				swapLetters(pw, x, y)
			} else if result := rotateAllRule.FindStringSubmatch(line); result != nil {
				amount := toInt(result[2])
				if result[1] == "left" {
					pw = rotateLeft(pw, amount)
				} else {
					pw = rotateRight(pw, amount)
				}
			} else if result := rotatePositionRule.FindStringSubmatch(line); result != nil {
				pw = rotateBasedOnPosition(pw, result[1])
			} else if result := reverseRule.FindStringSubmatch(line); result != nil {
				x, y := toInt(result[1]), toInt(result[2])
				reverse(pw, x, y)
			} else if result := moveRule.FindStringSubmatch(line); result != nil {
				x, y := toInt(result[1]), toInt(result[2])
				move(pw, x, y)
			} else {
				panic(line)
			}
		}

		fmt.Println(string(pw))
	}

	{
		fmt.Println("--- Part Two ---")

		pw := []byte("fbgdceah")

		for i := len(lines) - 1; i >= 0; i-- {
			line := lines[i]
			if result := swapPositionRule.FindStringSubmatch(line); result != nil {
				x, y := toInt(result[1]), toInt(result[2])
				pw[x], pw[y] = pw[y], pw[x]
			} else if result := swapLetterRule.FindStringSubmatch(line); result != nil {
				x, y := result[1][0], result[2][0]
				swapLetters(pw, x, y)
			} else if result := rotateAllRule.FindStringSubmatch(line); result != nil {
				amount := toInt(result[2])
				if result[1] == "left" {
					pw = rotateRight(pw, amount)
				} else {
					pw = rotateLeft(pw, amount)
				}
			} else if result := rotatePositionRule.FindStringSubmatch(line); result != nil {
				pw = inverseRotateBasedOnPosition(pw, result[1])
			} else if result := reverseRule.FindStringSubmatch(line); result != nil {
				x, y := toInt(result[1]), toInt(result[2])
				reverse(pw, x, y)
			} else if result := moveRule.FindStringSubmatch(line); result != nil {
				x, y := toInt(result[1]), toInt(result[2])
				move(pw, y, x)
			} else {
				panic(line)
			}
		}

		fmt.Println(string(pw))
	}
}

func swapLetters(data []byte, x, y byte) {
	for i := 0; i < len(data); i++ {
		if data[i] == x {
			data[i] = y
		} else if data[i] == y {
			data[i] = x
		}
	}
}

func rotateLeft(data []byte, amount int) []byte {
	amount %= len(data)
	result := make([]byte, len(data))
	for i := 0; i < len(data) - amount; i++ {
		result[i] = data[i + amount]
	}
	for i := 0; i < amount; i++ {
		result[len(data) - amount + i] = data[i]
	}
	return result
}

func rotateRight(data []byte, amount int) []byte {
	amount %= len(data)
	result := make([]byte, len(data))
	for i := 0; i < amount; i++ {
		result[i] = data[len(data) - amount + i]
	}
	for i := 0; i < len(data) - amount; i++ {
		result[i + amount] = data[i]
	}
	return result
}

func rotateBasedOnPosition(data []byte, letter string) []byte {
	index := strings.Index(string(data), letter)
	amount := 1 + index
	if index >= 4 { amount++ }
	return rotateRight(data, amount)
}

func inverseRotateBasedOnPosition(data []byte, letter string) []byte {
	for i := 0; i < len(data); i++ {
		candidate := rotateLeft(data, i)
		if bytes.Equal(rotateBasedOnPosition(candidate, letter), data) {
			return candidate
		}
	}
	return data
}

func reverse(data []byte, x, y int) {
	tmp := make([]byte, y+1 - x)
	copy(tmp, data[x:y+1])
	for i := 0; i < len(tmp); i++ {
		data[y - i] = tmp[i]
	}
}

func move(data []byte, x, y int) {
	tmp := data[x]
	for i := x; i + 1 < len(data); i++ {
		data[i] = data[i + 1]
	}
	for i := len(data) - 1; i > y; i-- {
		data[i] = data[i - 1]
	}
	data[y] = tmp
}

func toInt(v string) int {
	i, e := strconv.Atoi(v)
	if e != nil { panic(e) }
	return i
}
