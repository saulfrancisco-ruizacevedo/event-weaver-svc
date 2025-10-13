package models

type Event struct {
	Name          string      `crud:"pk,property:name"`
	Description   string      `crud:"property:description"`
	Deprecated    bool        `crud:"property:deprecated"`
	EffectiveFrom interface{} `crud:"property:effectiveFrom"`
	EffectiveTo   interface{} `crud:"property:effectiveTo"`
	Examples      string      `crud:"property:examples"`
	Schema        string      `crud:"property:schema"`
}
