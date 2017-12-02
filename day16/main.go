package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	input := scanner.Text()

	data := make([]bool, len(input))
	for i, c := range input {
		if c == '1' {
			data[i] = true
		} else {
			data[i] = false
		}
	}

	fmt.Println("--- Part One ---")
	fmt.Println(generateDataAndCalculateChecksum(data, 272))
	fmt.Println("--- Part Two ---")
	fmt.Println(generateDataAndCalculateChecksum(data, 35651584))
}

func generateDataAndCalculateChecksum(base []bool, target int) string {
	data := make([]bool, len(base))
	copy(data, base)

	for len(data) < target {
		next := make([]bool, 2*len(data) + 1)
		copy(next, data)
		next[len(data)] = false
		for i:= 0; i < len(data); i++ {
			next[2*len(data)-i] = !data[i]
		}
		data = next
	}

	data = data[0:target]

	for len(data) % 2 == 0 {
		next := make([]bool, len(data)/2)
		for i := 0; i < len(next); i++ {
			next[i] = (data[2*i+0] == data[2*i+1])
		}
		data = next
	}

	output := make([]byte, len(data))
	for i, d := range data {
		if d {
			output[i] = '1'
		} else {
			output[i] = '0'
		}
	}

	return string(output)
}
