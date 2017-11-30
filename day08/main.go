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

	rect := regexp.MustCompile(`rect (\d+)x(\d+)`)
	rrow := regexp.MustCompile(`rotate row y=(\d+) by (\d+)`)
	rcol := regexp.MustCompile(`rotate column x=(\d+) by (\d+)`)

	width := 50
	height := 6

	data := make([]bool, width * height)

	for scanner.Scan() {
		instruction := scanner.Text()
		if result := rect.FindStringSubmatch(instruction); result != nil {
			w := toInt(result[1])
			h := toInt(result[2])
			for x := 0; x < w; x++ {
				for y := 0; y < h; y++ {
					data[y * width + x] = true
				}
			}
		} else if result := rrow.FindStringSubmatch(instruction); result != nil {
			y := toInt(result[1])
			a := toInt(result[2])
			for i := 0; i < a; i++ {
				tmp := data[y * width + width-1]
				for x := width-1; x > 0; x-- {
					data[y * width + x] = data[y * width + x-1]
				}
				data[y * width + 0] = tmp
			}
		} else if result := rcol.FindStringSubmatch(instruction); result != nil {
			x := toInt(result[1])
			a := toInt(result[2])
			for i := 0; i < a; i++ {
				tmp := data[(height-1) * width + x]
				for y := height-1; y > 0; y-- {
					data[y * width + x] = data[(y-1) * width + x]
				}
				data[0 * width + x] = tmp
			}
		} else {
			panic(instruction)
		}
	}

	count := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if data[y * width + x] {
				count++
			}
		}
	}

	fmt.Println(count)
	print(data, width, height)
}

func toInt(v string) int {
	i, e := strconv.Atoi(v)
	if e != nil { panic(e) }
	return i
}

func print(data []bool, width, height int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if data[y * width + x] {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
