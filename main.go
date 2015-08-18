package main

import (
	"fmt"
	"math"
	"time"

	"github.com/ndowns/even_challenge/Types"
)

const simulatedSpendingMemo = "  -Simulated Spending-"

func main() {
	startDay, _ := time.Parse(Types.DateFormat, "2015.08.01")
	endDay, _ := time.Parse(Types.DateFormat, "2015.08.31")

	incomes := []Types.Income{}
	incomes = append(incomes,
		Types.Income{
			Amount:   500.,
			Name:     "Philz",
			Schedule: Types.Schedule{Period: Types.BiMonthly},
		},
		Types.Income{
			Amount:   175.,
			Name:     "Mission Cliffs",
			Schedule: Types.Schedule{Period: Types.BiWeekly, Weekday: time.Thursday},
		})
	fmt.Println("Incomes:  ", incomes)

	expenses := []Types.Expense{}
	expenses = append(expenses,
		Types.Expense{
			Amount:   42.34,
			Name:     "Utilities",
			Schedule: Types.Schedule{Period: Types.Monthly, Date: 25},
		},
		Types.Expense{
			Amount:   400.,
			Name:     "Rent",
			Schedule: Types.Schedule{Period: Types.Monthly, Date: 28},
		},
		Types.Expense{
			Amount:   40.,
			Name:     "Crossfit",
			Schedule: Types.Schedule{Period: Types.Weekly, Weekday: time.Tuesday},
		},
	)
	fmt.Println("Expenses: ", expenses)

	plan, _ := Plan(startDay, endDay, incomes, expenses)

	_, _ = Simulate(startDay, endDay, plan, true)
}

// Plan ...
func Plan(
	startDay time.Time,
	endDay time.Time,
	incomes []Types.Income,
	expenses []Types.Expense,
) (map[time.Time][]Types.Transaction, float64) {
	ledger := map[time.Time][]Types.Transaction{}
	totalIncome := 0.
	totalExpenses := 0.

	incomeTotals := map[time.Time]float64{}
	savingsPlan := map[time.Time]float64{}
	expenseTotals := map[time.Time]float64{}

	for _, income := range incomes {
		occurrances := income.Schedule.FindOccurrances(startDay, endDay)
		for _, date := range occurrances {
			ledger[date] = append(ledger[date], Types.Transaction{
				Date:  date,
				Delta: income.Amount,
				Memo:  fmt.Sprintf("Income: %s", income.Name),
				From:  Types.External,
				To:    Types.Checking,
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

	if totalExpenses > totalIncome {
		return map[time.Time][]Types.Transaction{}, 0.
	}

	idealDiscretionary := (totalIncome - totalExpenses) / (endDay.Sub(startDay).Hours() / 24)
	//fmt.Printf("Total: $%.2f in, $%.2f out, $%.2f ideally per day\n\n", totalIncome, totalExpenses, idealDiscretionary)

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
				ledger[simulatedSpendingDate] = append(ledger[simulatedSpendingDate], Types.Transaction{
					From:  Types.Checking,
					To:    Types.External,
					Memo:  simulatedSpendingMemo,
					Date:  simulatedSpendingDate,
					Delta: -1 * actualDiscretionary,
				})
				simulatedSpendingDate = simulatedSpendingDate.AddDate(0, 0, 1)
			}
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	//fmt.Println(idealDiscretionary)
	//fmt.Printf("Planned average discretionary spending levels: $%.2f (%.2f%%)\n\n", discretionaryAmount/discretionaryDays, discretionaryAmount/discretionaryDays/idealDiscretionary)

	for date, amount := range savingsPlan {
		if amount < 0 {
			ledger[date] = append(ledger[date], Types.Transaction{
				Date:  date,
				Delta: amount,
				Memo:  "Transfer to Savings",
				From:  Types.Checking,
				To:    Types.Savings,
			})
		} else {
			ledger[date] = append(ledger[date], Types.Transaction{
				Date:  date,
				Delta: amount,
				Memo:  "Transfer from Savings",
				From:  Types.Savings,
				To:    Types.Checking,
			})
		}
	}

	for _, expense := range expenses {
		occurrances := expense.Schedule.FindOccurrances(startDay, endDay)
		for _, date := range occurrances {
			ledger[date] = append(ledger[date], Types.Transaction{
				Date:  date,
				Delta: expense.Amount,
				Memo:  fmt.Sprintf("Transfer from Savings for: %s", expense.Name),
				From:  Types.Savings,
				To:    Types.Checking,
			},
				Types.Transaction{
					Date:  date,
					Delta: -1 * expense.Amount,
					Memo:  fmt.Sprintf("Expense: %s", expense.Name),
					From:  Types.Checking,
					To:    Types.External,
				},
			)
		}
	}
	return ledger, idealDiscretionary
}

// Simulate ...
func Simulate(startDay time.Time, endDay time.Time, ledger map[time.Time][]Types.Transaction, shouldPrintOutput bool) (accounts map[Types.Account]float64, averageSpending float64) {
	simulatedSpending := 0.
	numDays := 0.
	accounts = map[Types.Account]float64{Types.External: 0., Types.Checking: 0., Types.Savings: 0.}

	currentDate := startDay
	for {
		if currentDate.After(endDay) {
			break
		}

		transactions := ledger[currentDate]
		for _, transaction := range transactions {
			if transaction.Memo == simulatedSpendingMemo {
				simulatedSpending += math.Abs(transaction.Delta)
			}
			accounts[transaction.From] -= math.Abs(transaction.Delta)
			accounts[transaction.To] += math.Abs(transaction.Delta)
			if shouldPrintOutput {
				fmt.Printf("%s | %7.2f | %7.2f\n", transaction.ToString(), accounts[Types.Checking], accounts[Types.Savings])
			}
		}
		numDays += 1.
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return accounts, simulatedSpending / numDays
}

/*





































*/
