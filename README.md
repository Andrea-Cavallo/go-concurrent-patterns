# Gestione della Concorrenza in Go

## Introduzione

Go, spesso chiamato Golang, è un linguaggio di programmazione moderno creato da Google, progettato per semplificare lo sviluppo di software scalabile e ad alte prestazioni. Una delle caratteristiche più potenti di Go è il suo supporto integrato per la concorrenza, che consente di eseguire più operazioni simultaneamente in modo efficiente. In questo documento, esploreremo alcuni dei pattern più comuni per gestire la concorrenza in Go, fornendo esempi pratici per ciascuno.




# Pattern di Concorrenza in Go

Go offre potenti strumenti per la gestione della concorrenza, grazie alla sua capacità di gestire migliaia di goroutine in modo efficiente. Due dei pattern più comuni utilizzati in Go per sfruttare la concorrenza sono **FAN-OUT** e **FAN-IN**. In questo README, esploreremo questi due concetti con esempi pratici.

## FAN-OUT

Il pattern **FAN-OUT** si riferisce alla pratica di suddividere un carico di lavoro tra multiple goroutine. Questo approccio sfrutta il parallelismo e la concorrenza, permettendo a più unità di lavoro di essere eseguite simultaneamente. È particolarmente utile quando si ha un grande numero di compiti indipendenti che possono essere eseguiti in parallelo.

### Esempio di FAN-OUT

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Worker %d sta eseguendo il lavoro\n", id)
}

func main() {
    var wg sync.WaitGroup
    numWorkers := 5

    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait()
    fmt.Println("Tutti i worker hanno terminato il loro lavoro")
}
```
In questo esempio, creiamo 5 goroutine, ognuna delle quali esegue il proprio compito in parallelo. Utilizziamo un sync.WaitGroup per assicurarci che il programma principale attenda la terminazione di tutte le goroutine prima di concludere.


## FAN-IN
Il pattern FAN-IN è il concetto inverso rispetto a FAN-OUT. Si riferisce alla pratica di raccogliere i risultati da più goroutine in un unico canale. Questo pattern è utile per aggregare i risultati di compiti paralleli in un singolo punto di raccolta.

### Esempio di FAN-IN
```go
package main

import (
"fmt"
)

func worker(id int, jobs <-chan int, results chan<- int) {
for job := range jobs {
    fmt.Printf("Worker %d sta processando il job %d\n", id, job)
    results <- job * 2 // Simulazione di un lavoro
    }
}

func main() {
numWorkers := 3
numJobs := 5

    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= numJobs; a++ {
        result := <-results
        fmt.Printf("Risultato del job: %d\n", result)
    }
}
```

In questo esempio, abbiamo tre worker goroutine che leggono lavori da un canale jobs, processano ciascun lavoro, e inviano il risultato su un canale results. Il programma principale raccoglie tutti i risultati da results e li stampa.


FAN-OUT: Divide il carico di lavoro tra multiple goroutine, permettendo l'esecuzione parallela di più compiti.
FAN-IN: Raccoglie i risultati da multiple goroutine in un unico canale, aggregando i risultati in un singolo punto.

## WORKERS POOL


Worker Pool
Il pattern Worker Pool è una variante del pattern FAN-OUT che permette di gestire il carico di lavoro in modo più controllato ed efficiente. Invece di creare un numero arbitrario di goroutine, viene creato un pool fisso di worker che eseguono i compiti. Questo aiuta a limitare l'uso delle risorse e a mantenere un controllo migliore sulla concorrenza.

Esempio di Worker Pool

```go
package main

import (
"fmt"
"sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
    fmt.Printf("Worker %d sta processando il job %d\n", id, job)
    results <- job * 2 // Simulazione di un lavoro
    }
}

