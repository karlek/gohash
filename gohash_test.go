//Test cases for word list attack with different type of hashes
package main

import "github.com/karlek/gohash/attack"
import "testing"

func TestMD5(t *testing.T) {

	m := map[string]string{
		"d41d8cd98f00b204e9800998ecf8427e": "",          //Empty string ""
		"d41d8cd98f00b204e9800998ecf8427":  "",          //Invalid MD5 (1 character too short)
		"8b1a9953c4611296a827abf8c47804d7": "Hello",     //Titled test: Hello
		"9e076f5885f5cc16a4b5aeb8de4adff5": "Not found", //Unfound test: Not found
		"e8bb0b2e10d6706a0ae1a8633a9feace": "asdf0",     //Number suffix test: asdf0
		"6a47b3f8f52318528ab6438078b28ad4": "asdf9999",  //Number suffix test: asdf9999
	}

	mError := map[string]string{
		"d41d8cd98f00b204e9800998ecf8427":  "encoding/hex: odd length hex string",
		"9e076f5885f5cc16a4b5aeb8de4adff5": "Hash not found",
	}

	for hash, out := range m {
		found, err := attacks.WordList(hash, os.Getenv("GOPATH")+"/src/github.com/karlek/gohash/a.txt")
		if err != nil && mError[hash] != err.Error() {
			t.Errorf("%s - failed with error: %s\n", hash, err.Error())
		}
		if out != found {
			t.Errorf("wordListAttack(%v) = %v, want %v\n\n", hash, found, out)
		}
	}
}

func TestSHA1(t *testing.T) {

	m := map[string]string{
		"da39a3ee5e6b4b0d3255bfef95601890afd80709": "",          //Empty string ""
		"da39a3ee5e6b4b0d3255bfef95601890afd8070":  "",          //Invalid SHA-1 (1 character too short)
		"f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0": "Hello",     //Titled test: Hello
		"475c848673a3f79fa778f01c2bd5a721d4c41707": "Not found", //Unfound test: Not found
		"e7587ca621f0819e68f5b740ebec0b7c5f292fac": "asdf0",     //Number suffix test: asdf0
		"ce743a784552c605fbd94774348f600cc10c8d2c": "asdf9999",  //Number suffix test: asdf9999
	}

	mError := map[string]string{
		"da39a3ee5e6b4b0d3255bfef95601890afd8070":  "encoding/hex: odd length hex string",
		"475c848673a3f79fa778f01c2bd5a721d4c41707": "Hash not found",
	}

	for hash, out := range m {
		found, err := attacks.WordList(hash, os.Getenv("GOPATH")+"/src/github.com/karlek/gohash/a.txt")
		if err != nil && mError[hash] != err.Error() {
			t.Errorf("%s - failed with error: %s\n", hash, err.Error())
		}
		if out != found {
			t.Errorf("wordListAttack(%v) = %v, want %v\n\n", hash, found, out)
		}
	}
}
