package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	line := scanner.Text()
	maps := createMaps(len(line))
	updateMaps(maps, line)
	for scanner.Scan() {
		updateMaps(maps, scanner.Text())
	}

	mostCommonString  := ""
	leastCommonString := ""
	for i := 0; i < len(maps); i++ {
		most, least := analyze(maps[i])
		mostCommonString += string([]byte{ most })
		leastCommonString += string([]byte{ least })
	}
	fmt.Println(mostCommonString)
	fmt.Println(leastCommonString)
}

func createMaps(length int) []map[byte]int {
	result := make([]map[byte]int, length)
	for i := 0; i < length; i++ {
		result[i] = make(map[byte]int)
	}
	return result
}

func updateMaps(maps []map[byte]int, line string) {
	for i := 0; i < len(line); i++ {
		maps[i][line[i]]++
	}
}

func analyze(data map[byte]int) (byte, byte) {
	length := len(data)
	list := make([]Entry, length)
	index := 0
	for char, count := range data {
		list[index] = Entry{ char, count }
		index++
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Count < list[j].Count })
	return list[length-1].Character, list[0].Character
}

type Entry struct {
	Character byte
	Count int
}
