package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

const addr = "127.0.0.1:3000"
const bufferSize = 256 //tamaño de almacenamiento
const endLine = 10     //salto de linea

var nick string
var in *bufio.Reader //puntero analiza el string y te resuelve los caracteres especiales de catacter

func main() {
	in = bufio.NewReader(os.Stdin) // entrada estandar del sistema operativo

	for nick == "" {
		fmt.Printf("Ingresa tu nick: ")
		buf, _, _ := in.ReadLine()
		nick = string(buf)

	}

	var conn net.Conn
	var err error

	for {
		fmt.Printf("Conectando a %s...\n", addr)
		conn, err = net.Dial("tcp", addr) //conexion de lado del cliente
		if err == nil {
			//cuando se conecta sale del bucle
			break
		}
	}

	defer conn.Close()

	go reciveMessages(conn)
	handleConnection_(conn)
}

func handleConnection_(conn net.Conn) {

	for {
		buf, _, _ := in.ReadLine()
		if len(buf) > 0 {
			fmt.Println("mensaje enviado")
			//hay que enviar bytes no se puede enviar texto
			conn.Write(append([]byte(nick+"->"), append(buf, endLine)...))
			//"jaime ->" + string(buf) + "\n"

		}
	}
}

func reciveMessages(conn net.Conn) {
	fmt.Println("mensaje recibido")
	var data []byte
	buffer := make([]byte, bufferSize)

	for {
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			buffer = bytes.Trim(buffer[:n], "\x00")
			data = append(data, buffer...)
			if data[len(data)-1] == endLine {
				break
			}
		}
		fmt.Printf("%s", data)
		data = make([]byte, 0)
	}

}
