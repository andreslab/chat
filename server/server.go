package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const addr = "127.0.0.1:3000"
const bufferSize = 256 //tamaño de almacenamiento
const endLine = 10     //salto de linea

var clients []net.Conn

func main() {
	clients = make([]net.Conn, 0)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("cant listen on " + addr)
		os.Exit(1)
	}
	for {
		conn, _ := listener.Accept()
		clients = append(clients, conn)
		//gorutina
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var data []byte                    //slide de bytes
	buffer := make([]byte, bufferSize) //usamos el mismo variable y no habria que crar mas instancias ne memoria

	//si el mensaje que recibe es de mas de 256 (bufferSize) se deberá hacer otro ciclo de lectura, por eso hay 2 for infinitos
	for {
		for {
			//SE EJECUTARÁ CADA 256 (bufferSize) CARACTERES

			n, err := conn.Read(buffer) //leer un mensaje de la conexion (conn) y lo escribe en "buffer"
			//n cuandtos bytes recibe
			if err != nil {
				if err == io.EOF {
					//EOF .. cuando se acaba el archivo
					break
				}
			}
			buffer = bytes.Trim(buffer[:n], "\x00") //array de bytes (buffer) y representacion de byte (\x00)
			//si tenemos un buffer de 4 \x00\x00\x00\x00 y recibimos e 2 (el \x10 es salto de linea endLine)  \x97\x98\x10\x00   trim filtra y quita los 00 (\x00)    \x97\x98\x10

			data = append(data, buffer...) //append mete un elemto en un slice (data .. es un slide de bytes)
			/*
				// buffer... es equivalente a

				for i := 0; i < len(buffer); i ++{
					append(data, buffer[i])
				}
			*/

			//detectamos que se acabo el mensaje,  lo indicamos con el endLine
			if data[len(data)-1] == endLine {
				fmt.Println("final del mensaje leido")
				break
			}

			//el max que se puede eviar por socket es 2^16 - 1  (65.535)
		}
		sendToOtherClients(conn, data) //envia mensaje a todos los conectados
		data = make([]byte, 0)         //reseteamos data
	}
}

//sender persona que envia el mensaje
func sendToOtherClients(sender net.Conn, data []byte) {
	fmt.Println("enviando mensaje")
	for i := 0; i < len(clients); i++ {
		fmt.Println("enviado ...")

		if clients[i] != sender {
			//si cliente no es igual al que lo envia
			fmt.Println("enviado a contacto")

			clients[i].Write(data)
		}
	}
}
