package intro_ioutil
import (
  "fmt"
  "io/ioutil"
)

type Page struct {
  Title string
  Body []byte
}

func main (){
p1 := &Page{Title: "TestPage" , Body: []byte("This is a sample page")}
p1.save()
p2, _ := loadPage("TestPage")
fmt.Println(string(p2.Body))
}

// On enrengistre la page et le body dans un fichier.txt
func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename,p.Body,0600) // 0600 lecture ecriture
}
// Lit le contenu du fichier et le transmet a la seconde structure
func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if (err!=nil){
    return nil,err
  }
  return &Page{Title : title, Body : body},nil
}
