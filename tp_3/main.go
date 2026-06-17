package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Produit struct {
	ID        int
	Nom       string
	Prix      float64
	Categorie string
}

func categorieExiste(nom string, categories []string) bool {
	for _, c := range categories {
		if c == nom {
			return true
		}
	}
	return false
}

func supprimerCategorie(nom string, categories []string) []string {
	for i, c := range categories {
		if c == nom {
			return append(categories[:i], categories[i+1:]...)
		}
	}
	return categories
}

func obtenirProduit(id int, inventaire map[int]Produit, stock map[int]int) (Produit, int, bool) {
	produit, existe := inventaire[id]
	if !existe {
		return Produit{}, 0, false
	}
	return produit, stock[id], true
}

func vendreProduit(id int, quantite int, stock map[int]int) bool {
	if stock[id] < quantite {
		return false
	}
	stock[id] -= quantite
	return true
}

func reapprovisionnerProduit(id int, quantite int, stock map[int]int) {
	stock[id] += quantite
}

func listerProduitsParCategorie(categorie string, inventaire map[int]Produit, produitsParCategorie map[string][]int) {
	ids, existe := produitsParCategorie[categorie]
	if !existe || len(ids) == 0 {
		fmt.Printf("  Aucun produit dans la catégorie \"%s\"\n", categorie)
		return
	}
	for _, id := range ids {
		p := inventaire[id]
		fmt.Printf("  [%d] %s - %.2f€\n", p.ID, p.Nom, p.Prix)
	}
}

func trierParPrix(inventaire map[int]Produit, croissant bool) []Produit {
	produits := make([]Produit, 0, len(inventaire))
	for _, p := range inventaire {
		produits = append(produits, p)
	}
	sort.Slice(produits, func(i, j int) bool {
		if croissant {
			return produits[i].Prix < produits[j].Prix
		}
		return produits[i].Prix > produits[j].Prix
	})
	return produits
}

func valeurStockCategorie(categorie string, inventaire map[int]Produit, stock map[int]int, produitsParCategorie map[string][]int) float64 {
	total := 0.0
	ids := produitsParCategorie[categorie]
	for _, id := range ids {
		total += inventaire[id].Prix * float64(stock[id])
	}
	return total
}

func afficherInventaire(inventaire map[int]Produit, stock map[int]int) {
	for id, p := range inventaire {
		fmt.Printf("  [%d] %s | %.2f€ | Catégorie: %s | Stock: %d\n", id, p.Nom, p.Prix, p.Categorie, stock[id])
	}
}

