package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"so-crm/src/libs/ORM/lib"
	"strings"
)

type ORMMutation struct {
	DB     *sql.DB
	IsLog  bool
	fields ORMMutationField
	sql    bytes.Buffer
}

func (q *ORMMutation) Into(into string) *ORMMutation {
	q.fields.Into = ORMInto{Into: into}

	return q
}
func (q *ORMMutation) Where(column string, condition string, value interface{}) *ORMMutation {
	q.fields.Where = append(q.fields.Where, ORMWhere{
		Column:    column,
		Condition: condition,
		Value:     value,
	})
	return q
}
func (q *ORMMutation) getRecursiveDataSet(mutable reflect.Value, typeOfT reflect.Type) {
	for i := 0; i < mutable.NumField(); i++ {
		f := mutable.Field(i)
		if f.Type().String() != "ORM.Model" {
			q.Set(typeOfT.Field(i).Name, f.Type().String(), f.Interface())
		} else {
			t := f.Type()
			q.getRecursiveDataSet(f, t)
		}
	}
}
func (q *ORMMutation) TypeCreate(create bool) *ORMMutation {
	q.fields.IsCreate = create
	return q
}
func (q *ORMMutation) Set(column string, columnType string, value interface{}) *ORMMutation {
	q.fields.Set = append(q.fields.Set, ORMSet{
		Column:     column,
		ColumnType: columnType,
		Value:      value,
	})
	return q
}
func (q *ORMMutation) Return(isReturn bool) *ORMMutation {
	q.fields.IsReturn = isReturn
	return q
}
func (q *ORMMutation) build() (string, []interface{}) {
	sqlBuffer := q.sql
	into := q.fields.Into.Into
	if q.fields.IsCreate {
		sqlBuffer.WriteString(fmt.Sprintf("INSERT INTO %s", into))
	} else {
		sqlBuffer.WriteString(fmt.Sprintf("UPDATE %s", into))
	}
	var columns []string
	var allColumns []string
	var values []interface{}
	var valuesNumbers []string
	var whereData []string

	noUpdateFields := []string{"Id", "CreatedAt"}

	for _, set := range q.fields.Set {
		field := lib.ToSnakeCase(set.Column)
		allColumns = append(allColumns, field)
		value := set.Value
		if q.fields.IsCreate {
			if set.Column == "DeletedAt" || set.Column == "Id" {
				continue
			}
			columns = append(columns, field)
			valuesNumbers = append(valuesNumbers, fmt.Sprintf("$%d", len(valuesNumbers)+1))
		} else {
			if isExist, _ := lib.InArray(set.Column, noUpdateFields); isExist {
				continue
			}
			columns = append(columns, fmt.Sprintf("%s=$%d", field, len(values)+1))
		}
		values = append(values, value)
	}

	if q.fields.IsCreate {
		sqlBuffer.WriteString(
			fmt.Sprintf(
				"(%s) VALUES(%s)",
				strings.Join(columns, ", "),
				strings.Join(valuesNumbers, ", "),
			),
		)
	} else {
		sqlBuffer.WriteString(
			fmt.Sprintf(
				" SET %s",
				strings.Join(columns, ", "),
			),
		)
	}

	if len(q.fields.Where) > 0 {
		sqlBuffer.WriteString(" WHERE ")
	}

	for _, where := range q.fields.Where {
		whereData = append(whereData, fmt.Sprintf("%s %s $%d", lib.ToSnakeCase(where.Column), where.Condition, len(values)+1))
		values = append(values, where.Value)
	}

	if len(q.fields.Where) > 0 {
		sqlBuffer.WriteString(strings.Join(whereData, " AND "))
	}

	if q.fields.IsReturn {
		sqlBuffer.WriteString(
			fmt.Sprintf(
				" RETURNING %s",
				strings.Join(allColumns, ", "),
			),
		)
	}
	sqlBuilder := sqlBuffer.String()
	if q.IsLog {
		_ = fmt.Sprintf("SQL => %s %v", sqlBuilder, values)
	}
	return sqlBuilder, values
}
func (q *ORMMutation) Save(model interface{}) error {
	inputs := make([]reflect.Value, 0)
	of := reflect.TypeOf(model)
	value := reflect.ValueOf(model)
	mutable := value.Elem()

	typeOfT := mutable.Type()

	modelElement := mutable.FieldByName("Model")
	isExistModel := mutable.FieldByName("Id")
	isCreate := isExistModel.IsValid() && isExistModel.Int() == 0

	name := of.Elem().Name()
	table := lib.ToSnakeCase(name)

	var methodBefore = "BeforeCreate"

	if !q.fields.IsCreate {
		methodBefore = "BeforeUpdate"
	}

	beforeCreate := reflect.ValueOf(model).MethodByName(methodBefore)
	if beforeCreate.IsValid() {
		beforeCreate.Call(inputs)
	}

	if !q.fields.IsCreate && isExistModel.IsValid() {
		q.Where("Id", "=", isExistModel.Interface())
	}

	q.Into(table).TypeCreate(isCreate)
	q.getRecursiveDataSet(mutable, typeOfT)

	sqlBuilder, params := q.build()

	if q.fields.IsReturn {
		rows, err := q.DB.Query(sqlBuilder, params...)
		if err != nil {
			return err
		}
		columns, err := rows.Columns()
		colNum := len(columns)
		var values = make([]interface{}, colNum)

		for i, _ := range values {
			var ii interface{}
			values[i] = &ii
		}
		columnsModel := []string{"Id", "CreatedAt", "UpdatedAt", "DeletedAt"}
		for rows.Next() {
			rows.Scan(values...)
			for i, colName := range columns {
				var raw_value = *(values[i].(*interface{}))
				//var raw_type = reflect.TypeOf(raw_value)
				var fieldName = lib.ToCamelCase(colName)
				isExistModel := mutable.FieldByName(fieldName)
				if !isExistModel.IsZero() {
					isExistModel.Set(reflect.ValueOf(raw_value))
				}
				if isExist, _ := lib.InArray(fieldName, columnsModel); isExist {
					byName := modelElement.FieldByName(fieldName)
					if raw_value != nil {
						byName.Set(reflect.ValueOf(raw_value))
					}
				}
			}
		}
		return nil
	}
	_, err := q.DB.Exec(sqlBuilder, params...)

	return err
}

