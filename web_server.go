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

/* Pongo client con mayusculas porque esta en el tutorial asi,
lo que suecede es que puedo usar el struct Client desde otros modulos.
Si estuviese en minuscula le cambia el scope a este lugar. */    
type Client struct { 
    Remote_Addr string //ip del chabon, no me sirve aca porque ya tengo la referencia en clients
    URL_Path map[string]int // my key es el path, mi valor es el counter
}


var clients map[string]*Client

func main() {
    clients = make(map[string]*Client) //inicializo clients
	http.HandleFunc("/", handler)
    http.HandleFunc("/statistics", get_statistics) //para las estadisticas luego
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
    // chequeo que este cliente pueda hacer el fetch segun mis restricciones
    err := request_restrictions(writer, request)
    if err != nil {
        return // salgo
    }

	dominio := "https://api.mercadolibre.com"
	my_url := strings.Join([]string{dominio, request.URL.Path}, "")
	err = fetch(my_url, writer) // piso el err

    if err != nil {
        fmt.Fprint(writer, err.Error()) // S: resultado string, f: formato
        return // salgo
    }
	//request.Write(writer) // este metodo es parte de request, en contraposicion al otro
}

func request_restrictions(writer http.ResponseWriter, request *http.Request) error {
    client_data, ok := clients[request.RemoteAddr]
    if !ok { // si no lo encuentra, lo crea
        clients[request.RemoteAddr] = &Client{
            Remote_Addr: request.RemoteAddr, 
            URL_Path: make(map[string]int)}
        clients[request.RemoteAddr].URL_Path[request.URL.Path] = 1 //inicializa el path y pone el counter en 1
        fmt.Println("Client no existia.")
        return nil // si no existia no hay restricciones
    }
    
    my_count, ok := clients[request.RemoteAddr].URL_Path[request.URL.Path]
    if !ok { // si mi path no existe en mi client, lo crea
        clients[request.RemoteAddr].URL_Path[request.URL.Path] = 1
        fmt.Println("URL_Path no existia.")
        return nil // no existia, no hay restricciones
    }

    clients[request.RemoteAddr].URL_Path[request.URL.Path] ++
    
    fmt.Println(client_data)
    fmt.Println(my_count)

    return nil
}

func get_statistics(writer http.ResponseWriter, request *http.Request) {
    fmt.Fprintf(writer, "Coming soon!")
}