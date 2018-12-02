package main

import "fmt"

var strings2 = make([]string, 100)

func cenas() {
	strings2 = append(strings2, "test")
	strings2 = append(strings2, "test2")
	strings2 = append(strings2, "test3")

	for _, s := range strings2 {
		fmt.Printf("s is %d", s)
	}
}
