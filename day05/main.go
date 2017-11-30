package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"crypto/md5"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <input>\n", os.Args[0])
		return
	}

	input := os.Args[1]

	{ // --- Part One ---
		fmt.Print("Password: ")
		index := 0
		for i := 0; i < 8; i++ {
			var hash [md5.Size]byte
			for {
				value := input + strconv.Itoa(index)
				hash = md5.Sum([]byte(value))
				if hash[0] == 0 && hash[1] == 0 && (hash[2] & 0xf0) == 0 { break }
				index++
			}
			char := fmt.Sprintf("%x", hash[2])
			fmt.Print(char)
			index++
		}
		fmt.Println()
	}

	{ // --- Part Two ---
		index := 0
		password := []byte(strings.Repeat("_", 8))
		for i := 0; i < 8; {
			fmt.Print("\rPassword: ", string(password))
			var hash [md5.Size]byte
			for {
				value := input + strconv.Itoa(index)
				hash = md5.Sum([]byte(value))
				if hash[0] == 0 && hash[1] == 0 && (hash[2] & 0xf0) == 0 { break }
				index++
			}
			position := hash[2]
			value := hash[3] >> 4
			if position >= 0 && position < 8 && password[position] == '_' {
				password[position] = fmt.Sprintf("%x", value)[0]
				i++
			}
			index++
		}
		fmt.Println("\rPassword:", string(password))
	}
}
