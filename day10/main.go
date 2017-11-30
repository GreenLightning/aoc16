package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil { panic(err) }

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	bots := newBotList()

	valueRule := regexp.MustCompile(`value (\d+) goes to bot (\d+)`)
	botRule   := regexp.MustCompile(`bot (\d+) gives low to (bot|output) (\d+) and high to (bot|output) (\d+)`)

	for scanner.Scan() {
		instruction := scanner.Text()
		if data := valueRule.FindStringSubmatch(instruction); data != nil {
			value := toInt(data[1])
			id := toInt(data[2])
			bots.getOrCreateBot(id).appendValue(value)
		} else if data := botRule.FindStringSubmatch(instruction); data != nil {
			id := toInt(data[1])
			bot := bots.getOrCreateBot(id)
			bot.lowOut = (data[2] == "output")
			bot.lowTarget = toInt(data[3])
			bot.highOut = (data[4] == "output")
			bot.highTarget = toInt(data[5])
		} else {
			panic(instruction)
		}
	}

	maxOutput := -1
	for it := bots.iterator(); it.next(); {
		bot := it.getBot()
		if bot.lowOut && bot.lowTarget > maxOutput { maxOutput = bot.lowTarget }
		if bot.highOut && bot.highTarget > maxOutput { maxOutput = bot.highTarget }
	}

	outputs := make([][]int, maxOutput + 1)

	changed := true
	for changed {
		changed = false
		for it := bots.iterator(); it.next(); {
			bot := it.getBot()
			if len(bot.values) >= 2 {
				low, high := bot.values[0], bot.values[1]
				bot.values = bot.values[2:]

				if high < low { low, high = high, low }

				if low == 17 && high == 61 {
					fmt.Println("Bot that compares 61 and 17:", it.getId())
				}

				if bot.lowOut {
					outputs[bot.lowTarget] = append(outputs[bot.lowTarget], low)
				} else {
					bots.getBot(bot.lowTarget).appendValue(low)
				}
				if bot.highOut {
					outputs[bot.highTarget] = append(outputs[bot.highTarget], high)
				} else {
					bots.getBot(bot.highTarget).appendValue(high)
				}
				changed = true
			}
		}
	}

	product := 1
	for i := 0; i <= 2 && i < len(outputs); i++ {
		if len(outputs[i]) > 0 {
			product *= outputs[i][0]
			fmt.Print("Output ", i, ":")
			for _, value := range outputs[i] {
				fmt.Print(" ", value)
			}
			fmt.Println()
		}
	}
	fmt.Println("Product:", product)
}

type bot struct {
	values []int
	lowOut bool
	lowTarget int
	highOut bool
	highTarget int
}

type botList struct {
	slice []bot
}

type botIterator struct {
	parent *botList
	bot *bot
	id int
}

func (bot *bot) appendValue(value int) {
	bot.values = append(bot.values, value)
}

func newBotList() *botList {
	return &botList{ make([]bot, 0, 256) }
}

func (list *botList) getOrCreateBot(id int) *bot {
	if id >= len(list.slice) {
		list.slice = list.slice[0:id+1]
	}
	return list.getBot(id)
}

func (list *botList) getBot(id int) *bot {
	return &list.slice[id]
}

func (list *botList) iterator() botIterator {
	return botIterator{ list, nil, -1 }
}

func (it *botIterator) next() bool {
	if it.id + 1 >= len(it.parent.slice) {
		return false
	}
	it.id++
	it.bot = &it.parent.slice[it.id]
	return true
}

func (it *botIterator) getId() int {
	return it.id
}

func (it *botIterator) getBot() *bot {
	return it.bot
}

func toInt(v string) int {
	i, e := strconv.Atoi(v)
	if e != nil { panic(e) }
	return i
}
