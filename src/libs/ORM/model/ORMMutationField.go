package model

type ORMMutationField struct {
	Into     ORMInto
	IsCreate bool
	Set      []ORMSet
	IsReturn bool
	Where    []ORMWhere
}
