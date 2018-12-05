package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
)

var polymerChain []rune

type polymer struct {
	unit          rune
	originalIndex int
}

func hasReaction(a rune, b rune) bool {
	fmt.Printf("comparing char %c with %c \n", unicode.SimpleFold(a), b)
	if unicode.SimpleFold(a) == b {
		return true
	}
	return false
}

func doChainReaction(start int, end int) (int, int) {
	didChainReaction := false

	if start < 0 {
		fmt.Println("Chain reacting lt 0")
		return 0, end
	}
	if end >= len(polymerChain) {
		fmt.Println("Chain reacting gt end")
		return start + 1, len(polymerChain)
	}
	fmt.Println("Chain reacting ", start, polymerChain[start], end, polymerChain[end])
	for hasReaction(polymerChain[start], polymerChain[end]) {
		didChainReaction = true
		start--
		end++
		if start < 0 {
			return 0, end
		}
		if end >= len(polymerChain) {
			return start, len(polymerChain) - 1
		}
	}
	if !didChainReaction {
		return -1, -1
	}
	start++
	end--
	fmt.Printf("reaction with start %d end %d length %d found\n", start, end, end-start)
	return start, end
}

func findNewEndIndex(reactedChain []polymer, index int) int {
	fmt.Println("Looking for ", index)
	if len(reactedChain) == 0 {
		return 0
	}
	for i := range reactedChain {
		if reactedChain[i].originalIndex == index {
			return i
		}
	}
	log.Output(1, "Original index not found")
	return len(reactedChain)
}

func printChain(p []polymer) {
	if len(p) == 0 {
		fmt.Printf("Empty chain\n")
	}
	for i := range p {
		fmt.Printf("Unit %c, orig index %d\n", p[i].unit, p[i].originalIndex)
	}
}
func makeReactions() []polymer {
	reactedChain := make([]polymer, 0)
	curr := 0
	for range polymerChain {
		fmt.Println("cycle ", curr)
		if curr == 0 {
			reactedChain = append(reactedChain, polymer{unit: polymerChain[curr], originalIndex: curr})
			curr++
			continue
		}
		if curr > len(polymerChain)-1 {
			break
		}
		fmt.Println("reach1", curr)
		if hasReaction(polymerChain[curr], polymerChain[curr-1]) {
			fmt.Println("reach2")
			start, end := -1, -1
			if curr >= 2 {
				start, end = doChainReaction(curr-2, curr+1)
			}
			fmt.Println("reach3")

			if start == -1 && end == -1 { //no chain reaction
				fmt.Println("no chain reaction")
				reactedChain = reactedChain[0:findNewEndIndex(reactedChain, curr)]
			} else { //chain reaction
				fmt.Println("chain reaction")
				reactedChain = reactedChain[0:findNewEndIndex(reactedChain, start)]
				curr = end
			}

		} else {
			fmt.Println("no reaction")
			reactedChain = append(reactedChain, polymer{unit: polymerChain[curr], originalIndex: curr})
		}
		curr++
		//printChain(reactedChain)
	}
	return reactedChain
}

func main() {
	f, _ := os.Open("input.txt")
	scanner := bufio.NewScanner(f)
	polymerChain = make([]rune, 0)
	for scanner.Scan() {
		polymerChain = []rune(scanner.Text())
	}
	reactedChain := makeReactions()

	fmt.Println("part1: ", len(reactedChain))

	fmt.Println(len(polymerChain))

}

//compare i and i-1
//if they react, erase i-1 from new chain (if added) and skip i
//check for chain reactions from i-2 back and i+1 forward
//delete the subchain
//if they don't react, add i to the chain
