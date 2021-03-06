package sqlx

import (
	"fmt"
	"go/ast"
	"reflect"

	"github.com/artisanhe/tools/reflectx"
	"github.com/artisanhe/tools/sqlx/builder"
)

func ForEachStructField(structType reflect.Type, fn func(structField reflect.StructField, columnName string)) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		if ast.IsExported(field.Name) {
			fieldName, exists := field.Tag.Lookup("db")
			if exists {
				if fieldName != "-" {
					fn(field, fieldName)
				}
			} else if field.Anonymous {
				ForEachStructField(field.Type, fn)
				continue
			}
		}
	}
}

func ForEachStructFieldValue(rv reflect.Value, fn func(structFieldValue reflect.Value, structField reflect.StructField, columnName string)) {
	rv = reflect.Indirect(rv)
	structType := rv.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if ast.IsExported(field.Name) {
			fieldValue := rv.Field(i)

			columnName, exists := field.Tag.Lookup("db")
			if exists {
				if columnName != "-" {
					fn(fieldValue, field, columnName)
				}
			} else if field.Anonymous {
				ForEachStructFieldValue(fieldValue, fn)
				continue
			}
		}
	}
}

func FieldValuesFromStructBy(structValue interface{}, fieldNames []string) (fieldValues builder.FieldValues) {
	fieldValues = builder.FieldValues{}
	rv := reflect.Indirect(reflect.ValueOf(structValue))
	fieldMap := FieldNames(fieldNames).Map()
	ForEachStructFieldValue(rv, func(structFieldValue reflect.Value, structField reflect.StructField, columnName string) {
		if fieldMap != nil && fieldMap[structField.Name] {
			fieldValues[structField.Name] = structFieldValue.Interface()
		}
	})
	return fieldValues
}

func FieldValuesFromStructByNonZero(structValue interface{}, excludes ...string) (fieldValues builder.FieldValues) {
	fieldValues = builder.FieldValues{}
	rv := reflect.Indirect(reflect.ValueOf(structValue))
	fieldMap := FieldNames(excludes).Map()
	ForEachStructFieldValue(rv, func(structFieldValue reflect.Value, structField reflect.StructField, columnName string) {
		if !reflectx.IsEmptyValue(structFieldValue) || (fieldMap != nil && fieldMap[structField.Name]) {
			fieldValues[structField.Name] = structFieldValue.Interface()
		}
	})
	return
}

func ScanDefToTable(rv reflect.Value, table *builder.Table) {
	rv = reflect.Indirect(rv)
	structType := rv.Type()

	ForEachStructField(structType, func(structField reflect.StructField, columnName string) {
		sqlType, exists := structField.Tag.Lookup("sql")
		if !exists {
			panic(fmt.Errorf("%s.%s sql tag must defined for sql type", table.Name, structField.Name))
		}

		col := builder.Col(table, columnName).Type(sqlType).Field(structField.Name)

		if structField.Type.AssignableTo(reflect.TypeOf((*EnumTypeDescriber)(nil)).Elem()) {
			enumTypeDescriber := reflect.New(structField.Type).Interface().(EnumTypeDescriber)
			col = col.Enum(enumTypeDescriber.EnumType(), enumTypeDescriber.Enums())
		}

		finalSql := col.ColumnType.String()
		if sqlType != finalSql {
			panic(fmt.Errorf("%s.%s sql tag may be `%s`, current `%s`", table.Name, structField.Name, finalSql, sqlType))
		}
		table.Columns.Add(col)
	})

	if rv.CanAddr() {
		addr := rv.Addr()
		if addr.CanInterface() {
			i := addr.Interface()

			if primaryKeyHook, ok := i.(WithPrimaryKey); ok {
				primaryKey := builder.PrimaryKey()
				for _, fieldName := range primaryKeyHook.PrimaryKey() {
					if col := table.F(fieldName); col != nil {
						primaryKey = primaryKey.WithCols(col)
					} else {
						panic(fmt.Errorf("field %s for PrimaryKey is not defined in table model %s", fieldName, table.Name))
					}
				}
				table.Keys.Add(primaryKey)
			}

			if withComments, ok := i.(WithComments); ok {
				for fieldName, comment := range withComments.Comments() {
					field := table.F(fieldName)
					if field != nil {
						field.Comment = comment
					}
				}
			}

			if indexesHook, ok := i.(WithIndexes); ok {
				for name, indexes := range indexesHook.Indexes() {
					idx := builder.Index(name)
					for _, fieldName := range indexes {
						if col := table.F(fieldName); col != nil {
							idx = idx.WithCols(col)
						} else {
							panic(fmt.Errorf("field %s for key %s is not defined in table model %s", fieldName, name, table.Name))
						}
					}
					table.Keys.Add(idx)
				}
			}

			if uniqueIndexesHook, ok := i.(WithUniqueIndexes); ok {
				for name, indexes := range uniqueIndexesHook.UniqueIndexes() {
					idx := builder.UniqueIndex(name)
					for _, indexName := range indexes {
						if col := table.F(indexName); col != nil {
							idx = idx.WithCols(col)
						} else {
							panic(fmt.Errorf("field %s for unique indexes %s is not defined in table model %s", indexName, name, table.Name))
						}
					}
					table.Keys.Add(idx)
				}
			}
		}
	}
}
