// The du3 command computes disk usage of the files in a directory (uses goroutines for walkDir)
// go run main.go -v $GOPATH
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// sema is a counting sempahore for limiting
// concurrency in dirents (executed via buffered channel
// and remaining queue space). Must limit fs transactions
// as can have 1000s of goroutines running at once
var sema = make(chan struct{}, 20)

// channel for listening for cancelling program
var done = make(chan struct{})

// verbose logging with "-v" in CLI
var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {
	// Determine initial directories
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(done)                    // sends a message therefore satisfying the select { case <-done } below
		fmt.Println("Cancelling...")   // my own touch
	}()

	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}
	go func() {
		// https://en.wikipedia.org/wiki/Communicating_sequential_processes
		n.Wait()
		close(fileSizes)
	}()

	// Print the results periodically
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	// Print the results
	var nfiles, nbytes int64

loop: // loop label (break <label> will break both `for`` and `select`` statements)
	for {
		select {
		case <-done:
			// Drain fileSizes to allow existing goroutines to finish
			for range fileSizes {
				// Do nothing
			}
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)

	// How to test if `main` cleaned up after itself?
	// test with a `panic`. It will spit out a log of
	// all running goroutines. If `main` is the sole survivor
	// then cleanup was successful.
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	// performance implications of a `defer` statement?
	defer n.Done()
	if cancelled() {
		return // if cancelled, then new goroutines are no-ops
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			// within loop to reduce redundant goroutine ops
			// if cancelled during tenure
			if cancelled() {
				n.Add(1)
				subdir := filepath.Join(dir, entry.Name())
				go walkDir(subdir, n, fileSizes)
			}
		} else {
			fileSizes <- entry.Size()
		}
	}
}

func dirents(dir string) []os.FileInfo {
	// select to reduce bottleneck of
	// cancellation latency of the pgoram
	// from hundreds of milliseconds to tens
	select {
	case sema <- struct{}{}: // acquire token
	case <-done:
		return nil // cancelled
	}

	defer func() { <-sema }() // release token

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dul: %v\n", err)
		return nil
	}
	return entries
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
