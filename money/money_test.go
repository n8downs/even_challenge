package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	m := New(100.)
	assert.Equal(t, "100.00", m.String())

	m = New(123.45)
	assert.Equal(t, "123.45", m.String())
}

func TestTruncation(t *testing.T) {
	m := New(123.456)
	assert.Equal(t, "123.45", m.String())
}

func TestNegative(t *testing.T) {
	m := New(-42.99)
	assert.Equal(t, "(42.99)", m.String())
}

func TestAdd(t *testing.T) {
	m := New(123.45)
	n := New(456.78)

	o := m.Add(n)
	assert.Equal(t, "580.23", o.String())
}

func TestAddNegative(t *testing.T) {
	m := New(123.45)
	n := New(-456.78)

	o := m.Add(n)
	assert.Equal(t, "(333.33)", o.String())
}

func TestSubtract(t *testing.T) {
	m := New(123.45)
	n := New(456.78)

	o := m.Subtract(n)
	assert.Equal(t, "(333.33)", o.String())
}

func TestSubtractNegative(t *testing.T) {
	m := New(123.45)
	n := New(-456.78)

	o := m.Subtract(n)
	assert.Equal(t, "580.23", o.String())
}

func TestMultiply(t *testing.T) {
	m := New(1.03)
	n := m.Multiply(3.141592654)

	assert.Equal(t, "3.23", n.String())
}

func TestDivide(t *testing.T) {
	m := New(100.03)
	parts := m.Divide(2)

	assert.Equal(t, []Money{New(50.01), New(50.02)}, parts)
}

func TestDivideByOne(t *testing.T) {
	m := New(100.03)
	parts := m.Divide(1)

	assert.Equal(t, []Money{New(100.03)}, parts)
}

func TestDivideByZero(t *testing.T) {
	m := New(100.03)
	parts := m.Divide(0)

	assert.Equal(t, []Money{}, parts)
}

func TestDivideByNegative(t *testing.T) {
	m := New(100.03)
	parts := m.Divide(-2)

	// XXXnrd: there may be a more sensical behavior here, but I'm not quite sure what it is
	assert.Equal(t, []Money{}, parts)
}

func TestGreaterThan(t *testing.T) {
	m := New(1.99)

	assert.True(t, m.GreaterThan(New(1.98)))
	assert.False(t, m.GreaterThan(New(2.00)))
	assert.False(t, m.GreaterThan(New(1.99)))
}

func TestEqual(t *testing.T) {
	m := New(1.99)

	assert.False(t, m.EqualTo(New(1.98)))
	assert.True(t, m.EqualTo(New(1.99)))
	assert.False(t, m.EqualTo(New(2.00)))
}

func TestMax(t *testing.T) {
	assert.Equal(t, New(2.00), Max(New(2.00), New(1.99)))
	assert.Equal(t, New(2.00), Max(New(1.99), New(2.00)))
	assert.Equal(t, New(2.00), Max(New(1.98), New(1.99), New(2.00)))
}

func TestAbs(t *testing.T) {
	assert.Equal(t, New(2.00), New(2.00).Abs())
	assert.Equal(t, New(2.00), New(-2.00).Abs())
}

/*































*/
