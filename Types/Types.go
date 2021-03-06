package Types

import (
	"fmt"
	"time"

	"github.com/n8downs/even_challenge/money"
)

// Period ...
type Period int

// test
const (
	Monthly Period = iota
	BiMonthly
	Weekly
	BiWeekly
	OneTime
)

// Account ...
type Account int

// test
const (
	External Account = iota
	Checking
	Savings
)

func (p Period) String() string {
	switch p {
	case Monthly:
		return "Monthly"
	case BiMonthly:
		return "BiMonthly"
	case Weekly:
		return "Weekly"
	case BiWeekly:
		return "BiWeekly"
	case OneTime:
		return "OneTime"
	default:
		return "???"
	}
}

// DateFormat
const (
	DateFormat = "2006.01.02"
)

// Income ...
type Income struct {
	Name     string
	Amount   money.Money
	Schedule Schedule
}

// Expense ...
type Expense struct {
	Name     string
	Amount   money.Money
	Schedule Schedule
}

// Transaction ...
type Transaction struct {
	Date  time.Time
	Delta money.Money
	Memo  string
	From  Account
	To    Account
}

func (t Transaction) String() string {
	return fmt.Sprintf("%10s | %-40s | %15s", t.Date.Format(DateFormat), t.Memo, t.Delta.String())
}

// Schedule ...
type Schedule struct {
	Period  Period
	Weekday time.Weekday
	Date    int
	Time    time.Time
}

// FindRealOccurrances ...
func (s Schedule) FindRealOccurrances(from, to time.Time) (occurrances []time.Time) {
	switch s.Period {
	case Monthly:
		{
			fromYear, fromMonth, _ := from.Date()
			occ := time.Date(fromYear, fromMonth, s.Date, 0, 0, 0, 0, time.UTC)
			if occ.Before(from) {
				occ = occ.AddDate(0, 1, 0)
			}
			for {
				if occ.After(to) {
					break
				}
				occurrances = append(occurrances, occ)
				occ = occ.AddDate(0, 1, 0)
			}
		}
	case BiMonthly:
		{
			// TODO: handle the case where the 1st or 15th falls on a weekend/holiday
			// (payment should happen on the Friday before in that case)
			fromYear, fromMonth, _ := from.Date()
			occ1 := time.Date(fromYear, fromMonth, 1, 0, 0, 0, 0, time.UTC)
			occ15 := time.Date(fromYear, fromMonth, 15, 0, 0, 0, 0, time.UTC)
			if occ1.Before(from) {
				occ1.AddDate(0, 1, 0)
			}
			if occ15.Before(from) {
				occ15.AddDate(0, 1, 0)
			}

			for {
				if occ1.After(to) {
					break
				}
				occurrances = append(occurrances, occ1)
				occ1 = occ1.AddDate(0, 1, 0)

				if occ15.After(to) {
					break
				}
				occurrances = append(occurrances, occ15)
				occ15 = occ15.AddDate(0, 1, 0)
			}
		}
	case Weekly:
		{
			occ := from
			for {
				if occ.Weekday() == s.Weekday {
					break
				}
				occ = occ.AddDate(0, 0, 1)
			}

			for {
				if occ.After(to) {
					break
				}
				occurrances = append(occurrances, occ)
				occ = occ.AddDate(0, 0, 7)
			}
		}
	case BiWeekly:
		{
			// TODO: probably want a way to specify that the next occurance will not
			// be the first we come across, but rather the second from the start
			occ := from
			for {
				if occ.Weekday() == s.Weekday {
					break
				}
				occ = occ.AddDate(0, 0, 1)
			}

			for {
				if occ.After(to) {
					break
				}
				occurrances = append(occurrances, occ)
				occ = occ.AddDate(0, 0, 14)
			}
		}
	case OneTime:
		{
			if !(s.Time.Before(from) || s.Time.After(to)) {
				occurrances = append(occurrances, s.Time)
			}
		}
	}

	return
}

// FindVirtualOccurrances ...
func (e Expense) FindVirtualOccurrances(from, to time.Time) map[time.Time]money.Money {
	occurrances := map[time.Time]money.Money{}
	switch e.Schedule.Period {
	/*
		  // XXXnrd: It seems like this should help, but in practice, it doesn't seem to.
		  case Monthly:
			{
				v := Schedule{
					Period: Weekly,
				}
				realDates := e.Schedule.FindRealOccurrances(from, to)
				currentDate := from
				for _, realDate := range realDates {
					dates := v.FindRealOccurrances(currentDate, realDate)
					amounts := e.Amount.Divide(int64(len(dates)))
					for i := 0; i < len(dates); i++ {
						occurrances[dates[i]] = amounts[i]
					}
					currentDate = realDate.AddDate(0, 0, 1)
				}
			}
	*/
	case OneTime:
		{
			realDates := e.Schedule.FindRealOccurrances(from, to)
			if len(realDates) == 0 {
				return occurrances
			}
			v := Schedule{
				Period: Weekly,
			}
			dates := v.FindRealOccurrances(from, e.Schedule.Time)
			amounts := e.Amount.Divide(int64(len(dates)))
			for i := 0; i < len(dates); i++ {
				occurrances[dates[i]] = amounts[i]
			}
		}
	default:
		{
			for _, date := range e.Schedule.FindRealOccurrances(from, to) {
				occurrances[date] = e.Amount
			}
		}
	}
	return occurrances
}