func (q *ORMMutation) Delete(model interface{}) error {
	inputs := make([]reflect.Value, 0)
	of := reflect.TypeOf(model)
	value := reflect.ValueOf(model)
	mutable := value.Elem()

	fieldId := mutable.FieldByName("Id")
	fieldDeletedAt := mutable.FieldByName("DeletedAt")
	modelElement := mutable.FieldByName("Model")
	beforeDelete := reflect.ValueOf(model).MethodByName("BeforeDelete")

	if beforeDelete.IsValid() {
		beforeDelete.Call(inputs)
	}

	name := of.Elem().Name()
	table := lib.ToSnakeCase(name)

	q.Into(table)
	q.fields.IsCreate = false

	q.fields.Set = append(q.fields.Set, ORMSet{
		Column:     "DeletedAt",
		ColumnType: fieldDeletedAt.Type().String(),
		Value:      fieldDeletedAt.Interface(),
	})

	q.fields.Where = append(q.fields.Where, ORMWhere{
		Column:    "Id",
		Condition: "=",
		Value:     fieldId.Interface(),
		Type:      fieldId.Type().String(),
	})

	sqlBuilder, params := q.build()

	if q.fields.IsReturn {
		rows, err := q.DB.Query(sqlBuilder, params...)
		if err != nil {
			return err
		}
		columns, err := rows.Columns()
		colNum := len(columns)
		var values = make([]interface{}, colNum)

		for i, _ := range values {
			var ii interface{}
			values[i] = &ii
		}

		columnsModel := []string{"Id", "CreatedAt", "UpdatedAt", "DeletedAt"}
		for rows.Next() {
			rows.Scan(values...)
			for i, colName := range columns {
				var raw_value = *(values[i].(*interface{}))
				//var raw_type = reflect.TypeOf(raw_value)
				var fieldName = lib.ToCamelCase(colName)
				isExistModel := mutable.FieldByName(fieldName)
				if !isExistModel.IsZero() {
					isExistModel.Set(reflect.ValueOf(raw_value))
				}
				if isExist, _ := lib.InArray(fieldName, columnsModel); isExist {
					byName := modelElement.FieldByName(fieldName)
					if raw_value != nil {
						byName.Set(reflect.ValueOf(raw_value))
					}
				}
			}
		}

		return nil
	}
	_, err := q.DB.Exec(sqlBuilder, params...)

	return err
}
