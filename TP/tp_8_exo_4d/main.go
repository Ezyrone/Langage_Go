package main

import (
	"fmt"
	"math/rand"
	"time"
)

func dataProducer(dataCh chan<- string) {
	mesures := []string{
		"Température: 25°C",
		"Humidité: 60%",
		"Pression: 1013 hPa",
		"Température: 27°C",
		"Humidité: 55%",
	}
	for {
		delai := time.Duration(1+rand.Intn(3)) * time.Second
		time.Sleep(delai)
		dataCh <- mesures[rand.Intn(len(mesures))]
	}
}

func alertProducer(alertCh chan<- string) {
	alertes := []string{
		"Niveau critique atteint!",
		"Seuil de température dépassé!",
		"Capteur hors ligne!",
	}
	for {
		delai := time.Duration(5+rand.Intn(6)) * time.Second
		time.Sleep(delai)
		alertCh <- alertes[rand.Intn(len(alertes))]
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	dataChannel := make(chan string)
	alertChannel := make(chan string)
	quitChannel := make(chan struct{})

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go dataProducer(dataChannel)
	go alertProducer(alertChannel)

	go func() {
		time.Sleep(15 * time.Second)
		close(quitChannel)
	}()

	fmt.Println("Système de surveillance démarré.")

	for {
		select {
		case msg := <-dataChannel:
			fmt.Printf("[MESURE] %s\n", msg)
		case msg := <-alertChannel:
			fmt.Printf("[ALERTE CRITIQUE] ⚠️  %s\n", msg)
		case <-ticker.C:
			fmt.Println("[STATUS] Vérification système...")
		case <-quitChannel:
			fmt.Println("Signal d'arrêt reçu. Arrêt du système.")
			return
		}
	}
}
