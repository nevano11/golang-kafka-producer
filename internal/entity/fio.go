package entity

import "fmt"

type Fio struct {
	Surname   string `json:"surname"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (f Fio) String() string {
	return fmt.Sprintf("Fio:{Surname:%s, FirstName:%s, LastName:%s}", f.Surname, f.FirstName, f.LastName)
}
