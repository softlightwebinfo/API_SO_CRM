package model

type ORMField struct {
	from    ORMFrom
	where   []ORMWhere
	selects []ORMSelect
	join    []ORMJoin
}
