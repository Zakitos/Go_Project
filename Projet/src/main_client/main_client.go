package main

import "net"
import "fmt"
import "bufio"
import "os"
import "strings"
import "time"

func envoyer_message(conn net.Conn){
  for {
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n') // On lit jusqu'a \n -- Blocante
    // send to socket
    fmt.Fprintf(conn, text + "\n") // Envoie au socket
    fmt.Print("\n")
  }
}
func ecouter_serveur(conn net.Conn){
  for {
    message, _ := bufio.NewReader(conn).ReadString('\n')
    message = strings.Replace(message,"\n","",-1)
    parsing := strings.Split(message,"\t")
    t := time.Now()
    switch parsing [0]{
      case "TCCHAT_WELCOME":
        receive := parsing [1]
        fmt.Printf("S[%d/%d/%d %dh%d] : %s\n",t.Day(),t.Month(),t.Year(),t.Hour(),t.Minute(),receive)
      case "TCCHAT_USEROUT":
        receive := parsing [1]
        fmt.Printf("S[%d/%d/%d %dh%d]: Déconnection : %s\n",t.Day(),t.Month(),t.Year(),t.Hour(),t.Minute(),receive)
      case "TCCHAT_USERIN" :
        receive := parsing[1]
        fmt.Printf("S[%d/%d/%d %dh%d] : %s vient de rejoindre le serveur ! Souhaitez-lui la bienvenue ;)\n",t.Day(),t.Month(),t.Year(),t.Hour(),t.Minute(),receive)
      case "TCCHAT_ERROR_ID" :
        receive := parsing [1]
        fmt.Printf("S[%d/%d/%d %dh%d] : %s\n",t.Day(),t.Month(),t.Year(),t.Hour(),t.Minute(),receive)
        os.Exit(3)
      case "TCCHAT_MESSAGE":
        receive := parsing [1]
        fmt.Printf("%s\n",receive)
      default: // On recoit un message d'un client
        os.Exit(3)
    }
  }
}

func main() {
  // Connexion avec le serveur
  conn,err := net.Dial("tcp", "127.0.0.1:8081")
  if (err != nil){
  fmt.Println(err)
  os.Exit(3)
  }else{
    go ecouter_serveur(conn)
    go envoyer_message(conn)
  }
  for{
  }
}
