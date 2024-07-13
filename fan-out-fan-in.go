package main

import (
	"fmt"
)

// La funzione calcola il prodotto di id e i, poi invia il risultato nel canale out.
func calcola(id int, out chan<- int) {
	for i := 0; i < 5; i++ {
		out <- id * i
	}
}

func main() {
	// Creazione di due canali, c1 e c2.
	c1 := make(chan int)
	c2 := make(chan int)

	// Avvio di due goroutine separate, ciascuna esegue la funzione calcola con diversi valori di id.
	// Esempio di fan-out.
	go calcola(2, c1)
	go calcola(3, c2)

	// Creazione di un terzo canale, c, per il fan-in.
	c := make(chan int)

	// Avvio di una terza goroutine per raccogliere i risultati dai canali c1 e c2 e inviarli al canale c.
	go raccogliRisultati(c1, c2, c)

	// Lettura e stampa dei primi 10 valori dal canale c.
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
}

// esempio Fan-in
// raccogliRisultati riceve i valori dai canali c1 e c2 e li invia al canale c.
// a funzione raccogliRisultati implementa il pattern fan-in.
// Riceve i valori dai canali c1 e c2 e li invia al canale c.
// Questo significa che i risultati calcolati da due goroutine separate
// (una che esegue calcola(2, c1) e l'altra che esegue calcola(3, c2))
// vengono combinati in un unico flusso di dati nel canale c
func raccogliRisultati(c1, c2, c chan int) {
	for {
		select {
		case v := <-c1:
			c <- v
		case v := <-c2:
			c <- v
		}
	}
}
