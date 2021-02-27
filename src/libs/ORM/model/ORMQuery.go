package model

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"so-crm/src/libs/ORM/lib"
	"strconv"
	"strings"
)

type ORMQuery struct {
	DB     *sql.DB
	fields ORMField
	sql    bytes.Buffer
	params []interface{}
	IsLog  bool
}

func (q *ORMQuery) From(from string) *ORMQuery {
	q.fields.from = ORMFrom{From: from}
	return q
}

func (q *ORMQuery) Where(column string, columnType string, value interface{}) *ORMQuery {
	q.fields.where = append(q.fields.where, ORMWhere{
		Column: column,
		Type:   columnType,
		Value:  value,
	})
	return q
}

func (q *ORMQuery) Select(column string) *ORMQuery {
	split := strings.Split(column, ",")
	for _, s := range split {
		q.fields.selects = append(
			q.fields.selects,
			ORMSelect{column: s},
		)
	}
	return q
}

func (q *ORMQuery) build() (string, []interface{}) {
	buf := q.sql

	buf.WriteString("SELECT ")

	for _, ormSelect := range q.fields.selects {
		buf.WriteString(fmt.Sprintf("%s,", ormSelect.column))
	}

	if len(q.fields.selects) > 0 {
		buf.Truncate(buf.Len() - 1)
	}

	if len(q.fields.selects) == 0 {
		buf.WriteString("*")
	}

	buf.WriteString(" FROM ")
	buf.WriteString(q.fields.from.From)

	for _, join := range q.fields.join {
		buf.WriteString(
			fmt.Sprintf(
				" %s %s ON %s ",
				join.joinType,
				join.table,
				join.on,
			),
		)
	}

	buf.WriteString(" WHERE ")
	if len(q.fields.where) != 0 {
		for _, ormWhere := range q.fields.where {
			var value interface{} = "$" + strconv.Itoa(len(q.params)+1)
			if ormWhere.ValueX {
				value = ormWhere.Value
			}
			buf.WriteString(
				fmt.Sprintf(" %s %s %s %v",
					ormWhere.Condition,
					ormWhere.Column,
					ormWhere.Type,
					value,
				),
			)
			if !ormWhere.ValueX {
				q.params = append(q.params, ormWhere.Value)
			}
		}
	}

	delimiterAnd := ""
	if len(q.fields.where) > 0 {
		delimiterAnd = " AND"
	}

	buf.WriteString(fmt.Sprintf("%s %s.%s", delimiterAnd, q.fields.from.From, "deleted_at is null "))

	sqlBuild := buf.String()

	if q.IsLog {
		fmt.Println(sqlBuild)
		lib.PrintPre(q.params)
	}

	return sqlBuild, q.params
}

func (q *ORMQuery) FindAll() (string, []interface{}) {
	return q.build()
}

func (q *ORMQuery) Build() (string, []interface{}) {
	return q.build()
}

func (q *ORMQuery) Find(model interface{}) (interface{}, error) {
	sqlString, params := q.build()
	rows, _ := q.DB.Query(sqlString, params...)
	cols, _ := rows.Columns()

	var store []interface{}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		store = append(store, m)
	}

	if len(store) == 0 {
		return nil, nil
	}
	jsonbody, _ := json.Marshal(store[0])
	err := json.Unmarshal(jsonbody, &model)
	return model, err
}

func (q *ORMQuery) WhereX(column string, columnType string, value interface{}) *ORMQuery {
	q.fields.where = append(q.fields.where, ORMWhere{
		Column: column,
		Type:   columnType,
		Value:  value,
		ValueX: true,
	})
	return q
}

func (q *ORMQuery) SelectSubQuery(sql string, as string) *ORMQuery {
	q.fields.selects = append(
		q.fields.selects,
		ORMSelect{column: fmt.Sprintf(" (%s) as %s", sql, as)},
	)
	return q
}

func (q *ORMQuery) JOIN(table string, on string) *ORMQuery {
	q.fields.join = append(
		q.fields.join,
		ORMJoin{
			table:    table,
			on:       on,
			joinType: "INNER JOIN",
		},
	)
	return q
}

func (q *ORMQuery) AndWhere(column string, columnType string, value interface{}) *ORMQuery {
	q.fields.where = append(q.fields.where, ORMWhere{
		Column:    column,
		Type:      columnType,
		Value:     value,
		Condition: "AND",
	})
	return q
}
