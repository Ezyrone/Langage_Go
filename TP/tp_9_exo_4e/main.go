package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const nbGoroutines = 100
const incrementsParGoroutine = 1000

var compteur int
var mu sync.Mutex

func incrementerCompteurNonSynchro(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < incrementsParGoroutine; i++ {
		compteur++
	}
}

func incrementerCompteurSynchro(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < incrementsParGoroutine; i++ {
		mu.Lock()
		compteur++
		mu.Unlock()
	}
}

var compteurAtomic int64

func incrementerCompteurAtomic(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < incrementsParGoroutine; i++ {
		atomic.AddInt64(&compteurAtomic, 1)
	}
}

func main() {
	attendu := nbGoroutines * incrementsParGoroutine

	// === Étape 1 : Sans synchronisation ===
	fmt.Println("=== Étape 1 : Compteur NON synchronisé ===")
	compteur = 0
	var wg1 sync.WaitGroup

	start := time.Now()
	for i := 0; i < nbGoroutines; i++ {
		wg1.Add(1)
		go incrementerCompteurNonSynchro(&wg1)
	}
	wg1.Wait()
	duree1 := time.Since(start)

	fmt.Printf("Résultat : %d (attendu : %d)\n", compteur, attendu)
	fmt.Printf("Correct  : %v\n", compteur == attendu)
	fmt.Printf("Durée    : %v\n", duree1)

	// === Étape 2 : Avec Mutex ===
	fmt.Println("\n=== Étape 2 : Compteur synchronisé avec Mutex ===")
	compteur = 0
	var wg2 sync.WaitGroup

	start = time.Now()
	for i := 0; i < nbGoroutines; i++ {
		wg2.Add(1)
		go incrementerCompteurSynchro(&wg2)
	}
	wg2.Wait()
	duree2 := time.Since(start)

	fmt.Printf("Résultat : %d (attendu : %d)\n", compteur, attendu)
	fmt.Printf("Correct  : %v\n", compteur == attendu)
	fmt.Printf("Durée    : %v\n", duree2)

	// === Étape 3 : Avec sync/atomic ===
	fmt.Println("\n=== Étape 3 : Compteur synchronisé avec sync/atomic ===")
	compteurAtomic = 0
	var wg3 sync.WaitGroup

	start = time.Now()
	for i := 0; i < nbGoroutines; i++ {
		wg3.Add(1)
		go incrementerCompteurAtomic(&wg3)
	}
	wg3.Wait()
	duree3 := time.Since(start)

	fmt.Printf("Résultat : %d (attendu : %d)\n", compteurAtomic, attendu)
	fmt.Printf("Correct  : %v\n", compteurAtomic == int64(attendu))
	fmt.Printf("Durée    : %v\n", duree3)

	// === Comparaison ===
	fmt.Println("\n=== Comparaison des durées ===")
	fmt.Printf("Sans synchro : %v\n", duree1)
	fmt.Printf("Mutex        : %v\n", duree2)
	fmt.Printf("Atomic       : %v\n", duree3)
}
