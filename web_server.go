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

/*
    ch := make(chan string)
    for _, url := range os.Args[1:]{
        go fetch(url, ch) // arranca la go routine -> fetch de las urls que mande por parametro
    }
    for range os.Args[1:]{
        fmt.Println(<-ch) // recibe lo que venga de ch 
        }*/

}

func fetch(url string, ch chan<- string){
    resp, err := http.Get(url)
    if err != nil {
        ch <- fmt.Sprint(err) // por que usar sprint en vez de cualquiera de los otros prints?
        return
    }
    fmt.Sprintf("This is resp.Body: %s", resp.Body)
    nbytes, err := io.Copy(ioutil.Discard, resp.Body)
    resp.Body.Close() // "dont leak resources"
    if err != nil {
        ch <- fmt.Sprintf("While reading %s: %v", url, err)
        return
    }
    ch <- fmt.Sprintf("Esto vendria a ser lo que recibimos: %d de la url %s", nbytes, url)
}


func handler(w http.ResponseWriter, r *http.Request) {
    ch := make(chan string)
    dominio := "https://api.mercadolibre.com"
    my_url := strings.Join([]string{dominio, r.URL.Path}, "")
    fmt.Println("Entraste a handler por la url %s", my_url)
    fetch(my_url, ch)
}


