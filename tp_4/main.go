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

func (r *Rectangle) Move(dx, dy float64) {
	r.Min.X += dx
	r.Min.Y += dy
	r.Max.X += dx
	r.Max.Y += dy
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle[Min(%.2f, %.2f) Max(%.2f, %.2f)] largeur=%.2f, hauteur=%.2f",
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

// Receiver pointeur : Scale modifie le rayon en place
func (c *Circle) Scale(factor float64) {
	c.Radius *= factor
}

func (c Circle) String() string {
	return fmt.Sprintf("Cercle[centre(%.2f, %.2f), rayon=%.2f]", c.Center.X, c.Center.Y, c.Radius)
}

func NewRectangle(min, max Point) (Rectangle, error) {
	if max.X-min.X <= 0 || max.Y-min.Y <= 0 {
		return Rectangle{}, errors.New("les dimensions du rectangle doivent être positives")
	}
	return Rectangle{Min: min, Max: max}, nil
}

func NewCircle(center Point, radius float64) (Circle, error) {
	if radius <= 0 {
		return Circle{}, errors.New("le rayon du cercle doit être positif")
	}
	return Circle{Center: center, Radius: radius}, nil
}

func main() {

	// Exercice 1
	fmt.Println("Exercice 1 : Point et Rectangle")

	p1 := Point{X: 1, Y: 2}
	p2 := Point{X: 4, Y: 6}
	fmt.Printf("p1 = (%.2f, %.2f)\n", p1.X, p1.Y)
	fmt.Printf("p2 = (%.2f, %.2f)\n", p2.X, p2.Y)
	fmt.Printf("Distance p1 -> p2 : %.4f\n", p1.DistanceTo(p2))

	rect := Rectangle{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 5, Y: 3},
	}
	fmt.Println("\n" + rect.String())
	fmt.Printf("Largeur: %.2f\n", rect.Width())
	fmt.Printf("Hauteur: %.2f\n", rect.Height())
	fmt.Printf("Surface: %.2f\n", rect.Area())
	fmt.Printf("Périmètre: %.2f\n", rect.Perimeter())

	fmt.Println("\nDéplacement du rectangle de (2, 1)...")
	rect.Move(2, 1)
	fmt.Println(rect)
	fmt.Printf("Nouveau Min: (%.2f, %.2f), Nouveau Max: (%.2f, %.2f)\n", rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)

	// Exercice 2
	fmt.Println("\nExercice 2 : Cercle")

	cercle := Circle{Center: Point{X: 3, Y: 4}, Radius: 5}
	fmt.Println(cercle)
	fmt.Printf("Surface: %.4f\n", cercle.Area())
	fmt.Printf("Circonférence: %.4f\n", cercle.Circumference())

	fmt.Println("\nScale x2...")
	cercle.Scale(2)
	fmt.Println(cercle)
	fmt.Printf("Nouveau rayon: %.2f\n", cercle.Radius)
	fmt.Printf("Surface: %.4f\n", cercle.Area())
	fmt.Printf("Circonférence: %.4f\n", cercle.Circumference())

	// Exercice 3
	fmt.Println("\nExercice 3 : Constructeurs avec validation")

	r, err := NewRectangle(Point{X: 0, Y: 0}, Point{X: 4, Y: 3})
	if err != nil {
		fmt.Println("Erreur:", err)
	} else {
		fmt.Println("Rectangle valide:", r)
	}

	_, err = NewRectangle(Point{X: 5, Y: 5}, Point{X: 2, Y: 1})
	if err != nil {
		fmt.Println("Rectangle invalide:", err)
	}

	c, err := NewCircle(Point{X: 0, Y: 0}, 3)
	if err != nil {
		fmt.Println("Erreur:", err)
	} else {
		fmt.Println("Cercle valide:", c)
	}

	_, err = NewCircle(Point{X: 0, Y: 0}, -1)
	if err != nil {
		fmt.Println("Cercle invalide:", err)
	}

	// Exercice 3.3 : Réflexion sur les receivers
	fmt.Println("\nRéflexion sur les receivers:")
	fmt.Println("- Receiver de valeur : la méthode travaille sur une copie de la struct.")
	fmt.Println("  Utilisé pour Area(), Perimeter(), Width(), Height(), DistanceTo(), Circumference()")
	fmt.Println("  car ces méthodes lisent les données sans les modifier.")
	fmt.Println("- Receiver de pointeur : la méthode travaille sur l'instance originale.")
	fmt.Println("  Utilisé pour Move() et Scale() car elles modifient l'état de la struct.")
}
