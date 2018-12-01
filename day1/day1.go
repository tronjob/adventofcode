package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func dealWithIt(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func contains(history []int, freq int) bool {
	for _, val := range history {
		if val == freq {
			return true
		}
	}
	return false
}

func main() {
	freq := 0
	repeatFreq := -1   // first repeat freq
	oneCycleFreq := -1 // freq after one cycle
	history := make([]int, 100)
	f, err := os.Open("input.txt")
	dealWithIt(err)
	for repeatFreq == -1 {
		f.Seek(0, 0)
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			delta := scanner.Text()
			deltaInt, err := strconv.Atoi(delta)
			dealWithIt(err)
			freq += deltaInt
			if repeatFreq == -1 && contains(history, freq) {
				repeatFreq = freq
				break
			}

			history = append(history, freq)
		}
		if oneCycleFreq == -1 {
			oneCycleFreq = freq
		}

	}
	fmt.Println(oneCycleFreq)
	fmt.Println(repeatFreq)
}
