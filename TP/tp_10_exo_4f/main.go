package main

import (
	"context"
	"fmt"
	"time"
)

func effectuerOperationLongue(ctx context.Context, id string) error {
	fmt.Printf("[%s] Début de l'opération...\n", id)
	for i := 1; i <= 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] Opération annulée à l'étape %d : %v\n", id, i, ctx.Err())
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
			fmt.Printf("[%s] Traitement étape %d/5...\n", id, i)
		}
	}
	fmt.Printf("[%s] Opération terminée avec succès.\n", id)
	return nil
}

func main() {
	fmt.Println("Démarrage du programme principal.")

	fmt.Println("\n=== Scénario 1 : Timeout court (2s) ===")

	ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel1()

	resultChan1 := make(chan error, 1)

	go func() {
		resultChan1 <- effectuerOperationLongue(ctx1, "Tâche 1")
	}()

	select {
	case err := <-resultChan1:
		if err != nil {
			fmt.Printf("Main: L'opération s'est terminée avec une erreur : %v\n", err)
		} else {
			fmt.Println("Main: L'opération s'est terminée avec succès avant le timeout.")
		}
	case <-ctx1.Done():
		fmt.Printf("Main: Timeout atteint : %v\n", ctx1.Err())
	}

	<-resultChan1

	fmt.Println("\n=== Scénario 2 : Timeout long (3s) ===")

	ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel2()

	resultChan2 := make(chan error, 1)

	go func() {
		resultChan2 <- effectuerOperationLongue(ctx2, "Tâche 2")
	}()

	select {
	case err := <-resultChan2:
		if err != nil {
			fmt.Printf("Main: L'opération s'est terminée avec une erreur : %v\n", err)
		} else {
			fmt.Println("Main: L'opération s'est terminée avec succès avant le timeout.")
		}
	case <-ctx2.Done():
		fmt.Printf("Main: Timeout atteint : %v\n", ctx2.Err())
	}

	fmt.Println("\n=== Scénario 3 : Annulation manuelle (après 1s) ===")

	ctx3, cancel3 := context.WithCancel(context.Background())

	resultChan3 := make(chan error, 1)

	go func() {
		resultChan3 <- effectuerOperationLongue(ctx3, "Tâche 3")
	}()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Main: Annulation manuelle envoyée.")
		cancel3()
	}()

	err := <-resultChan3
	if err != nil {
		fmt.Printf("Main: L'opération a été annulée : %v\n", err)
	} else {
		fmt.Println("Main: L'opération s'est terminée avec succès.")
	}

	fmt.Println("\nFin du programme principal.")
}
