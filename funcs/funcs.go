//Contains functions for gohash
package funcs

//Inner
import "github.com/forsoki/gohash/attacks"

//Output
import "fmt"
import "errors"

//File
import "os"

//Http server
import "net/http"

//Returns help string
func Help() string {
	return "Gohash is a hash cracker.\n\nUsage:\n\tgohash command [arguments]\n\nThe commands are:\n\n\t-c\tCommand-line input for hashes\n\t-s\tStart server on port 8080 to listen for input"
}

func Queue(hash, fileName string) (err error) {

	file, err := os.Create(fileName)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	file.Close()

	file, err = os.OpenFile(fileName, os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err 	= file.Seek(0, 2)
	if err != nil {
		return err
	}

	ret, err := file.WriteString(hash)
	if err != nil || ret != len(hash) {
		return err
	}

	return nil
}

//Checks if port number is valid
func IsValidPortN(portN int) (checkedPortN int, err error) {
	if portN < 0 || portN > 65535 {
		return 0, errors.New("Port is outside port range")
	}

	return portN, nil
}

//Starts HTTP server and hooks all events to handler
func ListenOnHttp(portN int) (err error) {

	//Sets handler to httpRequestHandler
	http.HandleFunc("/", httpRequestHandler)

	//Checks if port number is valid
	portN, err = IsValidPortN(portN)
	if err != nil {
		return err
	}

	//Start HTTP server
	err = http.ListenAndServe(":"+fmt.Sprint(portN), nil)
	if err != nil {
		return err
	}

	return nil
}

//All requests made to the HTTP server are processed by this function
func httpRequestHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}

	//Prints webpage
	fmt.Fprintf(w, "<html><form action=\"\" method=\"post\"><input type=\"text\" name=\"hash\"><input type=\"submit\"></form></html>")

	//Retrive value from form of hash
	hash := r.Form.Get("hash")

	//If hash exist
	if len(hash) == 0 {
		return
	}

	//Make a word-list attack on entered hash
	found, err := attacks.WordList(hash, os.Getenv("GOPATH")+"/src/github.com/forsoki/gohash/a.txt")
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	} else {
		fmt.Fprintf(w, "<br>%s = %s\n", hash, found)
	}
}
