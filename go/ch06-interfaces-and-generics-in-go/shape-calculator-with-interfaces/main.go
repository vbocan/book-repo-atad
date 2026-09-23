// Exercise 6.2: Shape Calculator with Interfaces
//
// This program works with different shapes through a single Shape
// interface, demonstrating polymorphism: Circle, Rectangle, and Triangle
// each satisfy Shape independently, with no shared base type and no
// declared "implements" relationship.
package main

import (
	"cmp"
	"fmt"
	"math"
	"slices"
)

// Shape is satisfied by any type that can report its area and perimeter
// and describe itself as a string.
type Shape interface {
	Area() float64
	Perimeter() float64
	String() string
}

// Circle is a shape defined by its radius.
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

// Rectangle is a shape defined by its width and height.
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(width=%.2f, height=%.2f)", r.Width, r.Height)
}

// Triangle is a shape defined by the lengths of its three sides. Area is
// computed with Heron's formula, which needs only the side lengths.
type Triangle struct {
	A, B, C float64
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

func (t Triangle) Area() float64 {
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) String() string {
	return fmt.Sprintf("Triangle(sides=%.2f,%.2f,%.2f)", t.A, t.B, t.C)
}

// totalArea sums the area of every shape in shapes.
func totalArea(shapes []Shape) float64 {
	var total float64
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// largest returns the shape with the greatest area. It panics if shapes
// is empty, since there is then no largest shape to return.
func largest(shapes []Shape) Shape {
	if len(shapes) == 0 {
		panic("largest: no shapes")
	}
	best := shapes[0]
	for _, s := range shapes[1:] {
		if s.Area() > best.Area() {
			best = s
		}
	}
	return best
}

func main() {
	shapes := []Shape{
		Circle{Radius: 4},
		Rectangle{Width: 5, Height: 3},
		Triangle{A: 6, B: 8, C: 10},
	}

	fmt.Println("Shapes:")
	for _, s := range shapes {
		fmt.Printf("  %s -> area=%.2f, perimeter=%.2f\n", s, s.Area(), s.Perimeter())
	}

	fmt.Printf("Total area: %.2f\n", totalArea(shapes))

	big := largest(shapes)
	fmt.Printf("Largest shape: %s (area=%.2f)\n", big, big.Area())

	// Sort by area ascending. slices.SortFunc takes a comparison function
	// that returns a negative number, zero, or a positive number, exactly
	// like cmp.Compare.
	slices.SortFunc(shapes, func(a, b Shape) int {
		return cmp.Compare(a.Area(), b.Area())
	})

	fmt.Println("Sorted by area (ascending):")
	for _, s := range shapes {
		fmt.Printf("  %s -> area=%.2f\n", s, s.Area())
	}
}
