package main

import "fmt"

func jiafa() (int, int) {
	return 3, 8
}

func add(x, y int) int {
	return x + y
}

type addFunc func(x, y int) int

func Opt(x, y int, addFunc2 addFunc) int {
	return addFunc2(x, y)
}

func main() {
	//int1, int2 := jiafa()
	//fmt.Print(int1 + int2)

	x := 3
	y := 9
	//result := add(x, y)
	//fmt.Println(result)

	res := Opt(x, y, add)
	fmt.Println(res)
}
