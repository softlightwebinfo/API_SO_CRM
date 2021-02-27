package model

type ORMJoin struct {
	joinType string
	table    string
	on       string
}
