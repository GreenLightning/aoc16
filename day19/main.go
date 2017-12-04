package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

	length := toInt(input)

	{
		fmt.Println("--- Part One ---")
		elves := make([]int, length)

		for i := 0; i < length; i++ {
			elves[i] = 1
		}

		i := 0
		for elves[i] != length {
			index := (i + elves[i]) % length
			elves[i] += elves[index]
			elves[index] = 0
			i = (i + elves[i]) % length
		}

		fmt.Println(i + 1)
	}

	{
		fmt.Println("--- Part Two ---")
		elves := bucketList{}

		for i := 0; i < length; i++ {
			elves.append(i + 1)
		}

		i := 0
		for elves.length() > 1 {
			for iteration := 0; elves.length() > 1 && iteration < 1024; iteration++ {
				index := (i + elves.length() / 2) % elves.length()
				elves.removeIndex(index)
				if i < index { i++ }
				i %= elves.length()
			}
			fmt.Printf("\r%6.2f%%", (1.0 - float32(elves.length()) / float32(length)) * 100.0)
		}
		message := fmt.Sprint(elves.getIndex(0))
		if len(message) < 7 { message += strings.Repeat(" ", 7 - len(message)) }
		fmt.Printf("\r%s\n", message)
	}
}

const maxBucketSize = 1024

type bucketList struct {
	start, end *bucket
	total int
}

type bucket struct {
	previous, next *bucket
	data []int
}

func (list *bucketList) length() int {
	return list.total
}

func (list *bucketList) append(value int) {
	if list.end == nil {
		b := &bucket{ nil, nil, make([]int, 1, maxBucketSize) }
		b.data[0] = value
		list.start, list.end = b, b
	} else if len(list.end.data) == maxBucketSize {
		b := &bucket{ list.end, nil, make([]int, 1, maxBucketSize) }
		b.data[0] = value
		list.end.next = b
		list.end = b
	} else {
		list.end.data = append(list.end.data, value)
	}
	list.total++
}

func (list *bucketList) removeIndex(index int) {
	b := list.start
	for index >= len(b.data) {
		index -= len(b.data)
		b = b.next
	}

	for i, n := index, len(b.data)-1; i < n; i++ {
		b.data[i] = b.data[i+1]
	}
	b.data = b.data[:len(b.data)-1]
	list.total--

	if b.previous != nil && len(b.data) <= maxBucketSize - len(b.previous.data) {
		copyBucketData(b.previous, b)
		b.previous.next = b.next
		if b.next != nil { b.next.previous = b.previous }
		if b == list.end { list.end = b.previous }
	} else if b.next != nil && len(b.data) <= maxBucketSize - len(b.next.data) {
		copyBucketData(b.next, b)
		b.next.previous = b.previous
		if b.previous != nil { b.previous.next = b.next }
		if b == list.start { list.start = b.next }
	}
}

func copyBucketData(dest, source *bucket) {
	destLength, sourceLength := len(dest.data), len(source.data)
	dest.data = dest.data[:destLength + sourceLength]
	for i := 0; i < sourceLength; i++ {
		dest.data[destLength + i] = source.data[i]
	}
}

func (list *bucketList) getIndex(index int) int {
	b := list.start
	for index >= len(b.data) {
		index -= len(b.data)
		b = b.next
	}
	return b.data[index]
}

func toInt(v string) int {
	i, e := strconv.Atoi(v)
	if e != nil { panic(e) }
	return i
}
