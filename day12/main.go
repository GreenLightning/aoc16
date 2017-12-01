package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	data := make([]bool, width * height)

	for scanner.Scan() {
		line := scanner.Text()
	}
}

const (
	instrTypeCpy = iota
	instrTypeInc
	instrTypeDec
	instrTypeJnz
)

const (
	argTypeRegA = iota
	argTypeRegB
	argTypeRegC
	argTypeRegD
	argTypeVal
)

type instruction struct {
	instrType int
	xType int
	xVal int
	yType int
	yVal int
}

type machine struct {
	a, b, c, d int
}

func (m *machine) execute(i instruction) {
	switch (i.instrType) {
	case instrTypeCpy:
		*m.getReg(i.yType) = m.getVal(i.xType, i.xVal)

	}
}

func (m *machine) getReg(argType int) *int {
	switch (argType) {
		case argTypeRegA: return m.a
		case argTypeRegB: return m.b
		case argTypeRegC: return m.c
		case argTypeRegD: return m.d
		default: return nil
	}
}

func (m *machine) getVal(argType int, argVal int) int {
	if (argType == argTypeVal) {
		return argVal
	} else {
		return *m.getReg(argType)
	}
}
