package main

import (
	"fmt"
	// "bufio"
	// "os"
)


func main() {
	fmt.Print("¿Qué operación desea realizar? 1: suma, 2: resta, 3: multiplicación, 4: división ")
	var option int
	_, err := fmt.Scanf("%d", &option) // que onda el %d este?
    fmt.Println(option)
    if err != nil {
    	fmt.Println("Hubo un error.")
    	fmt.Println(err)
    }
    if option == 1 {
		fmt.Print("Ingrese los dos operandos: ")
		var x int
		_, err := fmt.Scanf("%d", &x)
		//_, err := fmt.Scanf("%d", &y)

		if err != nil {
	    	fmt.Println("Hubo un error.")
	    	fmt.Println(err)
    	}

		suma(x, x)
    }
	fmt.Println("Puto el que escribe.")
	
}

func suma(x, y int) int {
    fmt.Print("El resultado es: ")
	fmt.Println(x+y)
	return x + y
}

func resta(x, y int) int {
	return x - y
}

func multiply(x, y int) int {
	return x * y
}

func divide(x, y int) int {
	return x / y
}

func defino_cosas(){
	var operaciones [4]string
	operaciones : ["suma", "resta", "multiplicación", "división"]
}