package mutation

//Wordlist mutations
import "strings"

var conversionMap = map[string]string{
	"a": "4",
	"e": "3",
	"l": "1",
	"o": "0",
	"t": "7",
}

//Exchanges characters to their "leet" correspondence:
//a -> 4, e -> 3, l -> 1, t -> 7
func Leet(input string) string {

	//Range through conversionMap to replace letters with it's leet correspondance
	for letter, leet := range conversionMap {

		//The last argument is the limit of replacements, -1 means no limit
		input = strings.Replace(input, letter, leet, -1)
	}

	return input
}
