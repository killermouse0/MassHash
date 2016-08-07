package main

import (
	"fmt"
	"crypto/sha256"
	"os"
	"log"
	"encoding/hex"
	"io"
	"bufio"
	"sync"
)

var wg sync.WaitGroup
var quit chan bool
var res chan string

func computeHash(fileNames chan string) {
	for fname := range fileNames {
		f, err := os.Open(fname)
		hasher := sha256.New()
		if err != nil {
			log.Print(err)
		} else {
			defer f.Close()
			_, err = io.Copy(hasher, f)
			if err != nil {
				log.Print(err)
			} else {
				hash := hex.EncodeToString(hasher.Sum(nil))
				result := fmt.Sprintf("%s %s", hash, fname)
				res <- result
			}
		}
	}
	wg.Done()
}

func printHashes() {
	for str := range res {
		fmt.Println(str)
	}
	quit <- true
}

func main() {
	const nthreads int = 20

	res = make(chan string)
	fileNames := make(chan string)
	quit = make (chan bool)
	scanner := bufio.NewScanner(os.Stdin)

	/* Launching the hashing threads */
	wg.Add(nthreads)
	for i := 0; i < nthreads; i++ {
		go computeHash(fileNames)
	}

	/* Launching the printing thread */
	go printHashes()

	/* Feeding filenames to the hashing threads */
	for scanner.Scan() {
		fileNames <- scanner.Text()
	}
	close(fileNames) // No more files, we close the channel

	wg.Wait() // Waiting for the hashing threads to finish
	close(res)	// They won't be outputting anything anymore
	<- quit
}
