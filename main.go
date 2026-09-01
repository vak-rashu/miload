package main

import "fmt"

func rr(list []int) {
	i := 0
	for {
		val := list[i]
		fmt.Println(val)
		i += 1
		if i == len(list) {
			i = 0
		}
	}
}

func main() {
	rr([]int{1, 2, 3})
}
