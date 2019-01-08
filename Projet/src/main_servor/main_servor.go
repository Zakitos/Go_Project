package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"os"
	"time"
	"strconv"
	"sync"
)
var mapLock = &sync.Mutex{}
var nombre_clients int = 0

func accepter_connection(connexions chan net.Conn,l net.Listener) {
	fmt.Println("Lancement de la Go Routine : Accepte les requêtes")
	for {
		c, err := l.Accept() // fonction blocante
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		connexions <- c
	}
}

func connect(c net.Conn, d chan net.Conn, Clients map[net.Conn]string, Message chan string) {
	for {
		message, err := bufio.NewReader(c).ReadString('\n') /// Fonction blocante
		if (err != nil){
			d <- c
			break;
		}
		message = strings.Replace(message,"\n","",-1)
    parsed_args := strings.Split(message,"\t")
		if (len(message) > 0){
			fmt.Println("Message recu : ",parsed_args)
		}
    switch parsed_args[0] {
    case "TCCHAT_REGISTER": // Ajouter aux autres qu'un utilisateur s'est connécté au serveur
						username := "@" + parsed_args [1]
						flag := 0 // On suppose que de base l'username chosis n'est pas présent dans le serveur
						for _,j:=range Clients { // Je récupere chacune des clés de type net.Conn de tout les clients
							if (username == j){flag =1} // Si l'username est présent, flag passe à 1
						}
						if(flag == 0){ // Si username non présent , on peux le logger
							mapLock.Lock()
							Clients[c] = username
							mapLock.Unlock()
							fmt.Printf("Un nouvel utilisateur à rejoint le chat ! \nNom D'utilisateur : %s\n",username)
							fmt.Printf("Nombre Actuel de chatters : %d\n",nombre_clients)
							// Il faut notifier tout les utilisateurs de l'arrivée d'un nouveau tchateur
							fmt.Printf("Broadcast vers %d chatters\n",nombre_clients)
							send := "TCCHAT_USERIN\t" + username
							Message <- send
						}else{ // Sinon on lui dit de retenter de se connecter avec un autre username
							send := "TCCHAT_ERROR_ID\t"+"Un Chatter utilise déjà ce nom d'utilisateur"
							c.Write([]byte(send + "\n"))
							fmt.Printf("Message envoyé : %s \n",send)
							flag = 0
						}
						break;
    case "TCCHAT_MESSAGE":
			mapLock.Lock()
			tempo := Clients[c]
			mapLock.Unlock()
			if (tempo != ""){
				t := time.Now()
				y := t.Year() // retourner un int --,doit être convertie en string
				mo := t.Month() // retourner une variable de type t.Month
				d := t.Day()
				h := t.Hour()
				m := t.Minute()
				mapLock.Lock()
				message = "TCCHAT_MESSAGE\t" + Clients[c] + " [" + strconv.Itoa(d) + "/"+ strconv.Itoa(int(mo)) + "/" + strconv.Itoa(y) + " " + strconv.Itoa(h) + "h" + strconv.Itoa(m) + "] : " + parsed_args[1]
				mapLock.Unlock()
				Message <- message // Push
				fmt.Printf("Réception d'un message\n")
				fmt.Printf("Adresse IP : %s\n",c.RemoteAddr().String())
				fmt.Printf("Broadcast vers %d chatters\n",nombre_clients)
			}else{
				send := "TCCHAT_USEROUT\t"+"Vous n'êtes pas encore loggé sur le serveur"
				c.Write([]byte(send + "\n"))
				fmt.Printf("Message envoyé : %s \n",send)
			}
			break;
    case "TCCHAT_DISCONNECT":
							username := "@"+parsed_args[1]
							mapLock.Lock()
							copy_username := Clients[c]
							mapLock.Unlock()
							if (username == copy_username){
								fmt.Printf("Déconnecté : %s\n", c.RemoteAddr().String())
								// Il faut notifier tout les utilisateurs de la déconnexion d'un chatter
								send := "TCCHAT_USEROUT\t" + username + " à quitter le serveur"
								Message <- send
								fmt.Printf("Broadcast vers %d chatters\n",nombre_clients)
								c.Write([]byte("a" + "\n")) // on force la deconnexion du client
							}else {
								send := "TCCHAT_ERROR_ID\t"+"Il ne s'agit pas de votre nom d'utilisateur"
								c.Write([]byte(send + "\n"))
								fmt.Printf("Message envoyé : %s \n",send)
							}
							break;

		case "": // Evite de faire crasher le serveur quand un utilisateur se déconnecte
							fmt.Println("coucou")
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
}
func Broadcast(identifiant net.Conn, message string){
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
				mapLock.Lock()
				Clients[requetes_client] = "" // On définit un username éphemere
				mapLock.Unlock()
				nombre_clients += 1;
				fmt.Println("Un nouvel utilisateur à rejoint le serveur !")
				fmt.Println("Adresse IP :", requetes_client.RemoteAddr().String())
				requetes_client.Write([]byte("TCCHAT_WELCOME\tBONJOUR ET BIENVENUE SUR LE TCCHAT\n")) // J'envoie TCCHAT_Welcome
				fmt.Println("Connexion réussie ! ;)")
				go connect(requetes_client,deconnections_clients,Clients,messages)
		case deconnections := <- deconnections_clients :
				fmt.Println("Client Déconnecté:",Clients[deconnections])
				mapLock.Lock()
				delete(Clients,deconnections) // Nom_map ; key
				mapLock.Unlock()
				nombre_clients -= 1;
				fmt.Println("Nombre Actuel D'utilisateur :",nombre_clients)
		case reception_messages := <- messages:
				//Broadcast
				mapLock.Lock()
				for i,j:=range Clients { // Je récupere chacune des clés de type net.Conn de tout les clients
					if (j!= ""){
						Broadcast(i,reception_messages)
					}
				}
				mapLock.Unlock()
		default:
    }
	}
}
