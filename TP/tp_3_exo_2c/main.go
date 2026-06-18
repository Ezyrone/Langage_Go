package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Product struct {
	ID       int
	Name     string
	Price    float64
	Category string
}

func categoryExists(name string, categories []string) bool {
	for _, c := range categories {
		if c == name {
			return true
		}
	}
	return false
}

func removeCategory(name string, categories []string) []string {
	for i, c := range categories {
		if c == name {
			return append(categories[:i], categories[i+1:]...)
		}
	}
	return categories
}

func getProduct(id int, inventory map[int]Product, stock map[int]int) (Product, int, bool) {
	product, exists := inventory[id]
	if !exists {
		return Product{}, 0, false
	}
	return product, stock[id], true
}

func sellProduct(id int, quantity int, stock map[int]int) bool {
	if stock[id] < quantity {
		return false
	}
	stock[id] -= quantity
	return true
}

func restockProduct(id int, quantity int, stock map[int]int) {
	stock[id] += quantity
}

func listProductsByCategory(category string, inventory map[int]Product, productsByCategory map[string][]int) {
	ids, exists := productsByCategory[category]
	if !exists || len(ids) == 0 {
		fmt.Printf("  No products in category \"%s\"\n", category)
		return
	}
	for _, id := range ids {
		p := inventory[id]
		fmt.Printf("  [%d] %s - $%.2f\n", p.ID, p.Name, p.Price)
	}
}

func sortByPrice(inventory map[int]Product, ascending bool) []Product {
	products := make([]Product, 0, len(inventory))
	for _, p := range inventory {
		products = append(products, p)
	}
	sort.Slice(products, func(i, j int) bool {
		if ascending {
			return products[i].Price < products[j].Price
		}
		return products[i].Price > products[j].Price
	})
	return products
}

func categoryStockValue(category string, inventory map[int]Product, stock map[int]int, productsByCategory map[string][]int) float64 {
	total := 0.0
	ids := productsByCategory[category]
	for _, id := range ids {
		total += inventory[id].Price * float64(stock[id])
	}
	return total
}

func displayInventory(inventory map[int]Product, stock map[int]int) {
	for id, p := range inventory {
		fmt.Printf("  [%d] %s | $%.2f | Category: %s | Stock: %d\n", id, p.Name, p.Price, p.Category, stock[id])
	}
}