func main() {

	fmt.Println("PARTIE 1 : ")

	categories := []string{"Électronique", "Vêtements", "Livres"}
	categories = append(categories, "Alimentation", "Sport")
	fmt.Println("Catégories:", categories)

	fmt.Println("\"Livres\" existe ?", categorieExiste("Livres", categories))
	fmt.Println("\"Musique\" existe ?", categorieExiste("Musique", categories))

	categories = supprimerCategorie("Vêtements", categories)
	fmt.Println("Après suppression de \"Vêtements\":", categories)

	categories = supprimerCategorie("Musique", categories)
	fmt.Println("Après suppression de \"Musique\" (inexistante):", categories)

	fmt.Printf("Longueur: %d, Capacité: %d\n", len(categories), cap(categories))
	fmt.Println("-> len = nombre d'éléments actuels, cap = taille du tableau sous-jacent.")
	fmt.Println("   La capacité double quand le slice dépasse sa capacité lors d'un append.")

	fmt.Println("\n PARTIE 2 : ")

	inventaireProduits := map[int]Produit{
		1: {ID: 1, Nom: "Laptop", Prix: 999.99, Categorie: "Électronique"},
		2: {ID: 2, Nom: "Go Programming", Prix: 39.90, Categorie: "Livres"},
		3: {ID: 3, Nom: "Ballon de foot", Prix: 24.99, Categorie: "Sport"},
	}

	stockProduits := map[int]int{
		1: 15,
		2: 50,
		3: 30,
	}

	p := inventaireProduits[1]
	p.Prix = 899.99
	inventaireProduits[1] = p
	fmt.Println("Prix du Laptop modifié à 899.99€")

	stockProduits[2] = 45
	fmt.Println("Stock de \"Go Programming\" mis à jour à 45")

	fmt.Println("\nInventaire complet:")
	afficherInventaire(inventaireProduits, stockProduits)

	fmt.Println("\nRecherche produit ID=2:")
	if prod, qty, ok := obtenirProduit(2, inventaireProduits, stockProduits); ok {
		fmt.Printf("  Trouvé: %s, Stock: %d\n", prod.Nom, qty)
	}

	fmt.Println("Recherche produit ID=99:")
	if _, _, ok := obtenirProduit(99, inventaireProduits, stockProduits); !ok {
		fmt.Println("  Produit non trouvé")
	}

	delete(inventaireProduits, 3)
	delete(stockProduits, 3)
	fmt.Println("\nProduit ID=3 supprimé")
	if _, _, ok := obtenirProduit(3, inventaireProduits, stockProduits); !ok {
		fmt.Println("  Vérification: produit ID=3 n'existe plus")
	}

	fmt.Println("\nOpérations de stock:")
	fmt.Printf("  Stock Laptop avant vente: %d\n", stockProduits[1])
	if vendreProduit(1, 3, stockProduits) {
		fmt.Printf("  Vente de 3 Laptops réussie. Stock: %d\n", stockProduits[1])
	}
	if !vendreProduit(1, 100, stockProduits) {
		fmt.Println("  Vente de 100 Laptops échouée: stock insuffisant")
	}

	reapprovisionnerProduit(1, 10, stockProduits)
	fmt.Printf("  Réapprovisionnement de 10 Laptops. Stock: %d\n", stockProduits[1])

	fmt.Println("\n========== PARTIE 3 : Combinaison & Performance ==========")

	inventaireProduits[3] = Produit{ID: 3, Nom: "Ballon de foot", Prix: 24.99, Categorie: "Sport"}
	stockProduits[3] = 30
	inventaireProduits[4] = Produit{ID: 4, Nom: "Casque audio", Prix: 59.99, Categorie: "Électronique"}
	stockProduits[4] = 20

	produitsParCategorie := make(map[string][]int)
	for _, p := range inventaireProduits {
		produitsParCategorie[p.Categorie] = append(produitsParCategorie[p.Categorie], p.ID)
	}

	fmt.Println("Produits par catégorie \"Électronique\":")
	listerProduitsParCategorie("Électronique", inventaireProduits, produitsParCategorie)
	fmt.Println("Produits par catégorie \"Livres\":")
	listerProduitsParCategorie("Livres", inventaireProduits, produitsParCategorie)
	fmt.Println("Produits par catégorie \"Jardinage\":")
	listerProduitsParCategorie("Jardinage", inventaireProduits, produitsParCategorie)

	fmt.Println("\n--- Performance (100 000 produits) ---")
	n := 100_000
	cats := []string{"Électronique", "Livres", "Sport", "Alimentation", "Mode"}

	start := time.Now()
	bigInv := make(map[int]Produit)
	bigStock := make(map[int]int)
	for i := 0; i < n; i++ {
		bigInv[i] = Produit{
			ID:        i,
			Nom:       fmt.Sprintf("Produit_%d", i),
			Prix:      rand.Float64() * 500,
			Categorie: cats[rand.Intn(len(cats))],
		}
		bigStock[i] = rand.Intn(200)
	}
	dureesSansPrealloc := time.Since(start)
	fmt.Printf("Ajout SANS pré-allocation : %v\n", dureesSansPrealloc)

	start = time.Now()
	bigInv2 := make(map[int]Produit, n)
	bigStock2 := make(map[int]int, n)
	for i := 0; i < n; i++ {
		bigInv2[i] = Produit{
			ID:        i,
			Nom:       fmt.Sprintf("Produit_%d", i),
			Prix:      rand.Float64() * 500,
			Categorie: cats[rand.Intn(len(cats))],
		}
		bigStock2[i] = rand.Intn(200)
	}
	dureesAvecPrealloc := time.Since(start)
	fmt.Printf("Ajout AVEC pré-allocation : %v\n", dureesAvecPrealloc)
	fmt.Println("-> La pré-allocation évite les rehash/réallocations internes de la map.")

	start = time.Now()
	for i := 0; i < 10_000; i++ {
		_ = bigInv[rand.Intn(n)]
	}
	fmt.Printf("10 000 recherches par ID : %v\n", time.Since(start))

	start = time.Now()
	total := 0.0
	for _, p := range bigInv {
		total += p.Prix
	}
	fmt.Printf("Itération sur %d éléments : %v\n", n, time.Since(start))
	fmt.Println("-> L'accès par clé dans une map est O(1) en moyenne (hash table).")
	fmt.Println("   L'itération est O(n). La pré-allocation réduit le temps d'insertion")
	fmt.Println("   car la map n'a pas besoin de grandir dynamiquement.")

	fmt.Println("\n========== BONUS ==========")

	fmt.Println("Produits triés par prix croissant:")
	for _, p := range trierParPrix(inventaireProduits, true) {
		fmt.Printf("  %s - %.2f€\n", p.Nom, p.Prix)
	}

	fmt.Println("Produits triés par prix décroissant:")
	for _, p := range trierParPrix(inventaireProduits, false) {
		fmt.Printf("  %s - %.2f€\n", p.Nom, p.Prix)
	}

	valeur := valeurStockCategorie("Électronique", inventaireProduits, stockProduits, produitsParCategorie)
	fmt.Printf("Valeur totale du stock \"Électronique\": %.2f€\n", valeur)
}
