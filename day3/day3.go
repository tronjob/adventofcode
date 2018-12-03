package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var fabric [][]fabricMaterial

var uniqueCandidates []int

type fabricMaterial struct {
	elfID   int
	nClaims int
}

func initFabric() [][]fabricMaterial {
	fabric := make([][]fabricMaterial, 1000)
	for i := range fabric {
		fabric[i] = make([]fabricMaterial, 1000)
	}
	return fabric
}

func dealWithIt(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getIndex(a []int, val int) int {
	for i, v := range a {
		if v == val {
			return i
		}
	}
	return -1
}
func removeCandidates(nElves ...int) []int {
	for _, elf := range nElves {
		targetIndex := getIndex(uniqueCandidates, elf)
		if targetIndex < 0 {
			continue
		}
		uniqueCandidates[targetIndex] = uniqueCandidates[len(uniqueCandidates)-1]
		uniqueCandidates[len(uniqueCandidates)-1] = 0
		uniqueCandidates = uniqueCandidates[:len(uniqueCandidates)-1]
	}
	return uniqueCandidates
}

func addToFabric(elfID int, indexX int, indexY int, sizeX int, sizeY int) {
	for i := indexX; i < indexX+sizeX; i++ {
		for j := indexY; j < indexY+sizeY; j++ {
			if fabric[i][j].nClaims > 0 {
				uniqueCandidates = removeCandidates(fabric[i][j].elfID, elfID) //if square already claimed, disqualify the current and the previous claimants from uniqueness
			}
			fabric[i][j].elfID = elfID
			fabric[i][j].nClaims++
		}
	}
}

func checkDupes() int {
	dupes := 0
	for i := range fabric {
		for j := range fabric[i] {
			if fabric[i][j].nClaims > 1 {
				dupes++
			}
		}
	}
	return dupes
}

func atoiHack(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func getInput() {
	f, err := os.Open("input.txt")
	dealWithIt(err)
	in := bufio.NewScanner(f)
	for in.Scan() {
		line := in.Text()
		r := regexp.MustCompile("(?m)^#([0-9]+) @ ([0-9]+),([0-9]+): ([0-9]+)x([0-9]+)$")
		matches := r.FindStringSubmatch(line)
		nElf, indexX, indexY, sizeX, sizeY := atoiHack(matches[1]), atoiHack(matches[2]), atoiHack(matches[3]), atoiHack(matches[4]), atoiHack(matches[5])
		uniqueCandidates = append(uniqueCandidates, nElf)
		addToFabric(nElf, indexX, indexY, sizeX, sizeY)
	}
}
func main() {
	fabric = initFabric()
	getInput()
	fmt.Println(checkDupes())
	fmt.Println(uniqueCandidates[0])
}
