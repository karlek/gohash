//Mutate word to find common passwords
package mutation

//Leet
import "strings"

//NumberSuffix
import "strconv"

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

func NumberSuffix(input string) []string {

	numberList := []string{}

	//Compare hash with hashed number suffixed word
	for i := 0; i <= 9999; i++ {

		//Concatenate word with integer (0 - 9999)
		numberList = append(numberList, input+strconv.Itoa(i))
	}

	return numberList
}
