package main

import "fmt"

type schedule struct {
	period string
	day    string
}

type income struct {
	amount   int
	schedule schedule
}

func main() {
	fmt.Println("Hello World")

	incomes := []income{}
	fmt.Println("test")
	fmt.Println(incomes)
}
