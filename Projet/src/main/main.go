package main
import "fmt"

type personnage struct {
  nom string
  age int
  adulte bool
}

func main (){
  rachid := &personnage{nom : "Rachid", age: 27}
  majeur(rachid)
  fmt.Println(rachid.adulte)
}

func majeur(a *personnage){
  a.adulte = a.age >= 21
}
