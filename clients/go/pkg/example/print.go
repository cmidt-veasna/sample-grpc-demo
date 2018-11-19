package example

import "fmt"

func (e *Element) print() {
	if e == nil {
		return
	}
	fmt.Println("Id:\t\t", e.Id)
	fmt.Println("Name:\t\t", e.Name)
	fmt.Println("Age:\t\t", e.Age)
	fmt.Println("Statue:\t\t", e.Status)
	fmt.Println("CreateAt:\t", e.CreatedAt)
	fmt.Println("UpdateAt:\t", e.CreatedAt)
}

func (es *Elements) print() {
	if es == nil {
		return
	}
	for i := range es.Elements {
		es.Elements[i].print()
		fmt.Println("-------")
	}
}
