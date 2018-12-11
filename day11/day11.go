package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

var theGrid [300][300]int

/*

    --------------> X,i
   |
 Y |
   |
 j |
v

X is horizontal
Y is vertical

matrix[line][col] = matrix[y][x]
*/

func getPowerLevel(x int, y int, sn int) int {
	rackID := x + 10
	powerLevel := rackID * y
	powerLevel += sn
	powerLevel *= rackID
	if powerLevel >= 100 {
		powerLevel = (powerLevel / 100) % 10
	} else {
		powerLevel = 0
	}
	powerLevel -= 5
	return powerLevel
}

//calculates powerLevel in a 3x3 grid
func getGridPowerLevel(x int, y int, sn int, delta int) int {
	sum := 0
	for i := x; i < x+delta; i++ {
		if i > 300 {
			break
		}
		for j := y; j < y+delta; j++ {
			if j > 300 {
				break
			}
			sum += getPowerLevel(i, j, sn)
		}
	}
	/*
		1 2 3
		4 5 6
		7 8 9
	*/
	return sum
}
func findSum(sn int) (int, int, int) {
	maxPowaaa := math.MinInt64
	maxX, maxY := 0, 0
	for i := 0; i < 300; i++ {
		for j := 0; j < 300; j++ {
			localPower := getGridPowerLevel(i+1, j+1, sn, 3)
			if localPower > maxPowaaa {
				maxX = i + 1
				maxY = j + 1
				maxPowaaa = localPower
			}
		}
	}
	return maxX, maxY, maxPowaaa
}

func findSumAny(sn int) (int, int, int) {
	maxPowaaa := math.MinInt64
	maxX, maxY := 0, 0
	maxDelta := 0
	for i := 0; i < 300; i++ {
		for j := 0; j < 300; j++ {
			for delta := 1; delta < 300; delta++ {
				localPower := getGridPowerLevel(i+1, j+1, sn, delta)
				if localPower > maxPowaaa {
					maxX = i + 1
					maxY = j + 1
					maxPowaaa = localPower
					maxDelta = delta
				}
			}
		}
	}
	return maxX, maxY, maxDelta
}

func part1(sn int) {
	maxX, maxY, maxPower := findSum(sn)
	fmt.Println(maxX, maxY, maxPower)
}

func part2(sn int) {
	maxX, maxY, maxDelta := findSumAny(sn)
	fmt.Println(maxX, maxY, maxDelta)
}
func main() {
	f, _ := os.Open("input.txt")
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	sn, _ := strconv.Atoi(scanner.Text())
	part1(sn)
	part2(sn)

}
