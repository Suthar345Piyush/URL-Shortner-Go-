// URL shortner  application in golang
package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"Creation_date"`
}

//schema is something look like

/*
   d974843 --> {
		                        ID : "d974843"
														OriginalURL :  "random url"
														ShortURL : "d974843"
														CreationDate : time.Now()
	 }
*/

// variable for inline db to store the hash
var urlDB = make(map[string]URL)

//function to generate url

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()

	hasher.Write([]byte(OriginalURL)) // original url string to byte slice/array

	fmt.Println("hasher: ", hasher)

	data := hasher.Sum(nil)

	fmt.Println("hasher data: ", hasher)

	//here the hex code is encoded into the string format

	hash := hex.EncodeToString(data)

	fmt.Println("EncodeToString: ", hash)

	fmt.Println("Final string: ", hash[:8])

	return hash[:8]

}

//function to create short url

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

// function for getting the url

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("url not found")
	}
	return url, nil
}

func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Hello , world")
}

// http response construct by the response writer
// http request ,is recived to server or to send by the client

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}

	//this decoder reads the value from the json
	err := json.NewDecoder(r.Body).Decode(&data)
	//handling these error
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL := createURL(data.URL)

	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	//encoder writes to the w
	json.NewEncoder(w).Encode(response)
}

//redirect function , in this making a link , in which we have to just put the short hash (short url) then it will redirect to the particular webiste or page

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]

	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	fmt.Println("the url shortner project")

	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	//starting the http request on the server 8080
	fmt.Println("Starting the server on port 3000...")
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("Error on starting server: ", err)
	}
}
