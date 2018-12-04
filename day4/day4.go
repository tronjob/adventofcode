package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var dateLayout = "2006-01-02 15:04" //date format

type guardActivity struct {
	timestamp time.Time
	eventType string //SHIFTSTART, SLEEP, WAKE
	guardID   int
}

var guardLog []guardActivity

func dealWithIt(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func contains(s []int, el int) bool {
	for i := range s {
		if s[i] == el {
			return true
		}
	}
	return false
}

func getGuardList() []int {
	guardList := make([]int, 0)
	for i := range guardLog {
		if !contains(guardList, guardLog[i].guardID) {
			guardList = append(guardList, guardLog[i].guardID)
		}
	}
	return guardList
}

func addActivity(timeString string, guardAction string) {
	r, _ := regexp.Compile("(?m)^Guard #([0-9]+) begins shift$")
	t, _ := time.Parse(dateLayout, timeString)

	if r.MatchString(guardAction) {
		matches := r.FindStringSubmatch(guardAction)
		guardID, _ := strconv.Atoi(matches[1])
		guardLog = append(guardLog, guardActivity{timestamp: t, eventType: "SHIFTSTART", guardID: guardID})
	} else {
		switch guardAction {
		case "wakes up":
			guardLog = append(guardLog, guardActivity{timestamp: t, eventType: "WAKE", guardID: -1})

		case "falls asleep":
			guardLog = append(guardLog, guardActivity{timestamp: t, eventType: "SLEEP", guardID: -1})
		}
	}
}

func completeLog() { //Add guardID to every guardLog event
	lastID := -2
	for i := range guardLog {
		if guardLog[i].guardID != -1 { //guard started shift
			lastID = guardLog[i].guardID
		} else { //sleep/wake
			guardLog[i].guardID = lastID
		}
	}
}

func intMinutes(dif string) int {
	r := regexp.MustCompile("^([0-9]+).*$")
	mins, _ := strconv.Atoi(r.FindStringSubmatch(dif)[1])
	return mins
}

func getMax(s []int) (int, int) { //returns the minute most slept and how many times the guard was asleep in that minute
	max := s[0]
	for i := range s {
		if s[i] > s[max] {
			max = i
		}
	}
	return max, s[max]
}

func getSleepyMinute(guardID int) (int, int) { //finds the minute most spent sleeping by a certain guard and the number of days he spent sleeping on that minute
	var t1 int
	var t2 int
	var minutes = make([]int, 60)
	for i := range guardLog {
		if guardLog[i].guardID != guardID {
			continue
		}
		if guardLog[i].eventType == "WAKE" {
			t2 = guardLog[i].timestamp.Minute()
			t1 = guardLog[i-1].timestamp.Minute()
			for j := t1; j < t2; j++ {
				minutes[j]++
			}
		}
	}
	maxMinute, nDays := getMax(minutes)
	return maxMinute, nDays
}

func findSleepyHead() int { //finds guardID of the guard that sleeps the most
	cumulativeSleep := make(map[int]int)
	for i := range guardLog {
		if guardLog[i].eventType == "WAKE" {
			cumulativeSleep[guardLog[i].guardID] += intMinutes(guardLog[i].timestamp.Sub(guardLog[i-1].timestamp).String())
		}
	}
	maxSleep := 0
	sleepyID := -1

	for g, s := range cumulativeSleep {
		if s > maxSleep {
			sleepyID = g
			maxSleep = s
		}
	}
	return sleepyID
}

func part1() {
	completeLog()
	gID := findSleepyHead()
	sleepyMinute, _ := getSleepyMinute(gID)
	fmt.Println("part1: ", gID*sleepyMinute)
}

func part2() {
	maxMinute := 0
	maxSleeper := -1
	maxDaysAsleepOnSameMin := 0
	guardList := getGuardList()
	for i := range guardList {
		maxMin, nDays := getSleepyMinute(guardList[i])
		if nDays > maxDaysAsleepOnSameMin {
			maxDaysAsleepOnSameMin = nDays
			maxSleeper = guardList[i]
			maxMinute = maxMin
		}
	}
	fmt.Println("part2: ", maxMinute*maxSleeper)
}

func main() {
	f, err := os.Open("input.txt")
	dealWithIt(err)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	guardLog = make([]guardActivity, 0)
	for scanner.Scan() {
		line := scanner.Text()
		r := regexp.MustCompile("(?m)^\\[([-0-9]+ [0-9]{2}:[0-9]{2})\\] ([a-zA-Z0-9# ]+)$")
		matches := r.FindStringSubmatch(line)
		addActivity(matches[1], matches[2])
	}
	sort.Slice(guardLog, func(a, b int) bool { return guardLog[a].timestamp.Before(guardLog[b].timestamp) })
	part1()
	part2()

}
