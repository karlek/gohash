package wordlist

import (
	"strconv"
	"container/list"

	"github.com/karlek/gohash/str2hash"
	"github.com/mewkiz/pkg/bufioutil"
)
type Wordlist struct {
	MutateFuncs []func(string) string
	Words       *list.List
}

func New(fileName string) (wl Wordlist, err error) {
	l := list.New()
	ws, err := bufioutil.LoadLines(fileName)
	if err != nil {
		return Wordlist{}, err
	}
	for _, w := range ws {
		l.PushBack(w)
	}
	wl.Words = l
	return wl, nil
}

func (worder *Wordlist) Salt(prefix, suffix string) {
	f := func(word string)string {
		return prefix + word + suffix
	}
	worder.SaltFunc(f)
}

func (worder *Wordlist) SaltFunc(f func(string)string) {
	for e := worder.Words.Front(); e != nil; e = e.Next() {
		e.Value = f(e.Value.(string))
	}
}

func (worder Wordlist) Check(hash *str2hash.Hash, c chan string) {
	for e := worder.Words.Front(); e != nil; e = e.Next() {
		if hash.Hash == hash.HashFunc(e.Value.(string)) {
			c <- "Hash found: " + strconv.Quote(e.Value.(string))
		}
	}
	c <- "No wordlist found: " + hash.Hash
}

func (worder *Wordlist) Mutate() {
	for _, mutateFunc := range worder.MutateFuncs {
		for e := worder.Words.Front(); e != nil; e = e.Next() {
			worder.Words.PushBack(mutateFunc(e.Value.(string)))
		}
	}
}
