package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/ndowns/even_challenge/Types"
	"github.com/ndowns/even_challenge/money"
)

const simulatedSpendingMemo = "  -Simulated Spending-"

func main() {
	startDay, _ := time.Parse(Types.DateFormat, "2015.08.01")
	endDay, _ := time.Parse(Types.DateFormat, "2015.08.31")

	incomes := []Types.Income{
		Types.Income{
			Amount:   money.New(500.),
			Name:     "Philz",
			Schedule: Types.Schedule{Period: Types.BiMonthly},
		},
		Types.Income{
			Amount:   money.New(175.),
			Name:     "Mission Cliffs",
			Schedule: Types.Schedule{Period: Types.BiWeekly, Weekday: time.Thursday},
		},
	}

	expenses := []Types.Expense{
		Types.Expense{
			Amount:   money.New(42.34),
			Name:     "Utilities",
			Schedule: Types.Schedule{Period: Types.Monthly, Date: 25},
		},
		Types.Expense{
			Amount:   money.New(400.),
			Name:     "Rent",
			Schedule: Types.Schedule{Period: Types.Monthly, Date: 28},
		},
		Types.Expense{
			Amount:   money.New(40.),
			Name:     "Crossfit",
			Schedule: Types.Schedule{Period: Types.Weekly, Weekday: time.Tuesday},
		},
	}

	plan, _ := Plan(startDay, endDay, incomes, expenses)

	accounts, _, err := Simulate(startDay, endDay, plan, true)
	if err != nil {
		fmt.Println(err, accounts)
	}
}

// Plan ...
func Plan(
	startDay time.Time,
	endDay time.Time,
	incomes []Types.Income,
	expenses []Types.Expense,
) (map[time.Time][]Types.Transaction, money.Money) {
	ledger := map[time.Time][]Types.Transaction{}
	totalIncome := money.New(0.)
	totalExpenses := money.New(0.)

	incomeTotals := map[time.Time]money.Money{}
	savingsPlan := map[time.Time]money.Money{}
	expenseTotals := map[time.Time]money.Money{}

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
			totalIncome = totalIncome.Add(income.Amount)
			incomeTotals[date] = incomeTotals[date].Add(income.Amount)
			savingsPlan[date] = money.New(0.)
		}
	}

	for _, expense := range expenses {
		occurrances := expense.Schedule.FindOccurrances(startDay, endDay)
		for _, date := range occurrances {
			totalExpenses = totalExpenses.Add(expense.Amount)
			expenseTotals[date] = expenseTotals[date].Add(expense.Amount)
		}
	}

	if totalExpenses.GreaterThan(totalIncome) {
		return map[time.Time][]Types.Transaction{}, money.New(0.)
	}

	discretionaryDivided := totalIncome.Subtract(totalExpenses).Divide(int64(endDay.Sub(startDay).Hours() / 24))
	idealDiscretionary := discretionaryDivided[0]

	currentDate := startDay
	runningPlan := money.New(0.)
	for {
		if currentDate.After(endDay) {
			break
		}

		if !incomeTotals[currentDate].EqualTo(money.New(0.)) {
			nextIncomeDate := currentDate
			upcomingExpenses := money.New(0.)
			for {
				if nextIncomeDate.After(endDay) {
					break
				}

				upcomingExpenses = upcomingExpenses.Add(expenseTotals[nextIncomeDate])
				nextIncomeDate = nextIncomeDate.AddDate(0, 0, 1)
				if !incomeTotals[nextIncomeDate].EqualTo(money.New(0.)) {
					break
				}
			}
			daysUntilNextIncome := int64(nextIncomeDate.Sub(currentDate).Hours() / 24)
			mustTransfer := upcomingExpenses.Subtract(runningPlan)
			idealTransfer := incomeTotals[currentDate].Subtract(idealDiscretionary.Multiply(float64(daysUntilNextIncome)))
			transfer := money.Max(mustTransfer, idealTransfer)
			runningPlan = runningPlan.Add(transfer).Subtract(upcomingExpenses)

			savingsPlan[currentDate] = transfer.Multiply(-1.)
			actualDiscretionaries := incomeTotals[currentDate].Subtract(transfer).Divide(daysUntilNextIncome)

			simulatedSpendingDate := currentDate
			for _, discretionaryAmount := range actualDiscretionaries {
				ledger[simulatedSpendingDate] = append(ledger[simulatedSpendingDate], Types.Transaction{
					From:  Types.Checking,
					To:    Types.External,
					Memo:  simulatedSpendingMemo,
					Date:  simulatedSpendingDate,
					Delta: discretionaryAmount,
				})
				simulatedSpendingDate = simulatedSpendingDate.AddDate(0, 0, 1)
			}
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	for date, amount := range savingsPlan {
		if money.New(0.).GreaterThan(amount) {
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
					Delta: expense.Amount.Multiply(-1.),
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
func Simulate(
	startDay time.Time,
	endDay time.Time,
	ledger map[time.Time][]Types.Transaction,
	shouldPrintOutput bool,
) (accounts map[Types.Account]money.Money, averageSpending money.Money, err error) {
	simulatedSpending := money.New(0.)
	numDays := int64(0)
	accounts = map[Types.Account]money.Money{Types.External: money.New(0.), Types.Checking: money.New(0.), Types.Savings: money.New(0.)}

	currentDate := startDay
	for {
		if currentDate.After(endDay) {
			break
		}

		transactions := ledger[currentDate]
		for _, transaction := range transactions {
			if transaction.Memo == simulatedSpendingMemo {
				simulatedSpending = simulatedSpending.Add(transaction.Delta)
			}
			accounts[transaction.From] = accounts[transaction.From].Subtract(transaction.Delta.Abs())
			accounts[transaction.To] = accounts[transaction.To].Add(transaction.Delta.Abs())

			if shouldPrintOutput {
				fmt.Printf("%s | %9s | %9s\n", transaction.String(), accounts[Types.Checking].String(), accounts[Types.Savings].String())
			}

			if money.New(0.).GreaterThan(accounts[Types.Checking]) || money.New(0.).GreaterThan(accounts[Types.Savings]) {
				return accounts, money.New(0.), errors.New("Balance dipped below zero!")
			}
		}
		numDays += 1.
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	divided := simulatedSpending.Divide(numDays)

	return accounts, divided[0], nil
}

/*





































*/
