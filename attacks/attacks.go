package attacks

import "github.com/forsoki/gohash/str2hash"
import "github.com/forsoki/gohash/mutation"
import "encoding/hex"

//Error handling
import "errors"

//Google attack
import "net/http"
import "io/ioutil"

//Brute test
import "fmt"

//File
import "os"
import "bufio"
import "io"

//Mutations
import "strconv"
import "strings"

// WordList tries to find the hash in the word list and returns the found
// string on success.
func WordList(hash, wordListName string) (found string, err error) {

	//Validate hash by trying to decode it as hexadecimal
	_, err = hex.DecodeString(hash)
	if err != nil {
		return "", err
	}

	//Validate length of hash. 32 (MD5) 40 (SHA1)
	if len(hash) != 32 && len(hash) != 40 {
		return "", errors.New("Invalid hash - length mismatch")
	}

	//Use hashing function determined by length of hash
	hashFunc := str2hash.MD5
	if len(hash) == 40 {
		hashFunc = str2hash.SHA1
	}

	//Open word list
	file, err := os.Open(wordListName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//Make file into reader
	r := bufio.NewReader(file)

	//For each word in word list, hash it and compare it to the entered hash.
	//If they match, the password is the hashed word
	for {

		l, _, err := r.ReadLine()

		//Break if EOF to induce "Hash not found"
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		//[]byte to string
		word := string(l)

		//Compare hash with hashed word
		if hash == hashFunc(word) {
			return word, nil
		}

		//Compare hash with hashed titled word
		titled := strings.Title(word)
		if hash == hashFunc(titled) {
			return titled, nil
		}

		//Compare hash with hashed uppercase word
		upper := strings.ToUpper(word)
		if hash == hashFunc(upper) {
			return upper, nil
		}

		//Compare hash with hashed lowercase word
		lower := strings.ToLower(word)
		if hash == hashFunc(lower) {
			return lower, nil
		}

		//Compare hash with hashed leeted word
		leet := mutation.Leet(word)
		if hash == hashFunc(leet) {
			return leet, nil
		}

		//Compare hash with hashed number suffixed word
		for i := 0; i <= 9999; i++ {

			//Concatenate word with integer (0 - 9999)
			intString := word + strconv.Itoa(i)

			if hash == hashFunc(intString) {
				return intString, nil
			}
		}
	}
	return "", errors.New("Hash not found")
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

	//Find "About " where after that string the results in numerals are printed.
	//Increment offset by len to remove "About "
	html := body[strings.Index(string(body), "About ")+len("About "):]

	//No search term on google is longer than 15 characters
	for i := 0; i < 15; i++ {

		//If it's a number concatenation the number on to stringResults.
		//Else if it's a ",", continue, since google uses these as number separators
		if isNumeric(html[i]) {
			stringResults += string(html[i])
		} else if html[i] == 0x2c {
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

func isNumeric(input byte) bool {

	if input >= 0x30 && input <= 0x39 {
		return true
	}

	return false
}

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
