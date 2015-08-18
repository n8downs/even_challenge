package main

import (
	"fmt"
	"math"
	"time"

	t "github.com/ndowns/even_challenge/Types"
)

func main() {
	startDay, _ := time.Parse(t.DateFormat, "2015.08.01")
	endDay, _ := time.Parse(t.DateFormat, "2015.08.31")

	incomes := []t.Income{}
	incomes = append(incomes,
		t.Income{
			Amount:   500.,
			Name:     "Philz",
			Schedule: t.Schedule{Period: t.BiMonthly},
		},
		t.Income{
			Amount:   175.,
			Name:     "Mission Cliffs",
			Schedule: t.Schedule{Period: t.BiWeekly, Weekday: time.Thursday},
		})
	fmt.Println("Incomes:  ", incomes)

	expenses := []t.Expense{}
	expenses = append(expenses,
		t.Expense{
			Amount:   42.34,
			Name:     "Utilities",
			Schedule: t.Schedule{Period: t.Monthly, Date: 25},
		},
		t.Expense{
			Amount:   400.,
			Name:     "Rent",
			Schedule: t.Schedule{Period: t.Monthly, Date: 28},
		},
		t.Expense{
			Amount:   40.,
			Name:     "Crossfit",
			Schedule: t.Schedule{Period: t.Weekly, Weekday: time.Tuesday},
		},
	)
	fmt.Println("Expenses: ", expenses)

	_ = simulate(startDay, endDay, incomes, expenses)
}

func simulate(startDay time.Time, endDay time.Time, incomes []t.Income, expenses []t.Expense) map[t.Account]float64 {
	accounts := map[t.Account]float64{t.External: 0., t.Checking: 0., t.Savings: 0.}

	totalIncome := 0.
	totalExpenses := 0.

	incomeTotals := map[time.Time]float64{}
	savingsPlan := map[time.Time]float64{}
	expenseTotals := map[time.Time]float64{}
	ledger := map[time.Time][]t.Transaction{}

	for _, income := range incomes {
		occurrances := income.Schedule.FindOccurrances(startDay, endDay)
		for _, date := range occurrances {
			ledger[date] = append(ledger[date], t.Transaction{
				Date:  date,
				Delta: income.Amount,
				Memo:  fmt.Sprintf("Income: %s", income.Name),
				From:  t.External,
				To:    t.Checking,
			})
			totalIncome += income.Amount
			incomeTotals[date] += income.Amount
			savingsPlan[date] = 0.
		}
	}

	for _, expense := range expenses {
		occurrances := expense.Schedule.FindOccurrances(startDay, endDay)
		for _, date := range occurrances {
			totalExpenses += expense.Amount
			expenseTotals[date] += expense.Amount
		}
	}

	idealDiscretionary := (totalIncome - totalExpenses) / (endDay.Sub(startDay).Hours() / 24)
	fmt.Printf("Total: $%.2f in, $%.2f out, $%.2f ideally per day\n\n", totalIncome, totalExpenses, idealDiscretionary)
	if totalExpenses > totalIncome {
		fmt.Printf("Insolvent :(")
		return accounts
	}

	discretionaryAmount := 0.
	discretionaryDays := 0.
	currentDate := startDay
	runningPlan := 0.
	for {
		if currentDate.After(endDay) {
			break
		}

		if incomeTotals[currentDate] != 0. {
			nextIncomeDate := currentDate
			upcomingExpenses := 0.
			for {
				if nextIncomeDate.After(endDay) {
					break
				}

				upcomingExpenses += expenseTotals[nextIncomeDate]
				nextIncomeDate = nextIncomeDate.AddDate(0, 0, 1)
				if incomeTotals[nextIncomeDate] != 0. {
					break
				}
			}
			daysUntilNextIncome := nextIncomeDate.Sub(currentDate).Hours() / 24
			mustTransfer := upcomingExpenses - runningPlan
			idealTransfer := incomeTotals[currentDate] - idealDiscretionary*daysUntilNextIncome
			transfer := math.Max(mustTransfer, idealTransfer)
			runningPlan = runningPlan + transfer - upcomingExpenses
			savingsPlan[currentDate] = -1 * transfer
			actualDiscretionary := (incomeTotals[currentDate] - transfer) / daysUntilNextIncome

			discretionaryAmount += incomeTotals[currentDate] - transfer
			discretionaryDays += daysUntilNextIncome

			simulatedSpendingDate := currentDate
			for {
				if simulatedSpendingDate.Equal(nextIncomeDate) {
					break
				}
				ledger[simulatedSpendingDate] = append(ledger[simulatedSpendingDate], t.Transaction{
					From:  t.Checking,
					To:    t.External,
					Memo:  "Simulated Spending",
					Date:  simulatedSpendingDate,
					Delta: -1 * actualDiscretionary,
				})
				simulatedSpendingDate = simulatedSpendingDate.AddDate(0, 0, 1)
			}
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	fmt.Printf("Planned average discretionary spending levels: $%.2f (%.2f%%)\n\n", discretionaryAmount/discretionaryDays, discretionaryAmount/discretionaryDays/idealDiscretionary)

	for date, amount := range savingsPlan {
		if amount < 0 {
			ledger[date] = append(ledger[date], t.Transaction{
				Date:  date,
				Delta: amount,
				Memo:  "Transfer to Savings",
				From:  t.Checking,
				To:    t.Savings,
			})
		} else {
			ledger[date] = append(ledger[date], t.Transaction{
				Date:  date,
				Delta: amount,
				Memo:  "Transfer from Savings",
				From:  t.Savings,
				To:    t.Checking,
			})
		}
	}

	for _, expense := range expenses {
		occurrances := expense.Schedule.FindOccurrances(startDay, endDay)
		for _, date := range occurrances {
			ledger[date] = append(ledger[date], t.Transaction{
				Date:  date,
				Delta: -1 * expense.Amount,
				Memo:  fmt.Sprintf("Expense: %s", expense.Name),
				From:  t.Savings,
				To:    t.External,
			})
		}
	}

	currentDate = startDay
	inARow := 0
	for {
		if currentDate.After(endDay) {
			break
		}

		transactions := ledger[currentDate]
		if len(transactions) == 0 {
			inARow++
			if inARow <= 1 {
				fmt.Println("    ...    |")
			}
		} else {
			inARow = 0
		}

		for _, transaction := range transactions {
			accounts[transaction.From] -= math.Abs(transaction.Delta)
			accounts[transaction.To] += math.Abs(transaction.Delta)
			fmt.Printf("%s | %7.2f | %7.2f\n", transaction.ToString(), accounts[t.Checking], accounts[t.Savings])
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return accounts
}

/*





































*/
