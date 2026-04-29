package sqlconnect

import (
	"net/http"
	"restapi/internal/models"
	"slices"
	"strings"
)

func GetTeachers(dbParams map[string]string, r *http.Request) ([]models.Teacher, error) {
	db := ConnectDb()
	defer db.Close()
	teacherList := make([]models.Teacher, 0)
	query, args := buildQueryWithFilters(r, dbParams)

	query = buildQueryWithSorting(r, query, dbParams)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, err
		}
		teacherList = append(teacherList, teacher)
	}
	return teacherList, nil
}

func GetTeacherById(id int) (models.Teacher, error) {
	db := ConnectDb()
	defer db.Close()

	var teacher models.Teacher
	err := db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err != nil {
		return models.Teacher{}, err
	}
	return teacher, nil
}

func buildQueryWithFilters(r *http.Request, dbParams map[string]string) (string, []any) {
	var query strings.Builder
	query.WriteString("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1")
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

func buildQueryWithSorting(r *http.Request, query string, dbParams map[string]string) string {
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