func main() {
const numWorkers = 3
const numJobs = 5

    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    var wg sync.WaitGroup

    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    wg.Wait()
    close(results)

    for result := range results {
        fmt.Printf("Risultato del job: %d\n", result)
    }
}
```

In questo esempio, creiamo un pool di tre worker che leggono compiti dal canale jobs e inviano i risultati al canale results. Utilizziamo un sync.WaitGroup per attendere che tutti i worker completino i loro compiti prima di chiudere il canale results e raccogliere i risultati.


## Canali di Errori
In Go, è possibile creare canali di errori per gestire gli errori provenienti dalle goroutine in modo sicuro. Questo permette di centralizzare la gestione degli errori e di evitare che vengano ignorati o gestiti in modo inefficace. Utilizzando i canali di errori, è possibile inviare eventuali errori al programma principale o a una funzione di gestione dedicata.

Esempio di Canali di Errori
```go
package main

import (
"fmt"
"sync"
)

func worker(id int, jobs <-chan int, results chan<- int, errors chan<- error, wg *sync.WaitGroup) {
defer wg.Done()
for job := range jobs {
    if job%2 == 0 { // Simulazione di un errore per i lavori pari
        errors <- fmt.Errorf("worker %d ha incontrato un errore con il job %d", id, job)
        continue
    }
    fmt.Printf("Worker %d sta processando il job %d\n", id, job)
    results <- job * 2
    }
}

func main() {
const numWorkers = 3
const numJobs = 5

    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    errors := make(chan error, numJobs)

    var wg sync.WaitGroup

    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, errors, &wg)
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    go func() {
        wg.Wait()
        close(results)
        close(errors)
    }()

    for result := range results {
        fmt.Printf("Risultato del job: %d\n", result)
    }

    for err := range errors {
        fmt.Printf("Errore: %s\n", err)
    }
}
```

In questo esempio, creiamo un canale errors per gestire gli errori provenienti dalle goroutine. Ogni worker invia un errore sul canale errors se incontra un problema durante l'elaborazione di un job. Il programma principale raccoglie e stampa tutti i risultati dai canali results ed errors.


## Gestione degli Errori con `sync/errgroup`

In Go, il pacchetto `sync/errgroup` semplifica la gestione delle goroutine concorrenti e degli errori. Utilizzando `errgroup`, se una delle goroutine fallisce, tutte le altre goroutine vengono annullate. Questo rende il codice più semplice e robusto quando si ha a che fare con più goroutine che devono essere coordinate.

### Esempio con `sync/errgroup`

Ecco un esempio di utilizzo di `sync/errgroup`:

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "golang.org/x/sync/errgroup"
)

func main() {
    g, ctx := errgroup.WithContext(context.Background())

    // Simuliamo tre richieste HTTP concorrenti
    urls := []string{
        "http://example.com",
        "http://example.net",
        "http://example.org",
    }

    for _, url := range urls {
        // Creiamo una goroutine per ciascun URL
        url := url // catturiamo la variabile url
        g.Go(func() error {
            // Simuliamo una richiesta HTTP con un timeout
            req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
            if err != nil {
                return err
            }

            client := &http.Client{Timeout: 2 * time.Second}
            resp, err := client.Do(req)
            if err != nil {
                return err
            }
            defer resp.Body.Close()

            fmt.Printf("Risposta da %s: %s\n", url, resp.Status)
            return nil
        })
    }

    // Attende la conclusione di tutte le goroutine o l'errore di una di esse
    if err := g.Wait(); err != nil {
        fmt.Printf("Errore: %v\n", err)
    } else {
        fmt.Println("Tutte le richieste completate con successo")
    }
}
```

Spiegazione
- errgroup.WithContext: Crea un gruppo di goroutine associato a un contesto. Se una goroutine restituisce un errore, il contesto viene annullato, e tutte le altre goroutine associate a quel contesto vengono terminate.
- g.Go: Aggiunge una nuova goroutine al gruppo. Se una goroutine restituisce un errore, tutte le altre goroutine vengono annullate.
- g.Wait: Attende la conclusione di tutte le goroutine nel gruppo. Se una goroutine restituisce un errore, g.Wait restituirà quell'errore.