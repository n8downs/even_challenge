package Types

import "time"
import "fmt"
import "math"

const (
	dateFormat = "2006.01.02"
)

// Schedule ...
type Schedule struct {
	Period string
	Day    string
}

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
}

// ToString ... sigh
func (t Transaction) ToString() string {
	var d string
	if t.Delta >= 0 {
		d = fmt.Sprintf(" %.2f ", t.Delta)
	} else {
		d = fmt.Sprintf("(%.2f)", math.Abs(t.Delta))
	}
	return fmt.Sprintf("%10s | %-40s | %15s", t.Date.Format(dateFormat), t.Memo, d)
}
