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

	scanner.Scan()
	input := scanner.Text()

	fmt.Println(getDecompressedLength(input, false))
	fmt.Println(getDecompressedLength(input, true))
}

func getDecompressedLength(input string, recursive bool) int {
	output := 0
	marker := regexp.MustCompile(`\((\d+)x(\d+)\)`)
	pos := marker.FindStringIndex(input)
	for pos != nil {
		data := marker.FindStringSubmatch(input[pos[0]:pos[1]])
		sourceLength := toInt(data[1])
		repeats := toInt(data[2])
		decompressedLength := sourceLength
		if recursive {
			decompressedLength = getDecompressedLength(input[pos[1]:pos[1]+sourceLength], recursive)
		}
		output += pos[0] + repeats * decompressedLength
		input = input[pos[1]+sourceLength:]
		pos = marker.FindStringIndex(input)
	}
	output += len(input)
	return output
}

func toInt(v string) int {
	i, e := strconv.Atoi(v)
	if e != nil { panic(e) }
	return i
}
