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
			for iteration := 0; elves.length() > 1 && iteration < 4096; iteration++ {
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

const maxBucketSize = 2048

type bucketList struct {
	buckets []bucket
	total int
}

type bucket struct {
	data []int
}

func (list *bucketList) length() int {
	return list.total
}

func (list *bucketList) getIndex(index int) int {
	b, _, index := list.getBucket(index)
	return b.data[index]
}

func (list *bucketList) append(value int) {
	if len(list.buckets) == 0 || len(list.getLastBucket().data) == maxBucketSize {
		b := bucket{ make([]int, 1, maxBucketSize) }
		b.data[0] = value
		list.buckets = append(list.buckets, b)
	} else {
		b := list.getLastBucket()
		b.data = append(b.data, value)
	}
	list.total++
}

func (list *bucketList) removeIndex(index int) {
	b, bi, index := list.getBucket(index)
	n := len(b.data) - 1
	for i := index; i < n; i++ {
		b.data[i] = b.data[i+1]
	}
	b.data = b.data[:n]
	list.total--

	if bi-1 >= 0 && n <= maxBucketSize - len(list.buckets[bi-1].data) {
		list.copyBucketData(&list.buckets[bi-1], b)
		list.removeBucket(bi)
	} else if bi+1 < len(list.buckets) && n <= maxBucketSize - len(list.buckets[bi+1].data) {
		list.copyBucketData(b, &list.buckets[bi+1])
		list.removeBucket(bi+1)
	}
}

func (list *bucketList) getLastBucket() *bucket {
	length := len(list.buckets)
	return &list.buckets[length-1]
}

func (list *bucketList) getBucket(index int) (*bucket, int, int) {
	bi := 0
	for index >= len(list.buckets[bi].data) {
		index -= len(list.buckets[bi].data)
		bi++
	}
	return &list.buckets[bi], bi, index
}

func (list *bucketList) copyBucketData(dest, source *bucket) {
	destLength, sourceLength := len(dest.data), len(source.data)
	dest.data = dest.data[:destLength + sourceLength]
	for i := 0; i < sourceLength; i++ {
		dest.data[destLength + i] = source.data[i]
	}
}

func (list *bucketList) removeBucket(index int) {
	n := len(list.buckets) - 1
	for bi := index; bi < n; bi++ {
		list.buckets[bi] = list.buckets[bi+1]
	}
	list.buckets = list.buckets[:n]
}

func toInt(v string) int {
	i, e := strconv.Atoi(v)
	if e != nil { panic(e) }
	return i
}
