package main

import (
	"fmt"
	"time"

	t "github.com/ndowns/even_challenge/Types"
)

const (
	dateFormat = "2015.01.02"
)

func main() {
	fmt.Println("Hello World")

	incomes := []t.Income{}
	fmt.Println("test")
	fmt.Println(incomes)
	fmt.Println(time.Second)

	transaction := t.Transaction{}
	//transaction := t.Transaction{date: time.Parse(dateFormat, "2015.08.08"), delta: -100.32, memo: "SAVE for rent"}
	fmt.Println(transaction)
}
