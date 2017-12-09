package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	reg := regexp.MustCompile(`([-a-z]+)(\d+)\[(.*)\]`)

	sum := 0

	for scanner.Scan() {
		result := reg.FindStringSubmatch(scanner.Text())
		name := result[1]
		id := result[2]
		checksum := result[3]

		if checksum == calculateChecksum(name) {
			val, err := strconv.Atoi(id)
			if err != nil { panic(err) }
			sum += val
			fmt.Printf("%d: %s\n", val, decrypt(name, val))
		}
	}

	fmt.Println(sum)
}

func calculateChecksum(name string) string {
	counts := make(map[string]int)
	for i := 0; i < len(name); i++ {
		char := name[i:i+1]
		if char != "-" {
			counts[char]++
		}
	}

	length := len(counts)
	list := make([]Entry, length)
	index := 0
	for char, count := range counts {
		list[index] = Entry{ char, count }
		index++
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Count != list[j].Count {
			return list[i].Count > list[j].Count
		}
		return list[i].Character < list[j].Character
	})

	result := ""
	for i := 0; i < 5 && i < length; i++ {
		result += list[i].Character
	}
	return result
}

type Entry struct {
	Character string
	Count int
}

func decrypt(name string, id int) string {
	result := ""
	for i := 0; i < len(name); i++ {
		char := name[i:i+1]
		if char == "-" {
			result += " "
		} else {
			b := int('a') + (int(char[0]) - int('a') + id) % 26
			result += string([]byte{ byte(b) })
		}
	}
	return strings.TrimSpace(result)
}
