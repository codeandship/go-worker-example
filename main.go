package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

type job struct {
	ID        int
	Requested time.Time
	Finished  time.Time
}

type jobProcessor struct {
	queue chan job
	wg    *sync.WaitGroup
	stop  chan struct{}
}

func (r *jobProcessor) runWorker(id int) {
	defer r.wg.Done()
	fmt.Printf("starting worker-%d\n", id)
	var (
		jobCount int
		last     job
	)
workLoop:
	for {
		// listen for new jobs or stop
		select {
		case <-r.stop:
			// break workLoop to finish this goroutine
			break workLoop
		case job := <-r.queue:
			// simulate random processing time 1-500 ms
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+1))
			job.Finished = time.Now()
			last = job
			jobCount++
		}
	}
	fmt.Printf(
		"worker-%d stopped, %d jobs processed, last job: job.ID=%d processing_time=%s\n",
		id, jobCount, last.ID, last.Finished.Sub(last.Requested).String())
}

func main() {

	var (
		wg      = sync.WaitGroup{}
		workers = 10
		stop    = make(chan struct{})
	)

	rand.Seed(time.Now().Unix()) // seed random number generator

	jobQueue := make(chan job)

	// p shares jobQueue, wg, and stop with main
	p := jobProcessor{
		queue: jobQueue,
		wg:    &wg,
		stop:  stop,
	}

	// for every worker we add 1 to the wait group
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go p.runWorker(i) // start a worker asynchronically
	}

	// generate jobs and send them to jobQueue asynchronically
	wg.Add(1)
	go func() {
		defer wg.Done() // never forget to call wg.Done()
		i := 0
		for {
			job := job{
				ID:        i,
				Requested: time.Now(),
			}
			select {
			case <-stop:
				// leave this goroutine
				return
			case jobQueue <- job:
			}
			i++
		}
	}()

	// listen for interrupt signal (CTRL-C)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	wg.Add(1)
	go func() {
		defer wg.Done()  // never forget this ...
		<-sig            // blocks until we receive a signal
		close(stop)      // close stop channel to signal workers to quit
		fmt.Printf("\n") // print new line after CTRL-C
	}()

	// wait for all workers to be started before printing
	time.Sleep(time.Millisecond)
	fmt.Println("processing queue ... press CTRL-C to stop")

	// block until all goroutines have finished
	wg.Wait()
}
