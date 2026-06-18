package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type PersonneV1 struct {
	Nom   string
	Age   int
	Email string
	Actif bool
}

type Personne struct {
	Nom        string `json:"full_name"`
	Age        int    `json:"age_in_years"`
	Email      string `json:"contact_email,omitempty"`
	Actif      bool   `json:"is_active"`
	MotDePasse string `json:"-"`
}

type Produit struct {
	ID      int     `json:"product_id"`
	Nom     string  `json:"item_name"`
	Prix    float64 `json:"unit_price"`
	EnStock bool    `json:"in_stock"`
}

type Editeur struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Livre struct {
	ID             int       `json:"book_id"`
	Titre          string    `json:"title"`
	Auteur         string    `json:"author_name"`
	AnneePubli     int       `json:"publication_year"`
	Genres         []string  `json:"genres,omitempty"`
	ISBN           string    `json:"isbn_code,omitempty"`
	EstDisponible  bool      `json:"is_available"`
	EditeurInfo    *Editeur  `json:"publisher_info,omitempty"`
	DateAjout      UnixTime  `json:"date_ajout,omitempty"`
}

type UnixTime struct {
	time.Time
}

func (t UnixTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Unix())
}

func (t *UnixTime) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	t.Time = time.Unix(timestamp, 0)
	return nil
}

func prettyJSON(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("Erreur: %v", err)
	}
	return string(b)
}

func main() {
	fmt.Println("=== Exercice 1 : Sérialisation Basique ===")

	p1 := PersonneV1{
		Nom:   "Alice Dupont",
		Age:   30,
		Email: "alice.dupont@example.com",
		Actif: true,
	}

	jsonBytes, err := json.Marshal(p1)
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}
	fmt.Println(string(jsonBytes))

	fmt.Println("\n=== Exercice 2 : Struct Tags ===")

	pAvecEmail := Personne{
		Nom:        "Alice Dupont",
		Age:        30,
		Email:      "alice.dupont@example.com",
		Actif:      true,
		MotDePasse: "secret123",
	}

	pSansEmail := Personne{
		Nom:        "Bob Martin",
		Age:        25,
		Actif:      false,
		MotDePasse: "mdp456",
	}

	fmt.Println("Avec email:")
	fmt.Println(prettyJSON(pAvecEmail))

	fmt.Println("\nSans email:")
	fmt.Println(prettyJSON(pSansEmail))

	fmt.Println("\n=== Exercice 3 : Désérialisation ===")

	jsonString := `{
		"product_id": 101,
		"item_name": "Clavier Mécanique",
		"unit_price": 79.99,
		"in_stock": true
	}`

	var produit Produit
	err = json.Unmarshal([]byte(jsonString), &produit)
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	fmt.Printf("ID:      %d\n", produit.ID)
	fmt.Printf("Nom:     %s\n", produit.Nom)
	fmt.Printf("Prix:    %.2f\n", produit.Prix)
	fmt.Printf("EnStock: %v\n", produit.EnStock)

	fmt.Println("\n=== Exercice 4 : Gestion des Erreurs ===")

	malformedJSON := `{
		"product_id": 102,
		"item_name": "Souris Gaming",
		"unit_price": 49.99,
		"in_stock": true,
	`

	var p2 Produit
	err = json.Unmarshal([]byte(malformedJSON), &p2)
	fmt.Printf("JSON malformé -> Erreur: %v\n", err)

	wrongTypeJSON := `{
		"product_id": "103",
		"item_name": "Écran UltraWide",
		"unit_price": 399.99,
		"in_stock": true
	}`

	var p3 Produit
	err = json.Unmarshal([]byte(wrongTypeJSON), &p3)
	fmt.Printf("Type incorrect -> Erreur: %v\n", err)

	fmt.Println("\n=== Exercice 5A : Sérialisation de Livres ===")

	livre1 := Livre{
		ID:            1,
		Titre:         "Le Petit Prince",
		Auteur:        "Antoine de Saint-Exupéry",
		AnneePubli:    1943,
		Genres:        []string{"Fiction", "Jeunesse", "Philosophie"},
		ISBN:          "978-2-07-040850-4",
		EstDisponible: true,
		EditeurInfo:   &Editeur{Name: "Gallimard", Location: "Paris"},
		DateAjout:     UnixTime{time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)},
	}

	livre2 := Livre{
		ID:            2,
		Titre:         "1984",
		Auteur:        "George Orwell",
		AnneePubli:    1949,
		EstDisponible: false,
	}

	fmt.Println("Livre 1 (complet):")
	json1 := prettyJSON(livre1)
	fmt.Println(json1)

	fmt.Println("\nLivre 2 (genres et ISBN vides):")
	json2 := prettyJSON(livre2)
	fmt.Println(json2)

	fmt.Println("\n=== Exercice 5B : Désérialisation de Livres ===")

	var livreDecode1 Livre
	err = json.Unmarshal([]byte(json1), &livreDecode1)
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}
	fmt.Printf("Livre 1 décodé: %+v\n", livreDecode1)
	fmt.Printf("  DateAjout: %s\n", livreDecode1.DateAjout.Format(time.RFC3339))
	fmt.Printf("  Editeur: %+v\n", livreDecode1.EditeurInfo)

	var livreDecode2 Livre
	err = json.Unmarshal([]byte(json2), &livreDecode2)
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}
	fmt.Printf("Livre 2 décodé: %+v\n", livreDecode2)
	fmt.Printf("  Genres: %v (nil: %v)\n", livreDecode2.Genres, livreDecode2.Genres == nil)
	fmt.Printf("  ISBN: %q\n", livreDecode2.ISBN)
}
