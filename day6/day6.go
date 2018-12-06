package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type coord struct {
	x    int
	y    int
	area int
}

var infinityMap [][]rune
var coords [50]coord

func atoiHack(v string) int {
	i, _ := strconv.Atoi(v)
	return i
}

func manhattanDistance(x1 int, x2 int, y1 int, y2 int) int {
	return int(math.Abs(float64(x1-x2)) + math.Abs(float64(y1-y2)))
}

/*func initMap() { //setup start coordinates in map
	for i := range coords {
		fmt.Println(i)
		x := coords[i].x
		y := coords[i].y
		infinityMap[x][y] = rune(i)
	}

}*/
func main() {

	f, _ := os.Open("input.txt")
	scanner := bufio.NewScanner(f)
	var i int
	for i = 0; scanner.Scan(); i++ {
		line := scanner.Text()
		inputCoords := strings.Split(line, ", ")
		coords[i] = coord{x: atoiHack(inputCoords[0]), y: atoiHack(inputCoords[1]), area: 1}
	}
	infinityMap = make([][]rune, 500)
	for j := range infinityMap {
		infinityMap[j] = make([]rune, 500)
	}
	fmt.Println(infinityMap)
}

/*

Goroutine idea

8 gr checking each matrix coordinate for minimum distance to coords points
50 coords -> 50/5 for each goroutine = 10 coords to check
post distance to channel?
something gets distances from channels, gets min, updates matrix
rinse and repeat for each matrix coord

*/
