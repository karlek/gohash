package attack

import "github.com/forsoki/gohash/str2hash"
import "github.com/mewkiz/pkg/bufioutil"

//Google attack
import "net/http"
import "io/ioutil"
import "unicode"
import "unicode/utf8"

import "strings"
import "strconv"

import "runtime"

type Wordlist struct {
	Words       []string
	MutateFuncs []func(string) string
}

func New(fileName string) (words Wordlist, err error) {

	words.Words, err = bufioutil.ReadLines(fileName)
	if err != nil {
		return words, err
	}
	return words, nil
}

func (worder *Wordlist) Salt(prefix, suffix string) {

	for key, word := range worder.Words {
		worder.Words[key] = prefix + word + suffix
	}
}

func (worder Wordlist) Check(hash *str2hash.Hash, c chan string) {

	for _, word := range worder.Words {
		if hash.Hash == hash.HashFunc(word) {
			c <- "Hash found: " + strconv.Quote(word)
			return
		}
	}
	runtime.Gosched()
	c <- "No wordlist found: " + hash.Hash
}

func (worder *Wordlist) Mutate() {

	newWords := []string{}

	for _, mutateFunc := range worder.MutateFuncs {
		for _, word := range worder.Words {
			newWords = append(newWords, mutateFunc(word))
		}
	}

	worder.Words = newWords
}

func Google(hash string) (results int, err error) {

	//Get response from google
	//hl is language variable, nfpr is automatic no-redirect for more believable search, q is query
	resp, err := http.Get("https://encrypted.google.com/search?hl=en&nfpr=1&q=" + hash)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	//Make resp.Body into []byte 
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	//String to be converted into integer
	stringResults := ""
	searchString := "About "
	//Find "About " where after that string the results in numerals are printed.
	//Increment offset by len to remove "About "
	html := body[strings.Index(string(body), searchString)+len(searchString):]

	//No search term on google is longer than 15 characters
	for i := 0; i < 15; i++ {

		//If it's a number concatenation the number on to stringResults.
		//Else if it's a ",", continue, since google uses these as number separators
		rune, _ := utf8.DecodeRune([]byte{html[i]})

		if unicode.IsDigit(rune) {
			stringResults += string(html[i])
		} else if rune == ',' {
			continue
		} else {
			break
		}
	}

	//If some numbers where added to the string
	if len(stringResults) > 0 {

		//Convert stringResults to int
		nResults, err := strconv.Atoi(stringResults)
		if err != nil {
			return 0, err
		}

		return nResults, nil
	}

	return 0, nil
}

/*
func BruteForce(hash string) (found string, err error) {

   //Start with a-z
   //when z, reset to a and add z
   //when za, reset to a and increment to zb
   //when zz, reset to a's and add a -> aaa

   hashFunc := str2hash.MD5
   if len(hash) == 40 {
      hashFunc = str2hash.SHA1
   }

   s := []string{}

   for char := 97; char <= 122; char++ {
      s[0] = string(char)
   }
   fmt.Println(s)
   fmt.Println(hashFunc(string(s)))

   return "", errors.New("Hash not found")
}
*/
