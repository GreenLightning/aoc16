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
	tmp := make([]byte, splitSize)
	for i := 0; i < 2016; i++ {
		hash = md5.Sum(splitAndConvertToHex(hash, tmp))
	}
	return hash
}

const hextable = "0123456789abcdef"

func splitAndConvertToHex(hash hashValue, result []byte) []byte {
	for i := 0; i < hashSize; i++ {
		result[2*i + 0] = hextable[hash[i] >> 4]
		result[2*i + 1] = hextable[hash[i] & 0xf]
	}
	return result
}

func searchForKeys(stream *hashStream) {
	keys := 0
	triples := 0
	fmt.Printf("\rFound %d keys and %d triples", keys, triples)
	for {
		split := splitHash(stream.get(0))
		found, value := hasTriple(split)
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

func hasTriple(split splitValue) (bool, byte) {
	for i := 0; i+2 < splitSize; i++ {
		if split[i] == split[i+1] && split[i] == split[i+2] {
			return true, split[i]
		}
	}
	return false, 0
}

func findQuintuple(stream *hashStream, target byte) bool {
	for next := 1; next <= 1000; next++ {
		split := splitHash(stream.get(next))
		for i := 0; i+4 < splitSize; i++ {
			if split[i+0] == target &&
			   split[i+1] == target &&
			   split[i+2] == target &&
			   split[i+3] == target &&
			   split[i+4] == target {
				return true
			}
		}
	}
	return false
}

func splitHash(hash hashValue) splitValue {
	split := splitValue{}
	for i := 0; i < hashSize; i++ {
		split[2*i + 0] = hash[i] >> 4
		split[2*i + 1] = hash[i] & 0xf
	}
	return split
}

type hashStream struct {
	salt string
	index int
	function hashFunction
	values circularBuffer
}

func newHashStream(salt string, function hashFunction) *hashStream {
	return &hashStream{ salt, 0, function, makeCircularBuffer() }
}

func (stream *hashStream) get(offset int) hashValue {
	for offset >= stream.values.length {
		source := stream.salt + strconv.Itoa(stream.index + stream.values.length)
		stream.values.append(stream.function([]byte(source)))
	}
	return stream.values.get(offset)
}

func (stream *hashStream) advance() {
	stream.index++
	stream.values.advance()
}

type circularBuffer struct {
	start, length int
	data []hashValue
}

func makeCircularBuffer() circularBuffer {
	return circularBuffer{ 0, 0, make([]hashValue, 0x400) }
}

func (buffer * circularBuffer) get(index int) hashValue {
	pos := (buffer.start + index) & 0x3ff
	return buffer.data[pos]
}

func (buffer *circularBuffer) append(value hashValue) {
	pos := (buffer.start + buffer.length) & 0x3ff
	buffer.data[pos] = value
	buffer.length++
}

func (buffer *circularBuffer) advance() {
	if buffer.length > 0 {
		buffer.start = (buffer.start + 1) & 0x3ff
		buffer.length--
	}
}
