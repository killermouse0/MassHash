package main

import (
	"fmt"
	"crypto/sha256"
	"os"
	"log"
	"encoding/hex"
	"io"
	"time"
	"bufio"
)

func computeHash(fileNames chan string, res chan string, quit chan bool) {
	for fname := range fileNames {
		f, err := os.Open(fname)
		hasher := sha256.New()
		if err != nil {
			log.Print(err)
		}
		_, err = io.Copy(hasher, f)
		if err != nil {
			log.Print(err)
		}
		f.Close()
		hash := hex.EncodeToString(hasher.Sum(nil))
		result := fmt.Sprintf("%s %s", hash, fname)
		res <- result
	}
	quit <- true
}

func printHashes(res chan string, quit chan bool) {
	for str := range res {
		fmt.Println(str)
	}
	quit <- true
}

func main() {
	const nthreads int = 20

	res := make(chan string)
	fileNames := make(chan string)
	quitComputeHash := make (chan bool)
	quit := make (chan bool)
	scanner := bufio.NewScanner(os.Stdin)

	/* Launching the hashing threads */
	for i:=0; i < nthreads; i++ {
		go computeHash(fileNames, res, quitComputeHash)
	}

	/* Launching the printing thread */
	go printHashes(res, quit)

	/* Feeding filenames to the hashing threads */
	total := 0
	for scanner.Scan() {
		fileNames <- scanner.Text()
		total++
	}
	close(fileNames)
	for i := 0; i < nthreads; i++ {
		<- quitComputeHash
	}
	close(res)
	for {
		select {
		case <- quit:
			return
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}
