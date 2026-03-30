![or banner](assets/banner.png)

<h3 align="center">Zero-dependency Go utility for the classic "or-done" channel pattern.</h3>

<br>

## Installation

```bash
go get github.com/Pur1st2EpicONE/Or
```

<br>

## Usage

**Basic example**
```Go
package main

import (
	"fmt"
	"time"

	"github.com/Pur1st2EpicONE/Or"
)

func main() {

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			time.Sleep(after)
			close(c)
		}()
		return c
	}

	start := time.Now()
	<-or.Or(
		sig(2*time.Second),
		sig(300*time.Millisecond),
		sig(1*time.Second),
	)

	fmt.Printf("done after %v\n", time.Since(start).Round(100*time.Millisecond))

}
```

<br>

**Production example: Graceful Shutdown**

```Go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pur1st2EpicONE/Or"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)


	done := or.Or(
		ctx.Done(),                   
		sigChan,                   
		healthCheckFailed(),          
		time.After(24*time.Hour),     
	)

	log.Println("Service started...")

	go worker1(ctx)
	go worker2(ctx)
	go metricsCollector(ctx)

	<-done

	log.Println("Shutdown signal received, starting graceful shutdown...")
	cancel()

	time.Sleep(5 * time.Second)
	log.Println("Service stopped")

}

```