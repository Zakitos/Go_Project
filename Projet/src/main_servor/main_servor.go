package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const MIN = 1
const MAX = 100

func connect(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String()) // On établit une adresse à l'utilisateur

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
            break

    default :

            c.Write([]byte("error" + "\n"))
	   }

     fmt.Print("Message Received from : ", c.RemoteAddr().String(), " "+message)

     }
      c.Close()
}

func main() {

	l, err := net.Listen("tcp4", ":8081") // J'écoute une co sur le port 8081
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close() // A la fin meme si il y a des erreur tu arrête d'écouter le port
	for {
		c, err := l.Accept() // Accepte la requête de l'utilisateur
		if err != nil {
			fmt.Println(err)
			return
		}
		go connect(c) // On fait une go routine pour chacune des requêtes
	}
}
