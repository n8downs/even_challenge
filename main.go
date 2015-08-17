package main

import (
	"fmt"
	"time"

	t "github.com/ndowns/even_challenge/Types"
)

const (
	dateFormat = "2006.01.02"
)

func main() {
	fromDay, _ := time.Parse(dateFormat, "2015.08.01")
	toDay, _ := time.Parse(dateFormat, "2015.08.31")

	incomes := []t.Income{}
	incomes = append(incomes, t.Income{
		Amount:   500.,
		Name:     "Buckstars",
		Schedule: t.Schedule{Period: t.BiMonthly},
	})
	fmt.Println(incomes)

	expenses := []t.Expense{}
	expenses = append(expenses,
		t.Expense{
			Amount:   -42.34,
			Name:     "Utilities",
			Schedule: t.Schedule{Period: t.Monthly, Date: 25},
		},
		t.Expense{
			Amount:   -400.,
			Name:     "Rent",
			Schedule: t.Schedule{Period: t.Monthly, Date: 1},
		},
	)
	fmt.Println(expenses)

	totalIncome := 0.
	totalExpenses := 0.

	ledger := map[time.Time][]t.Transaction{}

	for _, income := range incomes {
		occurrances := income.Schedule.FindOccurrances(fromDay, toDay)
		for _, date := range occurrances {
			ledger[date] = append(ledger[date], t.Transaction{
				Date:  date,
				Delta: income.Amount,
				Memo:  fmt.Sprintf("Income: %s", income.Name),
			})
			totalIncome += income.Amount
		}
	}

	for _, expense := range expenses {
		occurrances := expense.Schedule.FindOccurrances(fromDay, toDay)
		for _, date := range occurrances {
			ledger[date] = append(ledger[date], t.Transaction{
				Date:  date,
				Delta: expense.Amount,
				Memo:  fmt.Sprintf("Expense: %s", expense.Name),
			})
			totalExpenses += expense.Amount
		}
	}

	currentDate := fromDay
	for {
		if currentDate.After(toDay) {
			break
		}

		fmt.Println(currentDate.Format(dateFormat))
		for _, transaction := range ledger[currentDate] {
			fmt.Println(transaction.ToString())
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	idealDiscretionary := (totalIncome + totalExpenses) / (toDay.Sub(fromDay).Hours() / 24)
	fmt.Printf("\n$%.2f in, $%.2f out, $%.2f ideally per day\n", totalIncome, totalExpenses, idealDiscretionary)

	//expenseSavings := 0.
	/*
		transactions := []t.Transaction{}
		date, _ := time.Parse(dateFormat, "2015.08.14")
		transactions = append(transactions, t.Transaction{Date: date, Delta: -100.32, Memo: "SAVER for rent"})
		date, _ = time.Parse(dateFormat, "2015.08.31")
		transactions = append(transactions, t.Transaction{Date: date, Delta: 12.34, Memo: "TRANS for rent"})

		for _, transaction := range transactions {
			expenseSavings = expenseSavings - transaction.Delta
			fmt.Printf("%s | %.2f\n", transaction.ToString(), expenseSavings)
		}
	*/
}
