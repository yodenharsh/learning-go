package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"
)

func GenerateInsertQuery(tableName string, model any) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag != "" && dbTag != "id" {
			if columns != "" {
				columns += " ," + dbTag
				placeholders += ",?"
			} else {
				columns += dbTag
				placeholders += "?"
			}
		}
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)
}

func GetStructValues(model any) []any {
	modelVal := reflect.ValueOf(model)
	modelType := modelVal.Type()
	var values []any
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag != "" && dbTag != "id" {
			values = append(values, modelVal.Field(i).Interface())
		}
	}

	return values
}

func BuildQueryWithFilters(r *http.Request, initialQuery string, dbParams map[string]string) (string, []any) {
	var query strings.Builder
	query.WriteString(initialQuery)

	var args []any

	for param, dbField := range dbParams {
		value := r.URL.Query().Get(param)
		if value != "" {
			query.WriteString(" AND " + dbField + " = ?")
			args = append(args, value)
		}
	}
	return query.String(), args
}

func BuildQueryWithSorting(r *http.Request, query string, dbParams map[string]string) string {
	sortParam := r.URL.Query()["sortBy"]
	if len(sortParam) > 0 {
		query += " ORDER BY "
		for i, param := range sortParam {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}

			field, order := parts[0], parts[1]
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ", "
			}
			query += dbParams[field] + " " + order
		}
	}
	return query
}

func isValidSortField(field string) bool {
	validFields := []string{"firstName", "lastName", "email", "class", "subject"}
	return slices.Contains(validFields, field)
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}
