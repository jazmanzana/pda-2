package main

import (
	"fmt"
	"log"
	"net/http"
	//"io/ioutil"
	//"os"
	//"io"
	"strings"
    "errors"
)

    
type client struct {
    RemoteAddr string //ip del chabon
    UrlPath string 
    Counter int 
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}

// You need to specify what you will return after specifying the input parameters, this is not python.
func fetch(url string, writer http.ResponseWriter) error {
	resp, err := http.Get(url)
	if err != nil {
        return err
	}

    if resp.StatusCode != 200 { // GET hace hasta 10 redirecciones por si solo
        writer.WriteHeader(resp.StatusCode)
        return errors.New(fmt.Sprintf("The request failed with %v.", resp.StatusCode))
    }
    // porque a los arrays les tengo que fijar el tamanio de entrada
	body_bytes := make([]byte, resp.ContentLength)

	_, err = resp.Body.Read(body_bytes) // _ ignora, err pisa al err anterior
	if err != nil {
        return err
	}
	writer.Write(body_bytes)
	resp.Body.Close() // "dont leak resources"

    return nil
}

func handler(writer http.ResponseWriter, request *http.Request) {
	dominio := "https://api.mercadolibre.com"
	my_url := strings.Join([]string{dominio, request.URL.Path}, "")
	err := fetch(my_url, writer)
    if err != nil {
        fmt.Fprint(writer, err.Error()) // S: resultado string, f: formato
        return // salgo
    }
	//request.Write(writer) // este metodo es parte de request, en contraposicion al otro
}
