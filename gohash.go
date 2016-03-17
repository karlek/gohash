//Hashcracker with server capability and multiple hash algorithms.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/karlek/gohash/mutation"
	"github.com/karlek/gohash/str2hash"
	"github.com/karlek/gohash/wordlist"
)

var (
	mutate bool
	help   bool
	path   string
)

func init() {
	flag.StringVar(&path, "p", os.Getenv("GOPATH")+"/src/github.com/karlek/gohash/a.txt", "wordlist path")
	flag.BoolVar(&mutate, "m", false, "Add mutatations (lower, upper, title and leet)")

	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]... [hash]...\n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	for _, s := range flag.Args() {
		hash, err := str2hash.New(s)
		if err != nil {
			log.Fatalln("str2hash.New:", err)
		}
		find(hash)
	}
}

func find(hash *str2hash.Hash) {

	/** Google attack
	* There's a high chance that the hash has already been cracked so we use google to find them!
	**/

	t0 := time.Now()

	// results, err := google.Google(hash.Hash)
	// if err != nil {
	// 	log.Fatalln("attack.Google:", err)
	// }

	t1 := time.Now()

	// fmt.Printf("\n%d results on google\n", results)
	fmt.Printf("Googled hash in %v.\n", t1.Sub(t0))

	/** Wordlist attack
	 * Most people don't use random characters as their passwords, they use common words.
	 * By hashing and comparing the words in the list, we can find the word if the hash is identical.
	 **/

	t0 = time.Now()

	//Make wordlist from file
	worder, err := wordlist.New(path)
	if err != nil {
		log.Fatalln("attack.New:", err)
	}

	//Set mutate functions, these will affect the wordlist.Mutate() method
	worder.MutateFuncs = []func(string) string{
		strings.Title,
		strings.ToUpper,
		strings.ToLower,
		mutation.Leet,
	}

	//Add all mutations of the words to the wordlist and keeping the original
	if mutate {
		go worder.Mutate()
	}

	//Salt wordlist, in this case nothing happens since both strings are empty
	worder.Salt("", "")

	c := make(chan string)

	///This comment might be wrong, but this is how I understood go channels
	//The check functions runs through it's wordlist and searches for the correct string. If it fails it waits for the other goroutines. By using return, the goroutine which succeeds prematurely finishes the channel and prints the found string.
	//I know now that it doesn't prematurely finish the channel, but it somehow chooses the goroutine which returns.
	go worder.Check(hash, c)

	fmt.Println(<-c)

	t1 = time.Now()

	fmt.Printf("Hash lookup in %v.\n", t1.Sub(t0))
}
