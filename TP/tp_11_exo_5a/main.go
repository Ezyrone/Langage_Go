package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var items = []Item{
	{ID: "1", Name: "Clavier mécanique", Description: "Clavier RGB switches blue"},
	{ID: "2", Name: "Souris gaming", Description: "Souris sans fil 16000 DPI"},
	{ID: "3", Name: "Écran 27 pouces", Description: "Écran IPS 144Hz QHD"},
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func extractID(path string) string {
	trimmed := strings.TrimPrefix(path, "/items/")
	return strings.TrimSuffix(trimmed, "/")
}

func findItemIndex(id string) int {
	for i, item := range items {
		if item.ID == id {
			return i
		}
	}
	return -1
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/items" && r.URL.Path != "/items/" {
		itemByIDHandler(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getItemsHandler(w, r)
	case http.MethodPost:
		createItemHandler(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

func itemByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItemHandler(w, r)
	case http.MethodPut:
		updateItemHandler(w, r)
	case http.MethodDelete:
		deleteItemHandler(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

func getItemsHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, items)
}

func getItemHandler(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	idx := findItemIndex(id)
	if idx == -1 {
		respondError(w, http.StatusNotFound, "Item non trouvé")
		return
	}
	respondJSON(w, http.StatusOK, items[idx])
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		respondError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	newItem.ID = uuid.New().String()
	items = append(items, newItem)
	respondJSON(w, http.StatusCreated, newItem)
}

func updateItemHandler(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	idx := findItemIndex(id)
	if idx == -1 {
		respondError(w, http.StatusNotFound, "Item non trouvé")
		return
	}

	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		respondError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	items[idx].Name = updatedItem.Name
	items[idx].Description = updatedItem.Description
	respondJSON(w, http.StatusOK, items[idx])
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	idx := findItemIndex(id)
	if idx == -1 {
		respondError(w, http.StatusNotFound, "Item non trouvé")
		return
	}

	items = append(items[:idx], items[idx+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/items", itemsHandler)
	mux.HandleFunc("/items/", itemsHandler)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	fmt.Println("Endpoints :")
	fmt.Println("  GET    /items       - Liste tous les items")
	fmt.Println("  GET    /items/{id}  - Récupère un item par ID")
	fmt.Println("  POST   /items       - Crée un nouvel item")
	fmt.Println("  PUT    /items/{id}  - Met à jour un item")
	fmt.Println("  DELETE /items/{id}  - Supprime un item")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
