package models

type Order struct {
	Number   string
	Products []*Product
}
