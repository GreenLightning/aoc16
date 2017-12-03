package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const MaxUint32 = ^uint32(0)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var list blacklist

	for scanner.Scan() {
		line := scanner.Text()
		index := strings.Index(line, "-")
		low := toUint32(line[:index])
		high := toUint32(line[index+1:])
		list.add(low, high)
	}

	fmt.Println("--- Part One ---")
	if len(list.blocked) == 0 || list.blocked[0].low > 0 {
		fmt.Println(0)
	} else {
		fmt.Println(list.blocked[0].high + 1)
	}

	fmt.Println("--- Part Two ---")
	allowed := uint32(0)
	if len(list.blocked) > 0 {
		allowed += list.blocked[0].low
		allowed += MaxUint32 - list.blocked[len(list.blocked)-1].high
	}
	for i := 0; i + 1 < len(list.blocked); i++ {
		allowed += list.blocked[i+1].low - list.blocked[i].high - 1
	}
	fmt.Println(allowed)
}

type interval struct {
	low, high uint32
}

type blacklist struct {
	blocked []interval
}

func (list *blacklist) add(low, high uint32) {
	// Target is the index we want to insert the new interval into,
	// based on the low boundary, in the range [0, len(list.blocked)].
	target := 0
	for target < len(list.blocked) && list.blocked[target].low <= low {
		target++
	}

	if target > 0 && connected(list.blocked[target - 1].high, low) {
		// The new interval starts inside or directly after an anlready existing
		// interval.
		if list.blocked[target - 1].high >= high {
			// If the new interval is completely enclosed by the already
			// existing interval, we don't have to update anything.
			return
		} else {
			// Else, the new interval expands the already existing interval, so
			// we update the bounds of the already existing interval. Also, we
			// have to adjust target to point to the updated interval.
			list.blocked[target - 1].high = high
			target--
		}
	} else {
		// The new interval starts outside of an already existing interval, so
		// we extend the blocked list, shift the other intervals up and put our
		// interval into the target slot.
		list.blocked = append(list.blocked, interval{})
		for i := len(list.blocked) - 1; i > target; i-- {
			list.blocked[i] = list.blocked[i - 1]
		}
		list.blocked[target] = interval{ low, high }
	}

	// If any of the following intervals are now fully enclosed in the target
	// interval, we remove them.
	for target + 1 < len(list.blocked) && high >= list.blocked[target + 1].high {
		for i := target + 1; i + 1 < len(list.blocked); i++ {
			list.blocked[i] = list.blocked[i + 1]
		}
		list.blocked = list.blocked[:len(list.blocked) - 1]
	}

	// If the new interval still overlaps with the next interval, we merge them.
	if target + 1 < len(list.blocked) && connected(high, list.blocked[target + 1].low) {
		list.blocked[target].high = list.blocked[target + 1].high
		for i := target + 1; i + 1 < len(list.blocked); i++ {
			list.blocked[i] = list.blocked[i + 1]
		}
		list.blocked = list.blocked[:len(list.blocked) - 1]
	}
}

func connected(high, low uint32) bool {
	return high >= low || high + 1 == low
}

func toUint32(v string) uint32 {
	result, err := strconv.ParseUint(v, 10, 32)
	if err != nil { panic(err) }
	return uint32(result)
}
