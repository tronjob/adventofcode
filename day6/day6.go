package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

type coord struct {
	x int
	y int
}

type mapCoord struct {
	minDist      int //minimum distance to nearest point
	minDistIndex int // index of the nearest point
	nearestPoint coord
	equidistant  bool
	init         bool
}

var infinityMap [][]mapCoord
var coords [50]coord

func atoiHack(v string) int {
	i, _ := strconv.Atoi(v)
	return i
}

func manhattanDistance(x1 int, x2 int, y1 int, y2 int) int {
	return int(math.Abs(float64(x1-x2)) + math.Abs(float64(y1-y2)))
}

func initMap() { //setup start coordinates in map
	for i := range coords {
		x := coords[i].x
		y := coords[i].y

		infinityMap[x][y] = mapCoord{
			minDist:      -1,
			minDistIndex: i,
			nearestPoint: coord{
				x: x,
				y: y,
			},
			equidistant: false,
			init:        false,
		}

	}
}

func calcPointMinDist(workerNum int, start int, end int, coords [50]coord, wg *sync.WaitGroup, ch chan bool) { //gets distance to a set of points for every element in the matrix and updates the matrix
	defer wg.Done()
	fmt.Println("Worker reporting : ", workerNum)
	for i := start; i < end; i++ {
		for j := range infinityMap[i] {
			for k := range coords {
				//critical zone
				//ch <- true
				if !infinityMap[i][j].init && infinityMap[i][j].minDist != -1 { //coordinate not initialized and not start point
					infinityMap[i][j] = mapCoord{
						minDist:      math.MaxInt64,
						minDistIndex: -1,
						nearestPoint: coord{
							x: -1,
							y: -1,
						},
						equidistant: false,
						init:        true,
					}
				}

				if infinityMap[i][j].minDist != -1 { //not a start coordinate
					dist := manhattanDistance(coords[k].x, i, coords[k].y, j)
					if dist == infinityMap[i][j].minDist { // equidistant with previous point
						infinityMap[i][j].equidistant = true
					} else {
						if !infinityMap[i][j].equidistant && dist < infinityMap[i][j].minDist { //not equidistant with any other point + new min distance
							infinityMap[i][j].minDist = dist
							infinityMap[i][j].minDistIndex = k
							infinityMap[i][j].nearestPoint = coords[k]
						}
					}
				} else {
					//<-ch
					continue
				}
				//critical zone end
				//<-ch
			}
		}
	}
}

func printMap() {
	for i := range infinityMap {
		for j := range infinityMap[i] {
			fmt.Print(i, j)
			fmt.Printf("minDist %d, minDistIndex: %d nearestPointX: %d nearestPontY: %d equidistant: %t init: %t\n", infinityMap[i][j].minDist, infinityMap[i][j].minDistIndex, infinityMap[i][j].nearestPoint.x, infinityMap[i][j].nearestPoint.y, infinityMap[i][j].equidistant, infinityMap[i][j].init)
		}
	}
}

func calcMaxArea(coords [50]coord) int {
	maxArea := math.MinInt64
	for k := range coords {
		infiniteArea := false
		currentArea := 0
		for i := range infinityMap {
			for j := range infinityMap[i] {
				if infinityMap[i][j].nearestPoint == coords[k] && !infinityMap[i][j].equidistant { //
					if i == 0 || i == len(infinityMap)-1 || j == 0 || j == len(infinityMap[i])-1 { //edge
						infiniteArea = true
						break
					}
					currentArea++
				}
				if infiniteArea {
					break
				}
			}
			if infiniteArea {
				break
			}
		}
		if !infiniteArea && currentArea > maxArea {
			fmt.Printf("Area of (%d,%d): %d\n", coords[k].x, coords[k].y, currentArea)
			maxArea = currentArea
		}
	}
	return maxArea
}

func main() {
	var wg sync.WaitGroup

	f, _ := os.Open("input.txt")
	scanner := bufio.NewScanner(f)
	var i int
	for i = 0; scanner.Scan(); i++ {
		line := scanner.Text()
		inputCoords := strings.Split(line, ", ")
		coords[i] = coord{x: atoiHack(inputCoords[0]), y: atoiHack(inputCoords[1])}
	}
	infinityMap = make([][]mapCoord, 500)
	for j := range infinityMap {
		infinityMap[j] = make([]mapCoord, 500)
	}
	initMap()
	nworkers := 10
	div := len(infinityMap) / nworkers //50
	ch := make(chan bool, 1)
	fmt.Printf("Firing up %d workers\n\n", nworkers)
	for i := 0; i < nworkers; i++ {
		startIndex := i * div        // 0 50 100 150 ... 450
		endIndex := startIndex + div //50 100 ... 500
		fmt.Printf("Firing up worker %d with params: start %d end %d\n", i, startIndex, endIndex)
		wg.Add(1)
		go calcPointMinDist(i, startIndex, endIndex, coords, &wg, ch)
	}
	wg.Wait()
	fmt.Println(calcMaxArea(coords))
	//printMap()
}

/*

Goroutine idea

8 gr checking each matrix coordinate for minimum distance to coords points
50 coords -> 50/5 for each goroutine = 10 coords to check
post distance to channel?
something gets distances from channels, gets min, updates matrix
rinse and repeat for each matrix coord

*/
