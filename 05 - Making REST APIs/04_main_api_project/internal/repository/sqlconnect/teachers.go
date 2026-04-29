package sqlconnect

import (
	"database/sql"
	"net/http"
	"reflect"
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

func AddTeacher(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db := ConnectDb()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			return nil, err
		}

		res.LastInsertId()
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		newTeacher.Id = int(lastId)
		addedTeachers[i] = newTeacher
	}
	return addedTeachers, nil
}

func UpdateTeacherById(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db := ConnectDb()
	defer db.Close()

	type IdHolder struct {
		id int
	}

	var existingTeacherId IdHolder
	err := db.QueryRow("SELECT id FROM teachers WHERE id = ?", id).Scan(&existingTeacherId.id)
	if err != nil {
		return models.Teacher{}, err
	}

	query := "UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?"
	_, err = db.Exec(query, &updatedTeacher.FirstName, &updatedTeacher.LastName, &updatedTeacher.Email, &updatedTeacher.Class, &updatedTeacher.Subject, id)

	if err != nil {
		return models.Teacher{}, err
	}
	return updatedTeacher, nil
}

func PatchTeacherById(id int, updates map[string]any) (models.Teacher, error) {
	db := ConnectDb()
	defer db.Close()

	var existingTeacher models.Teacher
	err := db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.Id, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		return models.Teacher{}, err
	}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					fieldVal := teacherVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	query := "UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?"
	_, err = db.Exec(query, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject, id)
	return existingTeacher, err
}

func DeleteTeacherById(id int) error {
	db := ConnectDb()
	defer db.Close()

	query := "DELETE FROM teachers WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
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
