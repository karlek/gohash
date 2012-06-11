//Hashcracker with server capability and multiple hash algorithms.
package main

//Functions
import "github.com/forsoki/gohash/funcs"
import "github.com/forsoki/gohash/attacks"

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

var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")
var cFlag = flag.String("c", "", "Command line input for hash cracking")
var sFlag = flag.Bool("s", false, "Server mode for hash cracking")
var hFlag = flag.Bool("help", false, "Prints help message")
var qFlag = flag.String("q", "", "Queues the hash to be cracked")

func main() {

	flag.Parse()

	//Create profiling file if filename was entered
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		//Start profiling
		pprof.StartCPUProfile(f)

		//Stop profiling when program returns
		defer pprof.StopCPUProfile()
	}

	hash := ""

	//Check which flag is inputted
	switch {

	//If hash mode (-c), program also needs a hash as input
	case *cFlag != "":
		hash = *cFlag

	//If server mode (-s), listen on port 8080 for form input
	case *sFlag:
		err := funcs.ListenOnHttp(8080)
		if err != nil {
			log.Fatalln("ListenOnHttp: ", err)
		}

	//If a hash is queued (-q), add it to the queue file
	case *qFlag != "":
		err := funcs.Queue(*qFlag, os.Getenv("GOPATH")+"/src/github.com/forsoki/gohash/queue.txt")
		if err != nil {
			log.Fatalln("Queue: ", err)
		}
		fmt.Println("Hash has been queued!")
		os.Exit(1)

	//If help message is requested (-help)
	case *hFlag:
		fallthrough

	//If no supported flags were entered
	default:
		fmt.Println(funcs.Help())
		os.Exit(1)
	}

	/** Google attack
	* There's a high chance that the hash has already been cracked so we use google to find them!
	**/

	t0 := time.Now()

	results, err := attacks.Google(hash)
	if err != nil {
		log.Fatalln("Google: ", err)
	}

	t1 := time.Now()

	fmt.Printf("\n%d results on google\n", results)
	fmt.Printf("Googled hash in %v.\n", t1.Sub(t0))

	/** Bruteforce attack
	*
	**/

	t0 = time.Now()

	found, err := attacks.BruteForce(hash)
	if err != nil {
		log.Fatalln("BruteForce: ", err)
	}

	t1 = time.Now()

	fmt.Printf("\n%s = %s\n", hash, found)
	fmt.Printf("Hash found in %v.\n", t1.Sub(t0))

	/** Wordlist attack
	* Most people don't use random characters as their passwords, they use common words.
	* By hashing and comparing the words in the list, we can find the word if the hash is identical.
	**/

	t0 = time.Now()

	found, err = attacks.WordList(hash, os.Getenv("GOPATH")+"/src/github.com/forsoki/gohash/a.txt")
	if err != nil {
		log.Fatalln("WordList: ", err)
	}

	t1 = time.Now()

	fmt.Printf("\n%s = %s\n", hash, found)
	fmt.Printf("Hash found in %v.\n", t1.Sub(t0))
}
