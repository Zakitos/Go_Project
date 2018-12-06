package main

import (
  "fmt"
  "net/http"
)
const SEND_MESSAGE = len("/MESSAGE/")
func handler (w http.ResponseWriter, r *http.Request){
  title := string(r.URL.Path[SEND_MESSAGE:])
  fmt.Println()
}

func main (){
  http.HandleFunc("/MESSAGE/",handler) // Si on a une requete sur le localhost:8080/ on fait handler
  http.ListenAndServe(":8080", nil) // On lance le serveur, on écoute en boucle
}
