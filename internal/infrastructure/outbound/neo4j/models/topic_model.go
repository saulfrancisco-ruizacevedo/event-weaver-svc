package models

type Topic struct {
	Name string `crud:"pk,property:name"`
	Type string `crud:"property:type"`
}
