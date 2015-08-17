package main

import (
	"fmt"
	"time"

	t "github.com/ndowns/even_challenge/Types"
)

const (
	dateFormat = "2006.01.02"

	startingBalance = 1000.
)

func main() {
	incomes := []t.Income{}
	fmt.Println(incomes)
	balance := startingBalance

	date, _ := time.Parse(dateFormat, "2015.08.14")
	transaction := t.Transaction{Date: date, Delta: -100.32, Memo: "SAVER for rent"}
	balance = balance + transaction.Delta
	fmt.Printf("%s | %.2f\n", transaction.ToString(), balance)

	date, _ = time.Parse(dateFormat, "2015.08.31")
	date = date.Add(26 * time.Hour)
	transaction = t.Transaction{Date: date, Delta: 12.34, Memo: "TRANS for rent"}
	balance = balance + transaction.Delta
	fmt.Printf("%s | %.2f\n", transaction.ToString(), balance)
}
