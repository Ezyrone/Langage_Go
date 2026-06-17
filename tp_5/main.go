package main

import (
	"fmt"
)

// --- Partie 1 : Définition et Implémentation Implicite d'Interfaces ---

// Exercice 1.1 : Interface Notifier
type Notifier interface {
	Send(message string) error
}

// Exercice 1.2 : EmailNotifier
type EmailNotifier struct {
	Recipient string
	Sender    string
}

func (e EmailNotifier) Send(message string) error {
	fmt.Printf("[EMAIL] De %s à %s : %s\n", e.Sender, e.Recipient, message)
	return nil
}

// Exercice 1.3 : SMSNotifier
type SMSNotifier struct {
	PhoneNumber string
}

func (s SMSNotifier) Send(message string) error {
	fmt.Printf("[SMS] Envoi à %s : %s\n", s.PhoneNumber, message)
	return nil
}

// Exercice 1.4 : ConsoleNotifier
type ConsoleNotifier struct{}

func (c ConsoleNotifier) Send(message string) error {
	fmt.Printf("[CONSOLE] Message : %s\n", message)
	return nil
}

// --- Partie 2 : L'Interface Vide ---

// Exercice 2.1 : processData avec interface{}
func processData(data interface{}) {
	switch v := data.(type) {
	case int:
		fmt.Printf("Donnée de type entier : %d\n", v)
	case string:
		fmt.Printf("Donnée de type chaîne : %s\n", v)
	case bool:
		fmt.Printf("Donnée de type booléen : %t\n", v)
	case *EmailNotifier:
		fmt.Printf("Donnée de type EmailNotifier pour %s\n", v.Recipient)
	default:
		fmt.Printf("Type de donnée inconnu : %T\n", v)
	}
}

// --- Partie 3 : Intégration et Réflexion ---

// Exercice 3.1 : Smart Notifier
type User struct {
	Name  string
	Email string
	Phone string
}

func sendSmartNotification(data interface{}, message string) error {
	switch v := data.(type) {
	case User:
		if v.Email != "" {
			notifier := EmailNotifier{Recipient: v.Email, Sender: "system@notifications.com"}
			return notifier.Send(message)
		}
		if v.Phone != "" {
			notifier := SMSNotifier{PhoneNumber: v.Phone}
			return notifier.Send(message)
		}
		notifier := ConsoleNotifier{}
		return notifier.Send(fmt.Sprintf("Aucune méthode de contact pour %s", v.Name))
	case string:
		notifier := ConsoleNotifier{}
		return notifier.Send("Message générique : " + message)
	default:
		fmt.Printf("Type non supporté : %T\n", v)
	}
	return nil
}

func main() {
	// === Partie 1 : Exercice 1.5 - Utilisation Polymorphique ===
	fmt.Println("=== Partie 1 : Utilisation Polymorphique ===")

	notifiers := []Notifier{
		EmailNotifier{Recipient: "alice@example.com", Sender: "bob@example.com"},
		SMSNotifier{PhoneNumber: "+33 6 12 34 56 78"},
		ConsoleNotifier{},
	}

	for _, n := range notifiers {
		n.Send("Votre commande a été expédiée !")
	}

	// === Partie 2 : Exercice 2.2 - Appel de processData ===
	fmt.Println("\n=== Partie 2 : processData ===")

	processData(42)
	processData("Bonjour le monde")
	processData(true)
	processData(&EmailNotifier{Recipient: "alice@example.com", Sender: "system@example.com"})
	processData([]int{1, 2, 3})
	processData(3.14)

	// === Partie 3 : Exercice 3.1 - Smart Notifier ===
	fmt.Println("\n=== Partie 3 : Smart Notifier ===")

	userWithAll := User{Name: "Alice", Email: "alice@example.com", Phone: "+33 6 12 34 56 78"}
	userPhoneOnly := User{Name: "Bob", Email: "", Phone: "+33 6 98 76 54 32"}
	userNoContact := User{Name: "Charlie", Email: "", Phone: ""}

	sendSmartNotification(userWithAll, "Bienvenue sur la plateforme !")
	sendSmartNotification(userPhoneOnly, "Votre code de vérification est 1234")
	sendSmartNotification(userNoContact, "Important : mise à jour requise")
	sendSmartNotification("broadcast", "Maintenance prévue ce soir")
	sendSmartNotification(123, "Ce message ne sera pas envoyé")

	// === Exercice 3.2 : Réflexion ===
	fmt.Println("\n=== Exercice 3.2 : Questions de Réflexion ===")
	fmt.Println(`
1. L'implémentation implicite des interfaces en Go permet un découplage total :
   un type satisfait une interface simplement en implémentant ses méthodes,
   sans import ni déclaration "implements". Cela favorise la composition
   et permet de définir des interfaces côté consommateur (pas côté producteur).

2. L'interface vide (interface{}) est utile pour les fonctions génériques
   (ex: fmt.Println) ou les collections hétérogènes. Ses inconvénients :
   perte de la sécurité de typage à la compilation, nécessité d'assertions
   de type à l'exécution, et code moins lisible. Depuis Go 1.18, les
   génériques sont souvent préférables.

3. Les interfaces permettent de substituer des implémentations (ex: mock
   pour les tests), de découpler les modules (dépendance sur un contrat,
   pas une implémentation concrète), et de composer des comportements
   via de petites interfaces (io.Reader, io.Writer, etc.).`)
}
