//Hashcracker with server capability and multiple hash algorithms.
package main

//Functions
import "github.com/karlek/gohash/google"
import "github.com/karlek/gohash/wordlist"
import "github.com/karlek/gohash/str2hash"
import "github.com/karlek/gohash/server"

//Mutatation
import "strings"
import "github.com/karlek/gohash/mutation"

//Output
import "fmt"
import "log"

//File
import "os"

//Input
import "flag"

//Profiling
import "time"
import "runtime/pprof"

//Flags
var cpuprofile string
var cFlag string
var sFlag bool
var portN int
var hFlag bool

func init() {

	flag.StringVar(&cpuprofile, "cpuprofile", "", "Write cpu profile to file")
	flag.StringVar(&cFlag, "c", "", "Command line input for hash cracking")
	flag.BoolVar(&sFlag, "s", false, "Server mode for hash cracking")
	flag.IntVar(&portN, "port", 8080, "Port number for server mode")
	flag.BoolVar(&hFlag, "help", false, "Prints help message")

	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]... [hash]...\n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	//Hash inputted via command-line
	inputHash := ""

	//Check flags is inputted
	switch {

	//If command-line mode (-c), program reads in a hash
	case cFlag != "":
		inputHash = cFlag

	//If server mode (-s), start service for hash cracking
	case sFlag:
		server.HttpServer(portN)

	//Create profiling file if filename was entered
	case cpuprofile != "":
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		//Start profiling
		pprof.StartCPUProfile(f)

		//Stop profiling when program returns
		defer pprof.StopCPUProfile()

	//Print help message
	case hFlag:
		fallthrough

	default:
		usage()
		os.Exit(1)
	}

	hash, err := str2hash.New(inputHash)
	if err != nil {
		log.Fatalln("str2hash.New: ", err)
	}

	/** Google attack
	* There's a high chance that the hash has already been cracked so we use google to find them!
	**/

	t0 := time.Now()

	results, err := google.Google(hash.Hash)
	if err != nil {
		log.Fatalln("attack.Google: ", err)
	}

	t1 := time.Now()

	fmt.Printf("\n%d results on google\n", results)
	fmt.Printf("Googled hash in %v.\n", t1.Sub(t0))

	/** Wordlist attack
	 * Most people don't use random characters as their passwords, they use common words.
	 * By hashing and comparing the words in the list, we can find the word if the hash is identical.
	 **/

	t0 = time.Now()

	//Make wordlist from file
	worder, err := wordlist.New(os.Getenv("GOPATH") + "/src/github.com/karlek/gohash/a.txt")
	if err != nil {
		log.Fatalln("attack.New: ", err)
	}

	//Salt wordlist, in this case nothing happens since both strings are empty
	worder.Salt("", "")

	//Set mutate functions, these will affect the wordlist.Mutate() method 
	worder.MutateFuncs = []func(string) string{
		strings.Title,
		strings.ToUpper,
		strings.ToLower,
		mutation.Leet,
	}

	//Add all mutations of the words to the wordlist and keeping the original
	worder.Mutate()

	newWorder := worder
	newWorder.Words = Obsc(worder.Words)

	c := make(chan string)

	///This comment might be wrong, but this is how I understood go channels
	//The check functions runs through it's wordlist and searches for the correct string. If it fails it waits for the other goroutines. By using return, the goroutine which succeeds prematurely finishes the channel and prints the found string.
	//I know now that it doesn't prematurely finish the channel, but it somehow chooses the goroutine which returns.
	go worder.Check(hash, c)
	go newWorder.Check(hash, c)

	fmt.Println(<-c)

	t1 = time.Now()

	fmt.Printf("Hash lookup in %v.\n", t1.Sub(t0))
}

func Reverse(words []string) (newWords []string) {
	for i := len(words) - 1; i > 0; i-- {
		newWords = append(newWords, words[i])
	}

	return newWords
}

func Obsc(words []string) (newWords []string) {
	for i := 0; i < 1000000; i++ {
		newWords = append(newWords, "a")
	}

	return newWords
}
