package main

import (
	"testing"
	"time"

	"github.com/ndowns/even_challenge/Types"
	"github.com/stretchr/testify/assert"
)

func TestBasicHappyPath(t *testing.T) {
	startDay, _ := time.Parse(Types.DateFormat, "2015.08.01")
	endDay, _ := time.Parse(Types.DateFormat, "2015.08.31")

	incomes := []Types.Income{
		Types.Income{
			Amount:   500.,
			Name:     "Philz",
			Schedule: Types.Schedule{Period: Types.BiMonthly},
		},
		Types.Income{
			Amount:   175.,
			Name:     "Mission Cliffs",
			Schedule: Types.Schedule{Period: Types.BiWeekly, Weekday: time.Thursday},
		},
	}

	expenses := []Types.Expense{
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
	}

	plan, idealSpending := Plan(startDay, endDay, incomes, expenses)
	accounts, avgSimulatedSpending := Simulate(startDay, endDay, plan, false)

	assert.InDelta(t, accounts[Types.External], 0, 0.001)
	assert.InDelta(t, accounts[Types.Checking], 0, 0.001)
	assert.InDelta(t, accounts[Types.Savings], 0, 0.001)

	assert.InDelta(t, avgSimulatedSpending/idealSpending, 1., 0.05)
}

func TestInsolvent(t *testing.T) {
	startDay, _ := time.Parse(Types.DateFormat, "2015.08.01")
	endDay, _ := time.Parse(Types.DateFormat, "2015.08.31")

	incomes := []Types.Income{
		Types.Income{
			Amount:   500.,
			Name:     "Philz",
			Schedule: Types.Schedule{Period: Types.BiMonthly},
		},
	}

	expenses := []Types.Expense{
		Types.Expense{
			Amount:   42.34,
			Name:     "Utilities",
			Schedule: Types.Schedule{Period: Types.Monthly, Date: 25},
		},
		Types.Expense{
			Amount:   1200.,
			Name:     "Rent",
			Schedule: Types.Schedule{Period: Types.Monthly, Date: 28},
		},
		Types.Expense{
			Amount:   40.,
			Name:     "Crossfit",
			Schedule: Types.Schedule{Period: Types.Weekly, Weekday: time.Tuesday},
		},
	}

	plan, idealSpending := Plan(startDay, endDay, incomes, expenses)
	assert.Equal(t, 0, len(plan))
	assert.Equal(t, 0., idealSpending)
}

func TestOneTimePayment(t *testing.T) {
	startDay, _ := time.Parse(Types.DateFormat, "2015.08.01")
	endDay, _ := time.Parse(Types.DateFormat, "2015.12.31")

	incomes := []Types.Income{
		Types.Income{
			Amount:   500.,
			Name:     "Philz",
			Schedule: Types.Schedule{Period: Types.BiMonthly},
		},
		Types.Income{
			Amount:   175.,
			Name:     "Mission Cliffs",
			Schedule: Types.Schedule{Period: Types.BiWeekly, Weekday: time.Thursday},
		},
	}

	christmas, _ := time.Parse(Types.DateFormat, "2015.12.25")
	expenses := []Types.Expense{
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
			Amount:   600.,
			Name:     "Vacation",
			Schedule: Types.Schedule{Period: Types.OneTime, Time: christmas},
		},
	}

	plan, idealSpending := Plan(startDay, endDay, incomes, expenses)
	accounts, avgSimulatedSpending := Simulate(startDay, endDay, plan, false)

	assert.InDelta(t, accounts[Types.External], 0, 0.001)
	assert.InDelta(t, accounts[Types.Checking], 0, 0.001)
	assert.InDelta(t, accounts[Types.Savings], 0, 0.001)

	assert.InDelta(t, avgSimulatedSpending/idealSpending, 1., 0.05)
}
