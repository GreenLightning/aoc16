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

	simpleLayout :=
		"     " +
		" 123 " +
		" 456 " +
		" 789 " +
		"     "

	complexLayout :=
		"       " +
		"   1   " +
		"  234  " +
		" 56789 " +
		"  ABC  " +
		"   D   " +
		"       "

	simplePad := keypad{ simpleLayout, 5, 5, 3, 3 }
	complexPad := keypad{ complexLayout, 7, 7, 2, 4 }

	simplePos := simplePad.start()
	complexPos := complexPad.start()

	simpleCode := ""
	complexCode := ""

	for scanner.Scan() {
		line := scanner.Text()
		for i := 0; i < len(line); i++ {
			dir := line[i:i+1]
			simplePos.move(dir)
			complexPos.move(dir)
		}
		simpleCode += simplePos.value()
		complexCode += complexPos.value()
	}
	fmt.Println("Simple Keypad: ", simpleCode)
	fmt.Println("Complex Keypad:", complexCode)
}

type keypad struct {
	layout string
	w, h int
	sx, sy int
}

func (k *keypad) start() position {
	return position{ k, k.sx, k.sy }
}

type position struct {
	keypad *keypad
	x, y int
}

func (p *position) move(dir string) {
	x, y := p.x, p.y
	switch dir {
		case "U": y -= 1;
		case "D": y += 1;
		case "L": x -= 1;
		case "R": x += 1;
	}
	if p.valueOf(x, y) != " " {
		p.x, p.y = x, y
	}
}

func (p *position) value() string {
	return p.valueOf(p.x, p.y)
}

func (p *position) valueOf(x, y int) string {
	index := p.keypad.w * y + x
	return p.keypad.layout[index:index+1]
}
