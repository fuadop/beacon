// internal/database/dml/dml.go
package dml

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"yourproject/internal/database/models"
)

func Insert[T models.TableRepresenter](db *sql.DB, item *T) error {
	val := reflect.ValueOf(item).Elem()
	typ := val.Type()
	tableName := (*item).TableName()

	var cols []string
	var placeholders []string
	var args []any

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.Tag.Get("pk") == "true" {
			continue
		}

		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToLower(field.Name)
		}
		cols = append(cols, dbTag)
		placeholders = append(placeholders, "?")
		args = append(args, val.Field(i).Interface())
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	result, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	if lastID, err := result.LastInsertId(); err == nil {
		for i := 0; i < val.NumField(); i++ {
			if typ.Field(i).Tag.Get("pk") == "true" && val.Field(i).CanSet() {
				val.Field(i).SetInt(lastID)
			}
		}
	}
	return nil
}

func Get[T models.TableRepresenter](db *sql.DB, id any) (*T, error) {
	var item T
	val := reflect.ValueOf(&item).Elem()
	typ := val.Type()
	tableName := item.TableName()

	var pkCol string
	var scanTargets []any

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToLower(field.Name)
		}
		if field.Tag.Get("pk") == "true" {
			pkCol = dbTag
		}
		scanTargets = append(scanTargets, val.Field(i).Addr().Interface())
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1", tableName, pkCol)
	err := db.QueryRow(query, id).Scan(scanTargets...)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func Update[T models.TableRepresenter](db *sql.DB, item *T) error {
	val := reflect.ValueOf(item).Elem()
	typ := val.Type()
	tableName := (*item).TableName()

	var setClauses []string
	var args []any
	var pkCol string
	var pkVal any

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToLower(field.Name)
		}

		if field.Tag.Get("pk") == "true" {
			pkCol = dbTag
			pkVal = val.Field(i).Interface()
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", dbTag))
		args = append(args, val.Field(i).Interface())
	}

	if pkCol == "" {
		return fmt.Errorf("missing explicit primary key")
	}
	args = append(args, pkVal)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", tableName, strings.Join(setClauses, ", "), pkCol)
	_, err := db.Exec(query, args...)
	return err
}

func Wipe[T models.TableRepresenter](db *sql.DB) error {
	var model T
	query := fmt.Sprintf("DELETE FROM %s", model.TableName())
	_, err := db.Exec(query)
	return err
}

// Count generic utility returns the total number of entries sitting inside the model's table
func Count[T models.TableRepresenter](db *sql.DB) (int, error) {
	var model T
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", model.TableName())

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetAll pulls all existing records out of a target model table dynamically
func GetAll[T models.TableRepresenter](db *sql.DB) ([]T, error) {
	var model T
	tableName := model.TableName()

	// Reflect structure values to dynamically prepare scan destinations
	val := reflect.ValueOf(&model).Elem()
	typ := val.Type()

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T

	for rows.Next() {
		// Initialize a completely new element instance instance for each row iteration
		var item T
		itemVal := reflect.ValueOf(&item).Elem()

		var scanTargets []any
		for i := 0; i < itemVal.NumField(); i++ {
			scanTargets = append(scanTargets, itemVal.Field(i).Addr().Interface())
		}

		if err := rows.Scan(scanTargets...); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
