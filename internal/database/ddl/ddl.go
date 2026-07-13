// internal/database/ddl/ddl.go
package ddl

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"github.com/fuadop/beacon/internal/database/models"
)

func Create[T models.TableRepresenter](db *sql.DB) error {
	var model T
	tableName := model.TableName()
	val := reflect.TypeOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var columns []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToLower(field.Name)
		}

		sqlType := "TEXT"
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int64, reflect.Bool:
			sqlType = "INTEGER"
		case reflect.Float64, reflect.Float32:
			sqlType = "REAL"
		}

		if field.Tag.Get("pk") == "true" {
			sqlType += " PRIMARY KEY AUTOINCREMENT"
		}
		columns = append(columns, fmt.Sprintf("%s %s", dbTag, sqlType))
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n\t%s\n);", tableName, strings.Join(columns, ",\n\t"))
	_, err := db.Exec(query)
	return err
}

func Drop[T models.TableRepresenter](db *sql.DB) error {
	var model T
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", model.TableName())
	_, err := db.Exec(query)
	return err
}
