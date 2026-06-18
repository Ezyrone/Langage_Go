package main

import "fmt"


const (
	PI       = 3.14159
	APP_NAME = "Go Manager"
	LAUNCH_YEAR = 2023
)

const (
	Monday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func main() {
	fmt.Println("Exercise 1")

	var userName string = "Jory"
	var userAge int = 30
	var isConnected bool = true
	var accountBalance float64 = 1542.75

	fmt.Println("Name:", userName)
	fmt.Println("Age:", userAge)
	fmt.Println("Connected:", isConnected)
	fmt.Println("Balance:", accountBalance)

	fmt.Println("\nExercise 2")

	city := "Grenoble"
	zipCode := 38000
	discountRate := 15.5

	fmt.Printf("City: %v (type: %T)\n", city, city)
	fmt.Printf("Zip code: %v (type: %T)\n", zipCode, zipCode)
	fmt.Printf("Discount rate: %v (type: %T)\n", discountRate, discountRate)

	fmt.Println("\nExercise 3")

	radius := 10.5
	circumference := 2 * PI * radius
	fmt.Printf("Circumference of a circle with radius %.2f: %.4f\n", radius, circumference)

	fmt.Println("PI:", PI)
	fmt.Println("Application:", APP_NAME)
	fmt.Println("Launch year:", LAUNCH_YEAR)

	fmt.Println("\nExercise 4")

	oldAge := userAge
	userAge = userAge + 1
	fmt.Printf("Birthday! Old age: %d -> New age: %d\n", oldAge, userAge)

	var message string
	fmt.Printf("message (uninitialized): \"%s\" (zero value of string = empty string)\n", message)

	var counter int
	fmt.Printf("counter (uninitialized): %d (zero value of int = 0)\n", counter)

	fmt.Println("\nBonus")

	var a, b, c int = 1, 2, 3
	fmt.Println("Multiple declaration:", a, b, c)

	fmt.Println("Monday =", Monday, "| Wednesday =", Wednesday, "| Sunday =", Sunday)

	var integer int = 42
	var decimal float64 = 3.14

	result := float64(integer) + decimal
	fmt.Printf("Conversion: %d (int) + %.2f (float64) = %.2f\n", integer, decimal, result)
}
