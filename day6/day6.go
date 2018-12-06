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
		}

	}

}

func calcPointMinDist(workerNum int, coords []coord, wg *sync.WaitGroup, ch chan bool) { //gets distance to a set of points for every element in the matrix and updates the matrix
	defer wg.Done()
	fmt.Println("Worker reporting : ", workerNum)
	for i := range infinityMap {
		for j := range infinityMap[i] {
			for k := range coords {
				//critical zone
				ch <- true
				if infinityMap[i][j].minDist == 0 && infinityMap[i][j].minDist != -1 { //coordinate not initialized and not start point
					infinityMap[i][j] = mapCoord{
						minDist:      math.MaxInt64,
						minDistIndex: -1,
						nearestPoint: coord{
							x: -1,
							y: -1,
						},
						equidistant: false,
					}
				}

				if infinityMap[i][j].minDist != -1 { //nao e uma coordenada input
					if manhattanDistance(coords[k].x, coords[k].y, i, j) == infinityMap[i][j].minDist { //ponto equidistante
						infinityMap[i][j].equidistant = true
					} else {
						dist := manhattanDistance(coords[k].x, coords[k].y, i, j)
						if !infinityMap[i][j].equidistant && dist < infinityMap[i][j].minDist { //dist e menor que a minima ja encontrada para algum ponto
							infinityMap[i][j].minDist = dist
							infinityMap[i][j].minDistIndex = k
							infinityMap[i][j].nearestPoint = coords[k]
						}
					}
				} else {
					<-ch
					continue
				}
				//critical zone end
				<-ch
			}
		}
	}
}

func printMap() {
	for i := range infinityMap {
		for j := range infinityMap[i] {
			fmt.Print(i, j)
			fmt.Printf("minDist %d, minDistIndex: %d nearestPointX: %d nearestPontY: %d equidistant: %t\n", infinityMap[i][j].minDist, infinityMap[i][j].minDistIndex, infinityMap[i][j].nearestPoint.x, infinityMap[i][j].nearestPoint.y, infinityMap[i][j].equidistant)
		}
	}
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
	div := len(coords) / nworkers
	ch := make(chan bool, 1)
	for i := 0; i < nworkers; i++ {
		startIndex := i * div // 0 5 10 15 ... 45
		endIndex := startIndex + div
		wg.Add(1)
		go calcPointMinDist(i, coords[startIndex:endIndex], &wg, ch)
	}
	wg.Wait()
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
