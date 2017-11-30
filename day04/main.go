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
	list := make(EntryList, length)
	index := 0
	for char, count := range counts {
		list[index] = Entry{ char, count }
		index++
	}
	sort.Sort(list)

	result := ""
	for i := 0; i < 5 && i < length; i++ {
		result += list[i].Character
	}
	return result
}

type EntryList []Entry

type Entry struct {
	Character string
	Count int
}

func (p *EntryList) Len() int {
	return len(p)
}

func (p *EntryList) Less(i, j int) bool {
	if p[i].Count != p[j].Count {
		return p[i].Count > p[j].Count
	}
	return p[i].Character < p[j].Character
}

func (p *EntryList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
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
