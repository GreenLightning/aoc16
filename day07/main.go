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

	tlsCount := 0
	sslCount := 0

	for scanner.Scan() {
		address := scanner.Text()
		if supportsTls(address) {
			tlsCount++
		}
		if supportsSsl(address) {
			sslCount++
		}
	}

	fmt.Println(tlsCount)
	fmt.Println(sslCount)
}

func supportsTls(address string) bool {
	hasAbba := false
	hasAbbaInHypernetSequence := false

	inHypernetSequence := false

	for i := 0; i + 3 < len(address); i++ {
		char := address[i]
		if char == '[' {
			inHypernetSequence = true
		} else if char == ']' {
			inHypernetSequence = false
		} else if address[i] != address[i+1] {
			if address[i] == address[i+3] && address[i+1] == address[i+2] {
				if inHypernetSequence {
					hasAbbaInHypernetSequence = true
				} else {
					hasAbba = true
				}
			}
		}
	}

	return hasAbba && !hasAbbaInHypernetSequence
}

func supportsSsl(address string) bool {
	inHypernetSequence := false
	for i := 0; i + 2 < len(address); i++ {
		char := address[i]
		if char == '[' {
			inHypernetSequence = true
		} else if char == ']' {
			inHypernetSequence = false
		} else if address[i] != address[i+1] && address[i] == address[i+2] {
			if !inHypernetSequence && hasBab(address, address[i], address[i+1]) {
				return true
			}
		}
	}
	return false
}

func hasBab(address string, a, b byte) bool {
	inHypernetSequence := false
	for i := 0; i + 2 < len(address); i++ {
		char := address[i]
		if char == '[' {
			inHypernetSequence = true
		} else if char == ']' {
			inHypernetSequence = false
		} else if inHypernetSequence && address[i] == b && address[i+1] == a && address[i+2] == b {
			return true
		}
	}
	return false
}