func main() {

	fmt.Println("PART 1: Slices")

	categories := []string{"Electronics", "Clothing", "Books"}
	categories = append(categories, "Food", "Sports")
	fmt.Println("Categories:", categories)

	fmt.Println("\"Books\" exists?", categoryExists("Books", categories))
	fmt.Println("\"Music\" exists?", categoryExists("Music", categories))

	categories = removeCategory("Clothing", categories)
	fmt.Println("After removing \"Clothing\":", categories)

	categories = removeCategory("Music", categories)
	fmt.Println("After removing \"Music\" (non-existent):", categories)

	fmt.Printf("Length: %d, Capacity: %d\n", len(categories), cap(categories))
	fmt.Println("-> len = current number of elements, cap = size of the underlying array.")
	fmt.Println("   Capacity doubles when the slice exceeds its capacity during an append.")

	fmt.Println("\nPART 2: Maps")

	productInventory := map[int]Product{
		1: {ID: 1, Name: "Laptop", Price: 999.99, Category: "Electronics"},
		2: {ID: 2, Name: "Go Programming", Price: 39.90, Category: "Books"},
		3: {ID: 3, Name: "Soccer Ball", Price: 24.99, Category: "Sports"},
	}

	productStock := map[int]int{
		1: 15,
		2: 50,
		3: 30,
	}

	p := productInventory[1]
	p.Price = 899.99
	productInventory[1] = p
	fmt.Println("Laptop price updated to $899.99")

	productStock[2] = 45
	fmt.Println("\"Go Programming\" stock updated to 45")

	fmt.Println("\nFull inventory:")
	displayInventory(productInventory, productStock)

	fmt.Println("\nSearching product ID=2:")
	if prod, qty, ok := getProduct(2, productInventory, productStock); ok {
		fmt.Printf("  Found: %s, Stock: %d\n", prod.Name, qty)
	}

	fmt.Println("Searching product ID=99:")
	if _, _, ok := getProduct(99, productInventory, productStock); !ok {
		fmt.Println("  Product not found")
	}

	delete(productInventory, 3)
	delete(productStock, 3)
	fmt.Println("\nProduct ID=3 deleted")
	if _, _, ok := getProduct(3, productInventory, productStock); !ok {
		fmt.Println("  Verification: product ID=3 no longer exists")
	}

	fmt.Println("\nStock operations:")
	fmt.Printf("  Laptop stock before sale: %d\n", productStock[1])
	if sellProduct(1, 3, productStock) {
		fmt.Printf("  Sold 3 Laptops. Stock: %d\n", productStock[1])
	}
	if !sellProduct(1, 100, productStock) {
		fmt.Println("  Sale of 100 Laptops failed: insufficient stock")
	}

	restockProduct(1, 10, productStock)
	fmt.Printf("  Restocked 10 Laptops. Stock: %d\n", productStock[1])

	fmt.Println("\nPART 3: Combination & Performance")

	productInventory[3] = Product{ID: 3, Name: "Soccer Ball", Price: 24.99, Category: "Sports"}
	productStock[3] = 30
	productInventory[4] = Product{ID: 4, Name: "Headphones", Price: 59.99, Category: "Electronics"}
	productStock[4] = 20

	productsByCategory := make(map[string][]int)
	for _, p := range productInventory {
		productsByCategory[p.Category] = append(productsByCategory[p.Category], p.ID)
	}

	fmt.Println("Products in category \"Electronics\":")
	listProductsByCategory("Electronics", productInventory, productsByCategory)
	fmt.Println("Products in category \"Books\":")
	listProductsByCategory("Books", productInventory, productsByCategory)
	fmt.Println("Products in category \"Gardening\":")
	listProductsByCategory("Gardening", productInventory, productsByCategory)

	fmt.Println("\n--- Performance (100,000 products) ---")
	n := 100_000
	cats := []string{"Electronics", "Books", "Sports", "Food", "Fashion"}

	start := time.Now()
	bigInv := make(map[int]Product)
	bigStock := make(map[int]int)
	for i := 0; i < n; i++ {
		bigInv[i] = Product{
			ID:       i,
			Name:     fmt.Sprintf("Product_%d", i),
			Price:    rand.Float64() * 500,
			Category: cats[rand.Intn(len(cats))],
		}
		bigStock[i] = rand.Intn(200)
	}
	withoutPrealloc := time.Since(start)
	fmt.Printf("Insert WITHOUT pre-allocation: %v\n", withoutPrealloc)

	start = time.Now()
	bigInv2 := make(map[int]Product, n)
	bigStock2 := make(map[int]int, n)
	for i := 0; i < n; i++ {
		bigInv2[i] = Product{
			ID:       i,
			Name:     fmt.Sprintf("Product_%d", i),
			Price:    rand.Float64() * 500,
			Category: cats[rand.Intn(len(cats))],
		}
		bigStock2[i] = rand.Intn(200)
	}
	withPrealloc := time.Since(start)
	fmt.Printf("Insert WITH pre-allocation: %v\n", withPrealloc)
	fmt.Println("-> Pre-allocation avoids internal rehash/reallocations in the map.")

	start = time.Now()
	for i := 0; i < 10_000; i++ {
		_ = bigInv[rand.Intn(n)]
	}
	fmt.Printf("10,000 lookups by ID: %v\n", time.Since(start))

	start = time.Now()
	total := 0.0
	for _, p := range bigInv {
		total += p.Price
	}
	fmt.Printf("Iteration over %d elements: %v\n", n, time.Since(start))
	fmt.Println("-> Key access in a map is O(1) on average (hash table).")
	fmt.Println("   Iteration is O(n). Pre-allocation reduces insertion time")
	fmt.Println("   because the map doesn't need to grow dynamically.")

	fmt.Println("\nBONUS")

	fmt.Println("Products sorted by price (ascending):")
	for _, p := range sortByPrice(productInventory, true) {
		fmt.Printf("  %s - $%.2f\n", p.Name, p.Price)
	}

	fmt.Println("Products sorted by price (descending):")
	for _, p := range sortByPrice(productInventory, false) {
		fmt.Printf("  %s - $%.2f\n", p.Name, p.Price)
	}

	value := categoryStockValue("Electronics", productInventory, productStock, productsByCategory)
	fmt.Printf("Total stock value for \"Electronics\": $%.2f\n", value)
}
