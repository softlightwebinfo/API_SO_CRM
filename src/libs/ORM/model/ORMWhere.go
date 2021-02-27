package model

type ORMWhere struct {
	Column    string
	Type      string
	Value     interface{}
	ValueX    bool
	Condition string
}
