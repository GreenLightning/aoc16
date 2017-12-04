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

	instructions := make([]instruction, 0)

	for scanner.Scan() {
		line := scanner.Text()
		i := parseInstruction(line)
		instructions = append(instructions, i)
	}

	fmt.Println("--- Part One ---")
	m := machine{}
	m.a = 7
	m.run(instructions)
	fmt.Println(m.a)

	fmt.Println("--- Part Two ---")
	m = machine{}
	m.a = 12
	m.run(instructions)
	fmt.Println(m.a)
}

const (
	argTypeRegA = iota
	argTypeRegB
	argTypeRegC
	argTypeRegD
	argTypeVal
)

const (
	instrTypeCpy = iota
	instrTypeInc
	instrTypeDec
	instrTypeJnz
	instrTypeTgl
)

type argument struct {
	aType int
	val int
}

type instruction struct {
	iType int
	x argument
	y argument
}

func split(line string) (string, string) {
	index := strings.Index(line, " ")
	if index == -1 {
		return line, ""
	} else {
		return line[:index], line[index+1:]
	}
}

func parseArg(line string) (argument, string) {
	value, line := split(line)
	switch value {
		case "a": return argument{ argTypeRegA, 0 }, line
		case "b": return argument{ argTypeRegB, 0 }, line
		case "c": return argument{ argTypeRegC, 0 }, line
		case "d": return argument{ argTypeRegD, 0 }, line
		default:  return argument{ argTypeVal, toInt(value) }, line
	}
}

func parseRegArg(line string) (argument, string) {
	arg, line := parseArg(line)
	if arg.aType == argTypeVal {
		panic("expected register put found literal value")
	}
	return arg, line
}

func parseInstruction(line string) instruction {
	if strings.HasPrefix(line, "cpy ") {
		line = line[4:]
		x, line := parseArg(line)
		y, _ := parseRegArg(line)
		return instruction{ instrTypeCpy, x, y }
	} else if strings.HasPrefix(line, "inc ") {
		line = line[4:]
		x, _ := parseRegArg(line)
		return instruction{ instrTypeInc, x, argument{} }
	} else if strings.HasPrefix(line, "dec ") {
		line = line[4:]
		x, _ := parseRegArg(line)
		return instruction{ instrTypeDec, x, argument{} }
	} else if strings.HasPrefix(line, "jnz ") {
		line = line[4:]
		x, line := parseArg(line)
		y, _ := parseArg(line)
		return instruction{ instrTypeJnz, x, y }
	} else if strings.HasPrefix(line, "tgl ") {
		line = line[4:]
		x, _ := parseArg(line)
		return instruction{ instrTypeTgl, x, argument{} }
	} else {
		panic(fmt.Sprint("cannot parse instruction: ", line))
	}
}

type machine struct {
	a, b, c, d int
}

func (m *machine) run(instructions []instruction) {
	tmp := make([]instruction, len(instructions))
	copy(tmp, instructions)
	ip := 0
	for ip >= 0 && ip < len(tmp) {
		ip += m.execute(tmp, ip)
	}
}

func (m *machine) execute(instructions []instruction, ip int) int {
	i := instructions[ip]
	switch i.iType {
		case instrTypeCpy:
			if reg := m.getReg(i.y); reg != nil {
				*reg = m.getVal(i.x)
			}
			return 1
		case instrTypeInc:
			if reg := m.getReg(i.x); reg != nil {
				*reg = *reg + 1
			}
			return 1
		case instrTypeDec:
			if reg := m.getReg(i.x); reg != nil {
				*reg = *reg - 1
			}
			return 1
		case instrTypeJnz:
			result := m.getVal(i.x)
			if result != 0 {
				return m.getVal(i.y)
			} else {
				return 1
			}
		case instrTypeTgl:
			index := ip + m.getVal(i.x)
			if index >= 0 && index < len(instructions) {
				switch instructions[index].iType {
					case instrTypeCpy: instructions[index].iType = instrTypeJnz
					case instrTypeInc: instructions[index].iType = instrTypeDec
					case instrTypeDec: instructions[index].iType = instrTypeInc
					case instrTypeJnz: instructions[index].iType = instrTypeCpy
					case instrTypeTgl: instructions[index].iType = instrTypeInc
				}
			}
			return 1
		default: panic(fmt.Sprintf("unknown instruction type %d", i.iType))
	}
}

func (m *machine) getReg(arg argument) *int {
	switch arg.aType {
		case argTypeRegA: return &m.a
		case argTypeRegB: return &m.b
		case argTypeRegC: return &m.c
		case argTypeRegD: return &m.d
		default: return nil
	}
}

func (m *machine) getVal(arg argument) int {
	if (arg.aType == argTypeVal) {
		return arg.val
	} else {
		return *m.getReg(arg)
	}
}

func toInt(v string) int {
	i, e := strconv.Atoi(v)
	if e != nil { panic(e) }
	return i
}
