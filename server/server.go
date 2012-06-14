package server

import "github.com/forsoki/gohash/attack"
import "github.com/forsoki/gohash/str2hash"

import "net/http"
import "fmt"
import "os"

//Starts HTTP server and hooks all events to handler
func HttpServer(portN int) (err error) {

	//Sets handler to httpRequestHandler
	http.HandleFunc("/", httpRequestHandler)

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
	strHash := r.Form.Get("hash")

	//If hash exist
	if len(strHash) == 0 {
		return
	}

	worder, err := attack.New(os.Getenv("GOPATH") + "/src/github.com/forsoki/gohash/a.txt")
	if err != nil {
		fmt.Println("New: ", err)
	}

	hash, err := str2hash.New(strHash)
	if err != nil {
		fmt.Println(err)
	}

	c := make(chan string)

	go worder.Check(hash, c)
	fmt.Println(<-c)
}
