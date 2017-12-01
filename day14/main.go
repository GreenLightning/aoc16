package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"crypto/md5"
)

const hashSize = md5.Size
const splitSize = 2 * hashSize

type hashValue [hashSize]byte
type splitValue [splitSize]byte

type hashFunction func([]byte)hashValue

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	input := scanner.Text()

	fmt.Println("--- Part One ---")
	streamOne := newHashStream(input, hashFunctionOne)
	searchForKeys(streamOne)

	fmt.Println("--- Part Two ---")
	streamTwo := newHashStream(input, hashFunctionTwo)
	searchForKeys(streamTwo)
}

func hashFunctionOne(data []byte) hashValue {
	return md5.Sum(data)
}

func hashFunctionTwo(data []byte) hashValue {
	hash := md5.Sum(data)
	for i := 0; i < 2016; i++ {
		hash = md5.Sum([]byte(fmt.Sprintf("%x", hash)))
	}
	return hash
}

func splitHash(hash hashValue) splitValue {
	split := splitValue{}
	for i := 0; i < hashSize; i++ {
		split[2*i + 0] = hash[i] >> 4
		split[2*i + 1] = hash[i] & 0xf
	}
	return split
}

func searchForKeys(stream *hashStream) {
	keys := 0
	triples := 0
	fmt.Printf("\rFound %d keys and %d triples", keys, triples)
	for {
		hex := splitHash(stream.get(0))
		found, value := hasTriple(hex)
		if  found {
			triples++
			fmt.Printf("\rFound %d keys and %d triples", keys, triples)
			if findQuintuple(stream, value) {
				keys++
				fmt.Printf("\rFound %d keys and %d triples", keys, triples)
				if keys == 64 {
					break
				}
			}
		}
		stream.advance()
	}
	newMessage     := fmt.Sprintf("\rThe Answer is %d", stream.index)
	currentMessage := fmt.Sprintf("\rFound %d keys and %d triples", keys, triples)
	newMessage += strings.Repeat(" ", len(currentMessage) - len(newMessage))
	fmt.Println(newMessage)
}

func hasTriple(hex splitValue) (bool, byte) {
	for i := 0; i+2 < splitSize; i++ {
		if hex[i] == hex[i+1] && hex[i] == hex[i+2] {
			return true, hex[i]
		}
	}
	return false, 0
}

func findQuintuple(stream *hashStream, target byte) bool {
	for next := 1; next <= 1000; next++ {
		hex := splitHash(stream.get(next))
		for i := 0; i+4 < splitSize; i++ {
			if hex[i+0] == target &&
			   hex[i+1] == target &&
			   hex[i+2] == target &&
			   hex[i+3] == target &&
			   hex[i+4] == target {
				return true
			}
		}
	}
	return false
}

type hashStream struct {
	salt string
	index int
	function hashFunction
	values []hashValue
}

func newHashStream(salt string, function hashFunction) *hashStream {
	return &hashStream{ salt, 0, function, nil }
}

func (stream *hashStream) get(offset int) hashValue {
	for offset >= len(stream.values) {
		source := stream.salt + strconv.Itoa(stream.index + len(stream.values))
		stream.values = append(stream.values, stream.function([]byte(source)))
	}
	return stream.values[offset]
}

func (stream *hashStream) advance() {
	stream.index++
	if len(stream.values) > 0 {
		stream.values = stream.values[1:]
	}
}
