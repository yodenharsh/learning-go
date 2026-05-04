package sqlconnect

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"restapi/internal/models"
	"restapi/pkg/utils"
	"slices"
	"strconv"
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
		return nil, utils.ErrorHandler(err, "Error querying DB")
	}
	defer rows.Close()

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Scan row failed")
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
		return models.Teacher{}, utils.ErrorHandler(err, "Querying teachers failed")
	}
	return teacher, nil
}

func AddTeacher(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db := ConnectDb()
	defer db.Close()

	fmt.Println(generateInsertQuery(models.Teacher{}))
	stmt, err := db.Prepare(generateInsertQuery(models.Teacher{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Preparing insert stmt failde")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		values := getStructValues(newTeacher)
		res, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to add teacher")
		}

		res.LastInsertId()
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to check last inserted ID")
		}

		newTeacher.Id = int(lastId)
		addedTeachers[i] = newTeacher
	}
	return addedTeachers, nil
}

func generateInsertQuery(model any) string {
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
	return fmt.Sprintf("INSERT INTO teachers (%s) VALUES (%s)", columns, placeholders)
}

func getStructValues(model any) []any {
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

func UpdateTeacherById(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db := ConnectDb()
	defer db.Close()

	type IdHolder struct {
		id int
	}

	var existingTeacherId IdHolder
	err := db.QueryRow("SELECT id FROM teachers WHERE id = ?", id).Scan(&existingTeacherId.id)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Error when querying for existing ID")
	}

	query := "UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?"
	_, err = db.Exec(query, &updatedTeacher.FirstName, &updatedTeacher.LastName, &updatedTeacher.Email, &updatedTeacher.Class, &updatedTeacher.Subject, id)

	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Failed to update teacher")
	}
	return updatedTeacher, nil
}

func PatchTeacherById(id int, updates map[string]any) (models.Teacher, error) {
	db := ConnectDb()
	defer db.Close()

	var existingTeacher models.Teacher
	err := db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.Id, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Failed to do existing check")
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

func PatchTeachers(updates []map[string]any) error {
	db := ConnectDb()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "Failed to start transaction")
	}

	for _, update := range updates {
		stringId, ok := update["id"].(string)

		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(errors.New("Invalid or missing teacher ID"), "Invalid or missing teacher ID")
		}

		id, err := strconv.Atoi(stringId)
		if err != nil {
			tx.Rollback()
			return err
		}

		var teacherFromDb models.Teacher

		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
			Scan(&teacherFromDb.Id, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Failed to get existing teachers")
		}
		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}

			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					if teacherVal.Field(i).CanSet() {
						fieldVal := teacherVal.Field(i)
						if fieldVal.CanSet() {
							val := reflect.ValueOf(v)
							if val.Type().ConvertibleTo(fieldVal.Type()) {
								fieldVal.Set(val.Convert(fieldVal.Type()))
							} else {
								tx.Rollback()
								return utils.ErrorHandler(errors.New("Error updating teacher in database"), "Error updating teacher in database")
							}
						}
					}
				}
			}
		}
		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, id)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error updating teacher in database")
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "Error updating transaction")
	}
	return nil
}

func DeleteTeachers(ids []int) ([]int, error) {
	db := ConnectDb()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	deletedIds := []int{}
	for _, id := range ids {
		execStmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Failed to prepare delete statement")
		}
		defer execStmt.Close()

		result, err := execStmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Failed to execute DELETE")
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Failed to get rows affected")
		}
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}

	}
	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Couldn't commit transaction")
	}
	return deletedIds, nil
}

func DeleteTeacherById(id int) error {
	db := ConnectDb()
	defer db.Close()

	query := "DELETE FROM teachers WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to delete teacher")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Failed to get rows affected")
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
