package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"os"
)

const MIN = 1
const MAX = 100
var nombre_clients int = 0

func accepter_connection(connexions chan net.Conn,l net.Listener) {
	fmt.Println("Lancement de la Go Routine : Accepte les requêtes")
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		connexions <- c
}
}

func connect(c net.Conn, d chan net.Conn, Clients map[net.Conn]string, Message chan string) {
	for {
		message, _ := bufio.NewReader(c).ReadString('\n') /// Fonction blocante
		message = strings.Replace(message,"\n","",-1)
    parsed_args := strings.Split(message,"\t")
		if (len(message)>0){
			fmt.Println("Message recu : ",parsed_args)
		}
    switch parsed_args[0] {
    case "TCCHAT_REGISTER": // Ajouter aux autres qu'un utilisateur s'est connécté au serveur
						username := parsed_args [1]
						Clients[c] = "@" + username
						fmt.Printf("Un nouvel utilisateur à rejoint le chat ! \nNom D'utilisateur : %s\n",username)
						fmt.Printf("Nombre Actuel de chatters : %d\n",nombre_clients)
						send := "TCCHAT_USERIN\tVotre nom d'utilisateur est @" + username
						c.Write([]byte(send + "\n"))
						fmt.Printf("Message envoyé : %s \n",send)
						fmt.Printf("Destinataire : %s \n", c.RemoteAddr().String())
						break;
    case "TCCHAT_MESSAGE":
						message = Clients[c] + " : " + parsed_args[1]
						Message <- message
						fmt.Printf("Réception d'un message\n")
						fmt.Printf("Adresse IP : %s\n",c.RemoteAddr().String())
						fmt.Printf("Broadcast vers %d chatters\n",nombre_clients)
						break;
    case "TCCHAT_DISCONNECT":
							username := parsed_args[1]
							d <- c
							fmt.Printf("Déconnecté : %s\n", c.RemoteAddr().String())
							send := "TCCHAT_USEROUT\t" + "@" + username
							c.Write([]byte(send+ "\n"))
							fmt.Printf("Message envoyé [%s] : %s\n",c.RemoteAddr().String(),send)
		case "": // Evite de faire crasher le serveur quand un utilisateur se déconnecte
							break;
    default:
						send := "S : Message non valide : " + message + "\n"
						fmt.Print("Message envoyé : ",send)
            c.Write([]byte(send))
						fmt.Print("Destinataire :", c.RemoteAddr().String(), " " +message)
						message = ""
						break;
	   }
     }
     c.Close()
}
func Broadcast(identifiant net.Conn,user string, message string){
	fmt.Printf("Message envoyé [%s] : %s\n",identifiant.RemoteAddr().String(),message)

	identifiant.Write([]byte(message + "\n"))
}

func main() {
	fmt.Println("\t\t\t\tTCCHAT SERVEUR")
	connections_entrantes := make(chan net.Conn) // Channel pour les connexions entrantes
	deconnections_clients := make(chan net.Conn) // Channel pour les déconnexions du clients
	messages := make(chan string) // Dès qu'il y a un message dans le channel il faudra l'envoyer aux clients
	Clients := make(map[net.Conn]string) // Dictionnaire, Permet de connaitre la liste des utilisateurs -- Key Net.Conn -- Value String == Username
	c, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer c.Close() // Meme Si il y a une erreur j'arrête découter le port
	go accepter_connection(connections_entrantes,c) // connections_entrantes
	for {
		select { // L'instruction suivante devient non blocante ! Yes
    case requetes_client := <- connections_entrantes : // Je vide le channel quand il y a du contenu donc une connexion
				Clients[requetes_client] = "" // On définit un username éphemere
				nombre_clients += 1;
				fmt.Println("Un nouvel utilisateur à rejoint le serveur !")
				fmt.Println("Adresse IP :", requetes_client.RemoteAddr().String())
				requetes_client.Write([]byte("TCCHAT_WELCOME\tBONJOUR ET BIENVENUE SUR LE TCCHAT\n")) // J'envoie TCCHAT_Welcome
				fmt.Println("Connexion réussie ! ;)")
				go connect(requetes_client,deconnections_clients,Clients,messages)
		case deconnections := <- deconnections_clients :
				fmt.Println("Client Déconnecté:",Clients[deconnections])
				delete(Clients,deconnections)
				nombre_clients -= 1;
				fmt.Println("Nombre Actuel D'utilisateur :",nombre_clients)
		case reception_messages := <- messages:
				//Broadcast
				for i,j:=range Clients { // Je récupere chacune des clés de type net.Conn de tout les clients
					go Broadcast(i,j,reception_messages)
				}
		default:
    }
	}
}
