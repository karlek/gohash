package wordlist

import "github.com/forsoki/gohash/str2hash"
import "github.com/mewkiz/pkg/bufioutil"
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
