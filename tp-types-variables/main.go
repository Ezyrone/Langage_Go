package main

import "fmt"

// Constantes
const (
	PI              = 3.14159
	NOM_APPLICATION = "Gestionnaire Go"
	ANNEE_LANCEMENT = 2023
)

const (
	Lundi = iota
	Mardi
	Mercredi
	Jeudi
	Vendredi
	Samedi
	Dimanche
)

func main() {
	// Exercice 1
	fmt.Println("Exercice 1")

	var nomUtilisateur string = "Jory"
	var ageUtilisateur int = 30
	var estConnecte bool = true
	var soldeCompte float64 = 1542.75

	fmt.Println("Nom :", nomUtilisateur)
	fmt.Println("Âge :", ageUtilisateur)
	fmt.Println("Connecté :", estConnecte)
	fmt.Println("Solde :", soldeCompte)

	// Exercice 2
	fmt.Println("Exercice 2")

	villeResidence := "Grenoble"
	codePostal := 38000
	tauxRemise := 15.5

	fmt.Printf("Ville : %v (type: %T)\n", villeResidence, villeResidence)
	fmt.Printf("Code postal : %v (type: %T)\n", codePostal, codePostal)
	fmt.Printf("Taux de remise : %v (type: %T)\n", tauxRemise, tauxRemise)

	// Exercice 3
	fmt.Println("Exercice 3")

	rayon := 10.5
	circonference := 2 * PI * rayon
	fmt.Printf("Circonférence d'un cercle de rayon %.2f : %.4f\n", rayon, circonference)

	fmt.Println("PI :", PI)
	fmt.Println("Application :", NOM_APPLICATION)
	fmt.Println("Année de lancement :", ANNEE_LANCEMENT)

	// Exercice 4
	fmt.Println("Exercice 4")

	ancienAge := ageUtilisateur
	ageUtilisateur = ageUtilisateur + 1
	fmt.Printf("Anniversaire ! Ancien âge : %d → Nouvel âge : %d\n", ancienAge, ageUtilisateur)

	var message string
	fmt.Printf("message (non initialisé) : \"%s\" (zero value de string = chaîne vide)\n", message)

	var compteur int
	fmt.Printf("compteur (non initialisé) : %d (zero value de int = 0)\n", compteur)

	// Bonus
	fmt.Println("Bonus")

	var a, b, c int = 1, 2, 3
	fmt.Println("Déclaration multiple :", a, b, c)

	fmt.Println("Lundi =", Lundi, "| Mercredi =", Mercredi, "| Dimanche =", Dimanche)

	var entier int = 42
	var decimal float64 = 3.14

	resultat := float64(entier) + decimal
	fmt.Printf("Conversion : %d (int) + %.2f (float64) = %.2f\n", entier, decimal, resultat)
}
