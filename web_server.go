package main

import (
	"fmt"
	"log"
	"net/http"
	//"io/ioutil"
	//"os"
	//"io"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))

}

// You need to specify what you will return after specifying the input parameters, this is not python.
func fetch(url string, writer http.ResponseWriter) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Errooooorrrrrr: ", err) // por que usar sprint en vez de cualquiera de los otros prints?
	}

	body_bytes := make([]byte, resp.ContentLength)

	_, err = resp.Body.Read(body_bytes) // _ ignora
	if err != nil {
		fmt.Println("While reading %s: %v", url, err)
	}
	//fmt.Println("Esto vendria a ser lo que recibimos: %d de la url %s", nbytes, url)
	writer.Write(body_bytes)
	resp.Body.Close() // "dont leak resources"

}

func handler(writer http.ResponseWriter, request *http.Request) {
	dominio := "https://api.mercadolibre.com"
	my_url := strings.Join([]string{dominio, request.URL.Path}, "")
	fetch(my_url, writer)
	//request.Write(writer)

}
