package models

type Domain struct {
	Name        string `crud:"pk,property:name"`
	Description string `crud:"property:description"`
}
