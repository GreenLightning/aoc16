package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

func main() {
	initial, count := parseFile("input.txt")

	{
		fmt.Println("--- Part One ---")
		target := makeTargetState(count)
		f := newFinder(initial, target, count)
		steps, found := f.find()
		if found { fmt.Println(steps) }
	}

	{
		fmt.Println("--- Part Two ---")
		initial.pairs = append(initial.pairs, pair{ 0, 0 })
		initial.pairs = append(initial.pairs, pair{ 0, 0 })
		count += 2
		target := makeTargetState(count)
		f := newFinder(initial, target, count)
		steps, found := f.find()
		if found { fmt.Println(steps) }
	}
}

func parseFile(name string) (state, int) {
	var result state

	file, err := os.Open(name)
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lineRule := regexp.MustCompile(`The \w+ floor contains (?:nothing relevant\.|(.*))`)
	sepRule  := regexp.MustCompile(`,? and |, |\.`)
	descRule := regexp.MustCompile(`a (\w+) generator|a (\w+)-compatible microchip`)

	positions := make(map[string]int)

	floor := byte(0)

	for scanner.Scan() {
		line := scanner.Text()
		matches := lineRule.FindStringSubmatch(line)
		descriptions := matches[1]
		for descriptions != "" {
			indices := sepRule.FindStringIndex(descriptions)
			description := descriptions[:indices[0]]
			descriptions = descriptions[indices[1]:]

			matches := descRule.FindStringSubmatch(description)
			name, isGenerator := matches[1], true
			if matches[1] == "" {
				name, isGenerator = matches[2], false
			}

			pos, ok := positions[name]
			if !ok {
				pos = len(positions)
				positions[name] = pos
				result.pairs = append(result.pairs, pair{ 0, 0 })
				// We pack states into 32 bits, which is enough room for the
				// elevator (2 bits) and 7 pairs (4 bits each). We need space
				// for 2 additional pairs for part two, which means the input
				// must not have more than 5 pairs.
				if pos >= 5 { panic("ERROR: Input too big!") }
			}

			if isGenerator {
				result.pairs[pos].generator = floor
			} else {
				result.pairs[pos].chip = floor
			}
		}
		floor++
	}

	return result, len(positions)
}

type state struct {
	elevator byte
	pairs    []pair
}

type pair struct {
	chip, generator byte
}

func (s *state) sort() {
	sort.Slice(s.pairs, func(i, j int) bool {
		if s.pairs[i].chip != s.pairs[j].chip {
			return s.pairs[i].chip < s.pairs[j].chip
		}
		return s.pairs[i].generator < s.pairs[j].generator
	})
}

func makeTargetState(count int) state {
	result := state{ 3, make([]pair, count) }
	for i := 0; i < count; i++ {
		result.pairs[i] = pair{ 3, 3 }
	}
	return result
}

type packed uint32

func pack(s state, count int) packed {
	result := packed(s.elevator)
	for i := 0; i < count; i++ {
		pair := s.pairs[i]
		result = (result << 4) | packed(pair.chip << 2) | packed(pair.generator)
	}
	return result
}

func unpack(p packed, count int) state {
	result := state{ 0, make([]pair, count) }
	for i := 0; i < count; i++ {
		result.pairs[i].chip = byte((p & 0xc) >> 2)
		result.pairs[i].generator = byte(p & 0x3)
		p >>= 4
	}
	result.elevator = byte(p)
	return result
}

type finder struct {
	steps int
	count int
	target packed
	checked map[packed]bool
	current []packed
	next []packed
}

func newFinder(initial, target state, count int) *finder {
	initial.sort()
	pi := pack(initial, count)
	target.sort()
	ti := pack(target, count)
	f := &finder{}
	f.count = count
	f.target = ti
	f.checked = make(map[packed]bool)
	f.checked[pi] = true
	f.current = []packed{ pi }
	return f
}

func (f *finder) find() (int, bool) {
	for len(f.current) > 0 {
		for _, p := range f.current {
			if p == f.target {
				return f.steps, true
			}
			s := unpack(p, f.count)
			if s.elevator > 0 { f.updateAll(s, s.elevator, s.elevator - 1) }
			if s.elevator < 3 { f.updateAll(s, s.elevator, s.elevator + 1) }
		}
		f.current, f.next = f.next, f.current
		f.next = f.next[:0]
		f.steps++
	}
	return 0, false
}

func (f* finder) updateAll(s state, floor, newFloor byte) {
	s.elevator = newFloor

	// move one or two chips
	for i := 0; i < f.count; i++ {
		for j := 0; j < f.count; j++ {
			if s.pairs[i].chip == floor && s.pairs[j].chip == floor {
				s.pairs[i].chip = newFloor
				s.pairs[j].chip = newFloor
				f.update(s)
				s.pairs[i].chip = floor
				s.pairs[j].chip = floor
			}
		}
	}
	// move one or two generators
	for i := 0; i < f.count; i++ {
		for j := 0; j < f.count; j++ {
			if s.pairs[i].generator == floor && s.pairs[j].generator == floor {
				s.pairs[i].generator = newFloor
				s.pairs[j].generator = newFloor
				f.update(s)
				s.pairs[i].generator = floor
				s.pairs[j].generator = floor
			}
		}
	}
	// move a chip and a generator together
	for i := 0; i < f.count; i++ {
		if s.pairs[i].chip == floor && s.pairs[i].generator == floor {
			s.pairs[i].chip = newFloor
			s.pairs[i].generator = newFloor
			f.update(s)
			s.pairs[i].chip = floor
			s.pairs[i].generator = floor
		}
	}

	s.elevator = floor
}

func (f *finder) update(s state) {
	// Check if a lone chip is on the same floor as another generator, which
	// fries the chip and makes the state invalid.
	for i := 0; i < f.count; i++ {
		if s.pairs[i].chip != s.pairs[i].generator {
			for j := 0; j < f.count; j++ {
				if s.pairs[i].chip == s.pairs[j].generator {
					return
				}
			}
		}
	}

	sCopy := state{ s.elevator, make([]pair, f.count) }
	copy(sCopy.pairs, s.pairs)
	sCopy.sort()
	p := pack(sCopy, f.count)
	if !f.checked[p] {
		f.checked[p] = true
		f.next = append(f.next, p)
	}
}
