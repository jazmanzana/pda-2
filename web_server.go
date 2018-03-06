package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    //"os"
    "io"
    "strings"
)

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))

}

func fetch(url string) string {
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Errooooorrrrrr: ", err) // por que usar sprint en vez de cualquiera de los otros prints?
        return "error"
    }
    fmt.Sprintf("This is resp.Body: %s", resp.Body)
    nbytes, err := io.Copy(ioutil.Discard, resp.Body) // utilidad?
    resp.Body.Close() // "dont leak resources"
    if err != nil {
        fmt.Println("While reading %s: %v", url, err)
        return "error"
    }
    fmt.Println("Esto vendria a ser lo que recibimos: %d de la url %s", nbytes, url)
    return "response"
}


func handler(w http.ResponseWriter, r *http.Request) {
    dominio := "https://api.mercadolibre.com"
    my_url := strings.Join([]string{dominio, r.URL.Path}, "")
    fmt.Println("Entraste a handler por la url ", my_url) // consola
    fmt.Fprintf(w, fetch(my_url))
}


