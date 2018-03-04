package main

import (
    "fmt"
    "log"
    "net/http"
)

var count int

func main() {
    http.HandleFunc("/count", counter)
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) { // por que pasa el puntero?
    count++
    fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
    fmt.Fprint(w, "EEeeeEee")
}

func counter(w http.ResponseWriter, r *http.Request){
    fmt.Fprint(w, "bUUUU")
    fmt.Fprintf(w, "Count %d\n", count)
}