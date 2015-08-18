package Types

import (
	"fmt"
	"math"
	"time"
)

// Period ...
type Period int

// test
const (
	Monthly Period = iota
	BiMonthly
	Weekly
	BiWeekly
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
	Amount   float64
	Schedule Schedule
}

// Expense ...
type Expense struct {
	Name     string
	Amount   float64
	Schedule Schedule
}

// Transaction ...
type Transaction struct {
	Date  time.Time
	Delta float64
	Memo  string
	From  Account
	To    Account
}

// ToString ...
func (t Transaction) ToString() string {
	var d string
	if t.Delta >= 0 {
		d = fmt.Sprintf(" %.2f ", t.Delta)
	} else {
		d = fmt.Sprintf("(%.2f)", math.Abs(t.Delta))
	}
	return fmt.Sprintf("%10s | %-40s | %15s", t.Date.Format(DateFormat), t.Memo, d)
}

// Schedule ...
type Schedule struct {
	Period  Period
	Weekday time.Weekday
	Date    int
}

// FindOccurrances ...
func (s Schedule) FindOccurrances(from, to time.Time) (occurrances []time.Time) {
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
	}

	return
}

/*

























*/
