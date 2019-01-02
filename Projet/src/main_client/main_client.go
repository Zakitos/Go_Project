package main

import "net"
import "fmt"
import "bufio"
import "os"

func envoyer_message(conn net.Conn){
  for {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Message à envoyer : ")
    text, _ := reader.ReadString('\n') // On lit jusqu'a \n -- Blocante
    // send to socket
    fmt.Fprintf(conn, text + "\n") // Envoie au socket
    fmt.Print("\n")
    // Ecoute du serveur
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Print("\nMessage from server: " + message)
  }
}
func ecouter_serveur(conn net.Conn){
  message, _ := bufio.NewReader(conn).ReadString('\n')
  switch (message){
    case "balababaz"
  }
  // BLABLA
}

func main() {
  // Connexion avec le serveur
  conn,err := net.Dial("tcp", "127.0.0.1:8081")
  if (err != nil){
  fmt.Println(err)
  os.Exit(3)
  }else {
    go envoyer_message(conn)
    go ecouter_serveur(conn)
  }
  for{
  }
}
