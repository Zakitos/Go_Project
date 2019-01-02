package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"os"
	//"reflect"
)

const MIN = 1
const MAX = 100

func accepter_connection(connexions chan net.Conn,l net.Listener) {
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		connexions <- c
	}
}

func connect(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	for {

    message, _ := bufio.NewReader(c).ReadString('\n')
    parsed_args := strings.Split(message, "/")

    switch parsed_args[0] {

    case "TCCHAT_REGISTER":

            c.Write([]byte("Bonjour " + parsed_args[1] + " bienvenue dans le chat" ))

    case "TCCHAT_MESSAGE":

            c.Write([]byte(parsed_args[1]))


    case "TCCHAT_DISCONNECT":
            c.Write([]byte("L'utilisateur " + parsed_args[1] + " c'est déconnecté"))
            fmt.Printf("Déconnecté : %s\n", c.RemoteAddr().String())
            break

    default :

            c.Write([]byte("error" + "\n"))

	   }

     fmt.Print("Message Received from : ", c.RemoteAddr().String(), " " +message)

     }
     c.Close()
}

func main() {

	connections := make(chan net.Conn)
	//messages := make(chan string) // bdcst
	c, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer c.Close() // Meme Si il y a une erreur j'arrête découter le port
	go accepter_connection(connections,c)
	for {
			//go connect(c)
	}
}
