package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func effectuerTache(id int, resultChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Goroutine %d: Début de la tâche...\n", id)
	duree := time.Duration(50+rand.Intn(451)) * time.Millisecond
	time.Sleep(duree)
	fmt.Printf("Goroutine %d: Tâche terminée. (durée: %v)\n", id, duree)
	resultChan <- fmt.Sprintf("Goroutine %d a terminé avec succès.", id)
}

func travailleur(id int, taches <-chan int, resultats chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for tacheID := range taches {
		fmt.Printf("Travailleur %d: traite la tâche %d...\n", id, tacheID)
		duree := time.Duration(50+rand.Intn(451)) * time.Millisecond
		time.Sleep(duree)
		resultats <- fmt.Sprintf("Travailleur %d: tâche %d terminée (durée: %v)", id, tacheID, duree)
	}
	fmt.Printf("Travailleur %d: plus de tâches, arrêt.\n", id)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// =============================================
	// Exercice 1 : Lancement sans synchronisation
	// =============================================
	fmt.Println("=== Exercice 1 : Sans synchronisation ===")
	for i := 1; i <= 5; i++ {
		go func(id int) {
			fmt.Printf("Goroutine %d: Début de la tâche...\n", id)
			duree := time.Duration(50+rand.Intn(451)) * time.Millisecond
			time.Sleep(duree)
			fmt.Printf("Goroutine %d: Tâche terminée.\n", id)
		}(i)
	}
	fmt.Println("Toutes les goroutines lancées.")

	// Réponse : Sans synchronisation, main() se termine immédiatement après
	// avoir lancé les goroutines. Le programme s'arrête avant que les goroutines
	// n'aient le temps de finir leur travail. On voit peu ou pas de messages
	// "Tâche terminée".

	// Petit délai pour laisser les goroutines de l'exo 1 s'afficher partiellement
	time.Sleep(600 * time.Millisecond)

	// =============================================
	// Exercice 2 : Synchronisation avec WaitGroup
	// =============================================
	fmt.Println("\n=== Exercice 2 : Avec sync.WaitGroup ===")

	var wg2 sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg2.Add(1)
		go func(id int) {
			defer wg2.Done()
			fmt.Printf("Goroutine %d: Début de la tâche...\n", id)
			duree := time.Duration(50+rand.Intn(451)) * time.Millisecond
			time.Sleep(duree)
			fmt.Printf("Goroutine %d: Tâche terminée.\n", id)
		}(i)
	}
	fmt.Println("Toutes les goroutines lancées.")
	wg2.Wait()
	fmt.Println("Toutes les goroutines ont terminé leur exécution.")

	// Réponse : Oui, toutes les goroutines terminent maintenant leur travail.
	// wg.Wait() bloque main() jusqu'à ce que chaque goroutine appelle wg.Done().

	// =============================================
	// Exercice 3 : Communication avec les canaux
	// =============================================
	fmt.Println("\n=== Exercice 3 : Canaux de résultats ===")

	var wg3 sync.WaitGroup
	resultChan := make(chan string, 5)

	for i := 1; i <= 5; i++ {
		wg3.Add(1)
		go effectuerTache(i, resultChan, &wg3)
	}

	fmt.Println("Toutes les goroutines lancées.")
	wg3.Wait()
	close(resultChan)

	fmt.Println("\nRésultats reçus :")
	for msg := range resultChan {
		fmt.Println(" -", msg)
	}

	// Réponse : L'ordre des résultats ne correspond pas nécessairement à l'ordre
	// des IDs. Les goroutines s'exécutent en concurrence avec des durées aléatoires,
	// donc celle qui finit en premier envoie son résultat en premier dans le canal.

	// =============================================
	// Exercice 4 : Pool de travailleurs
	// =============================================
	fmt.Println("\n=== Exercice 4 : Pool de travailleurs (3 workers, 10 tâches) ===")

	const nbTravailleurs = 3
	const nbTaches = 10

	taches := make(chan int, nbTaches)
	resultats := make(chan string, nbTaches)
	var wg4 sync.WaitGroup

	for w := 1; w <= nbTravailleurs; w++ {
		wg4.Add(1)
		go travailleur(w, taches, resultats, &wg4)
	}

	for t := 1; t <= nbTaches; t++ {
		taches <- t
	}
	close(taches)

	wg4.Wait()
	close(resultats)

	fmt.Println("\nRésultats du pool :")
	for msg := range resultats {
		fmt.Println(" -", msg)
	}

	fmt.Println("\nTous les travailleurs ont terminé.")

	// Réponse : Les tâches sont réparties entre les 3 travailleurs. L'ordre de
	// traitement dépend de la disponibilité de chaque travailleur. Avec plus de
	// travailleurs, le temps total diminue car davantage de tâches sont traitées
	// en parallèle. Avec 3 travailleurs pour 10 tâches, chacun traite ~3-4 tâches.
}
