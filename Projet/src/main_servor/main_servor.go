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
var nbr_users int = 0
func accepter_connection(connexions chan net.Conn,l net.Listener) {
	fmt.Println("\t\t\t\tTCCHAT SERVEUR")
	for {
		fmt.Println("a")
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		nbr_users += 1;
		fmt.Println("Un nouvel utilisateur à rejoint le serveur !")
		fmt.Println("Adresse IP :", c.RemoteAddr().String())
		fmt.Println("Nombre Actuel D'utilisateur :",nbr_users)
		connexions <- c
	}
}

func connect(c net.Conn) {
		fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {

    message, _ := bufio.NewReader(c).ReadString('\n')
    parsed_args := strings.Split(message, "/")

    switch parsed_args[0] {

    case "TCCHAT_REGISTER": // Ajouter aux autres qu'un utilisateur s'est connécté au serveur

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
