package main

import (
	"fmt"
	"time"
)

//Spiegazione del Concetto di Worker Pool
//In questo esempio, una Worker Pool viene implementata utilizzando goroutine (i worker)
//che eseguono compiti (i lavori) estratti da un canale condiviso (lavori).
//I risultati dei compiti vengono inviati a un altro canale condiviso (risultati).

//Passaggi Chiave:

//Creazione dei Canali:
//lavori: per inviare compiti ai worker.
//risultati: per ricevere i risultati dai worker.

//Avvio dei Worker:
//Vengono create numWorker goroutine, ciascuna delle quali esegue la funzione worker.

//Invio dei Lavori:
//I lavori vengono inviati al canale lavori.

//Chiusura del Canale Lavori:
//Una volta inviati tutti i lavori, il canale lavori viene chiuso per indicare ai worker che non ci sono più lavori.

//Raccolta dei Risultati:
//I risultati vengono raccolti dal canale risultati e stampati.

// worker esegue i lavori ricevuti dal canale lavori e invia i risultati al canale risultati.
func worker(id int, lavori <-chan int, risultati chan<- int) {
	for lavoro := range lavori {
		fmt.Printf("Worker %d: inizio lavoro %d\n", id, lavoro)
		// Simula un lavoro lungo un secondo.
		time.Sleep(time.Second)
		fmt.Printf("Worker %d: fine lavoro %d\n", id, lavoro)
		// Invia il risultato del lavoro al canale risultati.
		risultati <- lavoro * 2
	}
}

func main() {
	const numLavori = 5
	const numWorker = 3

	// Crea i canali per i lavori e per i risultati.
	lavori := make(chan int, numLavori)
	risultati := make(chan int, numLavori)

	// Avvia i worker.
	for i := 1; i <= numWorker; i++ {
		go worker(i, lavori, risultati)
	}

	// Inserisce i lavori nel canale lavori.
	for j := 1; j <= numLavori; j++ {
		lavori <- j
	}
	close(lavori) // Chiude il canale lavori poiché non ci sono più lavori da inviare.

	// Raccoglie i risultati.
	for i := 1; i <= numLavori; i++ {
		result := <-risultati
		fmt.Printf("Risultato: %d\n", result)
	}
}
