package models

type Team struct {
	Name  string `crud:"pk,property:name"`
	Lead  string `crud:"property:lead"`
	Email string `crud:"property:email"`
}
