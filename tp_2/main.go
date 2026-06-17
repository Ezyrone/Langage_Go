package main

import (
	"errors"
	"fmt"
)

// Exercice 1
func CalculerStatistiquesBase(nombres ...int) (int, int, float64) {
	if len(nombres) == 0 {
		return 0, 0, 0.0
	}

	somme := 0
	for _, n := range nombres {
		somme += n
	}

	count := len(nombres)
	moyenne := float64(somme) / float64(count)

	return somme, count, moyenne
}

// Exercice 2
func CalculerStatistiquesCompletes(nombres ...float64) (float64, float64, float64, float64, int, error) {
	if len(nombres) == 0 {
		return 0, 0, 0, 0, 0, errors.New("aucun argument fourni")
	}

	min := nombres[0]
	max := nombres[0]
	somme := 0.0

	for _, n := range nombres {
		somme += n
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	count := len(nombres)
	moyenne := somme / float64(count)

	return min, max, somme, moyenne, count, nil
}

// Exercice 3
func AnalyserDonneesCapteur(releves ...float64) (float64, float64, float64, int, int, error) {
	var valides []float64
	invalides := 0

	for _, r := range releves {
		if r > 0.0 && r <= 100.0 {
			valides = append(valides, r)
		} else {
			invalides++
		}
	}

	if len(valides) == 0 {
		return 0, 0, 0, 0, invalides, errors.New("aucun relevé valide trouvé")
	}

	min, max, _, moyenne, validCount, _ := CalculerStatistiquesCompletes(valides...)

	return min, max, moyenne, validCount, invalides, nil
}

func main() {
	// ===== Exercice 1 =====
	fmt.Println("=== Exercice 1 : Statistiques de Base ===")

	somme, count, moyenne := CalculerStatistiquesBase(10, 20, 30, 40)
	fmt.Printf("Somme: %d, Count: %d, Moyenne: %.2f\n", somme, count, moyenne)

	sommeVide, countVide, moyenneVide := CalculerStatistiquesBase()
	fmt.Printf("Somme (vide): %d, Count (vide): %d, Moyenne (vide): %.2f\n", sommeVide, countVide, moyenneVide)

	sommeUn, countUn, moyenneUn := CalculerStatistiquesBase(42)
	fmt.Printf("Somme (un seul): %d, Count (un seul): %d, Moyenne (un seul): %.2f\n", sommeUn, countUn, moyenneUn)

	// ===== Exercice 2 =====
	fmt.Println("\n=== Exercice 2 : Statistiques Complètes ===")

	min, max, sum, avg, cnt, err := CalculerStatistiquesCompletes(1.5, 2.8, 0.7, 3.1)
	if err != nil {
		fmt.Println("Erreur:", err)
	} else {
		fmt.Printf("Min: %.2f, Max: %.2f, Somme: %.2f, Moyenne: %.2f, Count: %d\n", min, max, sum, avg, cnt)
	}

	_, _, _, _, _, errVide := CalculerStatistiquesCompletes()
	if errVide != nil {
		fmt.Println("Erreur pour arguments vides:", errVide)
	}

	// ===== Exercice 3 =====
	fmt.Println("\n=== Exercice 3 : Analyse de Données de Capteur ===")

	minTemp, maxTemp, avgTemp, validCnt, invalidCnt, errCapteur := AnalyserDonneesCapteur(22.5, 23.1, -5.0, 101.0, 21.9, 0.0, 24.0)
	if errCapteur != nil {
		fmt.Println("Erreur d'analyse:", errCapteur)
	} else {
		fmt.Printf("Temp Min: %.2f, Max: %.2f, Moyenne: %.2f, Valides: %d, Invalides: %d\n", minTemp, maxTemp, avgTemp, validCnt, invalidCnt)
	}

	_, _, _, _, _, errToutInvalide := AnalyserDonneesCapteur(-10.0, 105.0, 0.0)
	if errToutInvalide != nil {
		fmt.Println("Erreur pour données toutes invalides:", errToutInvalide)
	}
}
