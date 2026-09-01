package main

import (
	"fmt"
	"slices"
)

var intT []int = []int{1, 2, 3}

// func rr(list []int) {
// 	i := 0
// 	for {
// 		val := list[i]
// 		fmt.Println(val)
// 		i += 1
// 		if i == len(list) {
// 			i = 0
// 		}
// 	}
// }

func main() {
	i := 0
	intT = slices.Concat(intT[:i], intT[i:])
	fmt.Println(intT)
	// rr([]int{1, 2, 3})
}
