package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/ndowns/even_challenge/Types"
	"github.com/ndowns/even_challenge/money"
	"github.com/stretchr/testify/assert"
)

func TestBasicHappyPath(t *testing.T) {
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

	plan, idealSpending := Plan(startDay, endDay, incomes, expenses)
	accounts, avgSimulatedSpending, err := Simulate(startDay, endDay, plan, false)
	assert.Equal(t, nil, err)

	assert.Equal(t, money.New(0.), accounts[Types.External])
	assert.Equal(t, money.New(0.), accounts[Types.Checking])
	assert.Equal(t, money.New(0.), accounts[Types.Savings])

	fmt.Println(idealSpending, avgSimulatedSpending)
	//assert.InDelta(t, avgSimulatedSpending/idealSpending, 1., 0.05)
}

func TestInsolvent(t *testing.T) {
	startDay, _ := time.Parse(Types.DateFormat, "2015.08.01")
	endDay, _ := time.Parse(Types.DateFormat, "2015.08.31")

	incomes := []Types.Income{
		Types.Income{
			Amount:   money.New(500.),
			Name:     "Philz",
			Schedule: Types.Schedule{Period: Types.BiMonthly},
		},
	}

	expenses := []Types.Expense{
		Types.Expense{
			Amount:   money.New(42.34),
			Name:     "Utilities",
			Schedule: Types.Schedule{Period: Types.Monthly, Date: 25},
		},
		Types.Expense{
			Amount:   money.New(1200.),
			Name:     "Rent",
			Schedule: Types.Schedule{Period: Types.Monthly, Date: 28},
		},
		Types.Expense{
			Amount:   money.New(40.),
			Name:     "Crossfit",
			Schedule: Types.Schedule{Period: Types.Weekly, Weekday: time.Tuesday},
		},
	}

	plan, idealSpending := Plan(startDay, endDay, incomes, expenses)
	assert.Equal(t, 0, len(plan))
	assert.Equal(t, money.New(0.), idealSpending)
}

func TestOneTimePayment(t *testing.T) {
	startDay, _ := time.Parse(Types.DateFormat, "2015.08.01")
	endDay, _ := time.Parse(Types.DateFormat, "2015.12.31")

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

	christmas, _ := time.Parse(Types.DateFormat, "2015.12.25")
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
			Amount:   money.New(600.),
			Name:     "Vacation",
			Schedule: Types.Schedule{Period: Types.OneTime, Time: christmas},
		},
	}

	plan, idealSpending := Plan(startDay, endDay, incomes, expenses)
	accounts, avgSimulatedSpending, err := Simulate(startDay, endDay, plan, false)
	assert.Equal(t, nil, err)

	assert.Equal(t, money.New(0.), accounts[Types.External])
	assert.Equal(t, money.New(0.), accounts[Types.Checking])
	assert.Equal(t, money.New(0.), accounts[Types.Savings])

	fmt.Println(idealSpending, avgSimulatedSpending)
	//assert.InDelta(t, avgSimulatedSpending/idealSpending, 1., 0.05)
}
