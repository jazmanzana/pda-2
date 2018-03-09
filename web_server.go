package main

import (
	"fmt"
	"log"
	"net/http"
	"errors"
	"strings"
	"time"
    "sync"
)

/* Pongo client con mayusculas porque esta en el tutorial asi,
lo que suecede es que puedo usar el struct Client desde otros modulos.
Si estuviese en minuscula le cambia el scope a este lugar. */
type Restrict struct {
    URL_Count int
    URL_Time string
}

type Client struct {
	Remote_Addr string         //no me sirve aca porque ya tengo la referencia a la ip en clients -key-
	URL_Path    map[string]Restrict // my key es el path, mi valor es el {counter, date}
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
		return                          // salgo
	}
	//request.Write(writer) // este metodo es parte de request, en contraposicion al otro
}

func request_restrictions(writer http.ResponseWriter, request *http.Request) error { //usar mutex aca para que bloquee
    var mutex = &sync.Mutex{} // para que clients sea modificado de a un client por vez
    request_time := time.Now().Format(time.UnixDate)

	client_data, ok := clients[request.RemoteAddr]
	if !ok { // si no lo encuentra, lo crea
        mutex.Lock()
		clients[request.RemoteAddr] = &Client{ //hacer que el key no tenga el puerto
			Remote_Addr: request.RemoteAddr,
			URL_Path:    make(map[string]Restrict)}

		clients[request.RemoteAddr].URL_Path[request.URL.Path] = Restrict{1, request_time} //inicializa el path y pone el counter en 1
        mutex.Unlock()
		fmt.Println("Client no existia.")
		return nil // si no existia no hay restricciones
	}

	_, ok = clients[request.RemoteAddr].URL_Path[request.URL.Path]
	if !ok { // si mi path no existe en mi client, lo crea
        mutex.Lock()
		clients[request.RemoteAddr].URL_Path[request.URL.Path] = Restrict{1, request_time}
        mutex.Unlock()
		fmt.Println("URL_Path no existia.")
		return nil // no existia, no hay restricciones
	}

	fmt.Println(clients[request.RemoteAddr].URL_Path[request.URL.Path])

    // si deja pasar 60 segundos entre llamadas, se resetea el counter y se pisa la fecha
    if reset_counter(clients[request.RemoteAddr].URL_Path[request.URL.Path].URL_Time) {
        mutex.Lock()
        clients[request.RemoteAddr].URL_Path[request.URL.Path] = Restrict{0, request_time}
        mutex.Unlock()
    }

	if clients[request.RemoteAddr].URL_Path[request.URL.Path].URL_Count == 5 {
		// me gustaria que para la 6ta le ponga un timer mas grande antes de devolver el mensaje
		fmt.Fprintf(writer, "Muchas requests en poco tiempo para esta url, intente mas tarde.")
		return errors.New(fmt.Sprintf("Hola, soy un mensaje de error poco claro."))
	} else {
        modified_strict := Restrict{
            clients[request.RemoteAddr].URL_Path[request.URL.Path].URL_Count + 1, 
            clients[request.RemoteAddr].URL_Path[request.URL.Path].URL_Time}
        mutex.Lock()
		clients[request.RemoteAddr].URL_Path[request.URL.Path] = modified_strict
        mutex.Unlock()
	}
	fmt.Println(client_data)

	return nil
}

func reset_counter(struct_time string) bool { // si paso mas de t tiempo entre las fechas, true
    const date_format = "Mon Jan  2 15:04:05 -07 2006"
    first_time, _ := time.Parse(date_format, struct_time)
    fmt.Println("struct_time:", struct_time)
    fmt.Println("first_time:", first_time)
    if time.Since(first_time) < 60 { //si intentaste en menos de 60 segundos 
        fmt.Println("false")
        return false
    }
    fmt.Println("true")
    return true
}

func get_statistics(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Coming soon!")
}
