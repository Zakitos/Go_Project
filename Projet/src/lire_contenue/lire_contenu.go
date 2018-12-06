package lire_contenue

import (
  "fmt"
  "net/http"
  "io/ioutil"
)

type Page struct {
  Title string
  Body []byte
}
const lenPath = len("/carl/")

// Lit le contenu du fichier et le transmet a la seconde structure
func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if (err!=nil){
    return nil,err
  }
  return &Page{Title : title, Body : body},nil
}

func handler (w http.ResponseWriter, r *http.Request){
  title := r.URL.Path[lenPath:]
  p,_ := loadPage(title)
  fmt.Fprintf(w, "<h1>%s</h1> <div> %s</div>", p.Title, p.Body)
}

func main (){
  http.HandleFunc("/carl/",handler) // On traite toutes les a partir du /view
  http.ListenAndServe(":8080", nil) // On lance le serveur, on écoute en boucle
}
