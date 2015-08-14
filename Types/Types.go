package Types

// Schedule ...
type Schedule struct {
	period string
	day    string
}

// Income ...
type Income struct {
	name     string
	amount   float64
	schedule Schedule
}

// Expense ...
type Expense struct {
	name     string
	amount   float64
	schedule Schedule
}
