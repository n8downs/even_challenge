package money

import (
	"fmt"
	"math"
)

// New ...
func New(m float64) Money {
	return Money{int64(m * 100)}
}

// Money ...
type Money struct {
	pennies int64
}

func (m Money) String() string {
	if m.pennies < 0 {
		return fmt.Sprintf("(%.2f)", math.Abs(float64(m.pennies)/100.))
	}
	return fmt.Sprintf("%.2f", float64(m.pennies)/100.)
}

// Add ...
func (m Money) Add(n Money) Money {
	return Money{m.pennies + n.pennies}
}

// Subtract ...
func (m Money) Subtract(n Money) Money {
	return Money{m.pennies - n.pennies}
}

// Multiply ...
func (m Money) Multiply(n float64) Money {
	return Money{int64(float64(m.pennies) * n)}
}

// Divide ...
func (m Money) Divide(n int64) []Money {
	if n <= 0 {
		return []Money{}
	}
	equalPart := m.pennies / n
	remaining := m.pennies
	result := []Money{}
	for i := int64(0); i < n-1; i++ {
		result = append(result, Money{equalPart})
		remaining -= equalPart
	}
	result = append(result, Money{remaining})
	return result
}

// GreaterThan ...
func (m Money) GreaterThan(right Money) bool {
	return m.pennies > right.pennies
}

// EqualTo ...
func (m Money) EqualTo(right Money) bool {
	return m.pennies == right.pennies
}

// Max ...
func Max(args ...Money) Money {
	max := args[0]
	for _, m := range args {
		if m.GreaterThan(max) {
			max = m
		}
	}
	return max
}

// Abs ...
func (m Money) Abs() Money {
	return Money{int64(math.Abs(float64(m.pennies)))}
}

// Float ...
func (m Money) Float() float64 {
	return float64(m.pennies) / 100.
}
