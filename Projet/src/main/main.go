package main

import (
  "fmt"
  "net/http"
  "strings"
)
const SEND_MESSAGE = len("/MESSAGE/")
func handler_message (w http.ResponseWriter, r *http.Request){
  message_received := r.URL.Path[SEND_MESSAGE:]
  fmt.Println(message_received)
  message := strings.Split(message_received," ")
  fmt.Println(message)
}

func main (){
  http.HandleFunc("/MESSAGE/",handler_message) // Si on a une requete sur le localhost:8080/ on fait handler
  http.ListenAndServe(":8080",nil) // On lance le serveur, on écoute en boucle
}
