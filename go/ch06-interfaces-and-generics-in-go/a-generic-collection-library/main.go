// Exercise 6.1: A Generic Collection Library
//
// SortedList[T] is a type-safe collection that keeps its elements in
// ascending order at all times. It works for any type that satisfies
// cmp.Ordered, such as int, float64, or string, and supports insertion
// at the correct position, membership tests, and removal.
package main

import (
	"cmp"
	"fmt"
)

// SortedList keeps its elements sorted in ascending order. It is generic
// over any element type T that satisfies cmp.Ordered.
type SortedList[T cmp.Ordered] struct {
	items []T
}

// NewSortedList creates an empty SortedList for element type T.
func NewSortedList[T cmp.Ordered]() *SortedList[T] {
	return &SortedList[T]{}
}

// search performs a binary search for value. It returns the index where
// value was found, or the index where it should be inserted to keep the
// list sorted, along with whether the value was found.
func (l *SortedList[T]) search(value T) (index int, found bool) {
	lo, hi := 0, len(l.items)
	for lo < hi {
		mid := lo + (hi-lo)/2
		switch c := cmp.Compare(l.items[mid], value); {
		case c == 0:
			return mid, true
		case c < 0:
			lo = mid + 1
		default:
			hi = mid
		}
	}
	return lo, false
}

// Insert adds value to the list, keeping the elements in ascending order.
// Duplicate values are allowed and are inserted next to any existing
// equal element.
func (l *SortedList[T]) Insert(value T) {
	index, _ := l.search(value)
	var zero T
	l.items = append(l.items, zero)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = value
}

// Contains reports whether value is present in the list.
func (l *SortedList[T]) Contains(value T) bool {
	_, found := l.search(value)
	return found
}

// Remove deletes the first occurrence of value from the list and reports
// whether an element was actually removed.
func (l *SortedList[T]) Remove(value T) bool {
	index, found := l.search(value)
	if !found {
		return false
	}
	l.items = append(l.items[:index], l.items[index+1:]...)
	return true
}

// Values returns a copy of the list's elements in ascending order. It is
// a copy so callers cannot corrupt the list's internal ordering.
func (l *SortedList[T]) Values() []T {
	values := make([]T, len(l.items))
	copy(values, l.items)
	return values
}

// String implements fmt.Stringer so a *SortedList prints the same way a
// plain slice of its elements would.
func (l *SortedList[T]) String() string {
	return fmt.Sprint(l.items)
}

func main() {
	ints := NewSortedList[int]()
	for _, n := range []int{5, 3, 8, 1, 9, 4, 7} {
		ints.Insert(n)
	}
	fmt.Println("Integers:", ints)
	fmt.Println("Contains 4:", ints.Contains(4))
	ints.Remove(3)
	fmt.Println("After removing 3:", ints)

	strs := NewSortedList[string]()
	for _, s := range []string{"cherry", "apple", "date", "banana"} {
		strs.Insert(s)
	}
	fmt.Println("Strings:", strs)
}
