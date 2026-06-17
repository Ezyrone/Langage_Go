package main

import (
	"errors"
	"fmt"
	"math"
)

type Point struct {
	X float64
	Y float64
}

func (p Point) DistanceTo(other Point) float64 {
	return math.Sqrt(math.Pow(other.X-p.X, 2) + math.Pow(other.Y-p.Y, 2))
}

type Rectangle struct {
	Min Point
	Max Point
}

func (r Rectangle) Width() float64 {
	return r.Max.X - r.Min.X
}

func (r Rectangle) Height() float64 {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle) Area() float64 {
	return r.Width() * r.Height()
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width() + r.Height())
}

// Pointer receiver: Move modifies the rectangle in place
func (r *Rectangle) Move(dx, dy float64) {
	r.Min.X += dx
	r.Min.Y += dy
	r.Max.X += dx
	r.Max.Y += dy
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle[Min(%.2f, %.2f) Max(%.2f, %.2f)] width=%.2f, height=%.2f",
		r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, r.Width(), r.Height())
}

type Circle struct {
	Center Point
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.Radius
}

// Pointer receiver: Scale modifies the radius in place
func (c *Circle) Scale(factor float64) {
	c.Radius *= factor
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle[center(%.2f, %.2f), radius=%.2f]", c.Center.X, c.Center.Y, c.Radius)
}

func NewRectangle(min, max Point) (Rectangle, error) {
	if max.X-min.X <= 0 || max.Y-min.Y <= 0 {
		return Rectangle{}, errors.New("rectangle dimensions must be positive")
	}
	return Rectangle{Min: min, Max: max}, nil
}

func NewCircle(center Point, radius float64) (Circle, error) {
	if radius <= 0 {
		return Circle{}, errors.New("circle radius must be positive")
	}
	return Circle{Center: center, Radius: radius}, nil
}

func main() {

	// Exercise 1
	fmt.Println("Exercise 1: Point and Rectangle")

	p1 := Point{X: 1, Y: 2}
	p2 := Point{X: 4, Y: 6}
	fmt.Printf("p1 = (%.2f, %.2f)\n", p1.X, p1.Y)
	fmt.Printf("p2 = (%.2f, %.2f)\n", p2.X, p2.Y)
	fmt.Printf("Distance p1 -> p2: %.4f\n", p1.DistanceTo(p2))

	rect := Rectangle{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 5, Y: 3},
	}
	fmt.Println("\n" + rect.String())
	fmt.Printf("Width: %.2f\n", rect.Width())
	fmt.Printf("Height: %.2f\n", rect.Height())
	fmt.Printf("Area: %.2f\n", rect.Area())
	fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())

	fmt.Println("\nMoving rectangle by (2, 1)...")
	rect.Move(2, 1)
	fmt.Println(rect)
	fmt.Printf("New Min: (%.2f, %.2f), New Max: (%.2f, %.2f)\n", rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)

	// Exercise 2
	fmt.Println("\nExercise 2: Circle")

	circle := Circle{Center: Point{X: 3, Y: 4}, Radius: 5}
	fmt.Println(circle)
	fmt.Printf("Area: %.4f\n", circle.Area())
	fmt.Printf("Circumference: %.4f\n", circle.Circumference())

	fmt.Println("\nScale x2...")
	circle.Scale(2)
	fmt.Println(circle)
	fmt.Printf("New radius: %.2f\n", circle.Radius)
	fmt.Printf("Area: %.4f\n", circle.Area())
	fmt.Printf("Circumference: %.4f\n", circle.Circumference())

	// Exercise 3
	fmt.Println("\nExercise 3: Constructors with validation")

	r, err := NewRectangle(Point{X: 0, Y: 0}, Point{X: 4, Y: 3})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Valid rectangle:", r)
	}

	_, err = NewRectangle(Point{X: 5, Y: 5}, Point{X: 2, Y: 1})
	if err != nil {
		fmt.Println("Invalid rectangle:", err)
	}

	c, err := NewCircle(Point{X: 0, Y: 0}, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Valid circle:", c)
	}

	_, err = NewCircle(Point{X: 0, Y: 0}, -1)
	if err != nil {
		fmt.Println("Invalid circle:", err)
	}

	// Exercise 3.3: Reflection on receivers
	fmt.Println("\nReflection on receivers:")
	fmt.Println("- Value receiver: the method works on a copy of the struct.")
	fmt.Println("  Used for Area(), Perimeter(), Width(), Height(), DistanceTo(), Circumference()")
	fmt.Println("  because these methods read data without modifying it.")
	fmt.Println("- Pointer receiver: the method works on the original instance.")
	fmt.Println("  Used for Move() and Scale() because they modify the struct's state.")
}
