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
	"strconv"
	"strings"
)

func GetStudents(dbParams map[string]string, r *http.Request) ([]models.Student, error) {
	db := ConnectDb()
	defer db.Close()
	studentList := make([]models.Student, 0)

	query := "SELECT id, first_name, last_name, email, class FROM students WHERE 1=1"
	query, args := utils.BuildQueryWithFilters(r, query, dbParams)

	query = utils.BuildQueryWithSorting(r, query, dbParams)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error querying DB")
	}
	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Scan row failed")
		}
		studentList = append(studentList, student)
	}
	return studentList, nil
}

func GetStudentById(id int) (models.Student, error) {
	db := ConnectDb()
	defer db.Close()

	var student models.Student
	err := db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Querying students failed")
	}
	return student, nil
}

func AddStudent(newStudents []models.Student) ([]models.Student, error) {
	db := ConnectDb()
	defer db.Close()

	fmt.Println(utils.GenerateInsertQuery("students", models.Student{}))
	stmt, err := db.Prepare(utils.GenerateInsertQuery("students", models.Student{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Preparing insert stmt failde")
	}
	defer stmt.Close()

	addedStudents := make([]models.Student, len(newStudents))
	for i, newStudent := range newStudents {
		values := utils.GetStructValues(newStudent)
		res, err := stmt.Exec(values...)
		if err != nil {
			if strings.Contains(err.Error(), "Error 1452") {
				return nil, utils.ErrorHandler(err, "Class does not exist")
			}
			return nil, utils.ErrorHandler(err, "Failed to add student")
		}

		res.LastInsertId()
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to check last inserted ID")
		}

		newStudent.Id = int(lastId)
		addedStudents[i] = newStudent
	}
	return addedStudents, nil
}

func UpdateStudentById(id int, updatedStudent models.Student) (models.Student, error) {
	db := ConnectDb()
	defer db.Close()

	type IdHolder struct {
		id int
	}

	var existingStudentId IdHolder
	err := db.QueryRow("SELECT id FROM students WHERE id = ?", id).Scan(&existingStudentId.id)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error when querying for existing ID")
	}

	query := "UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?"
	_, err = db.Exec(query, &updatedStudent.FirstName, &updatedStudent.LastName, &updatedStudent.Email, &updatedStudent.Class, id)

	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Failed to update student")
	}
	return updatedStudent, nil
}

func PatchStudentById(id int, updates map[string]any) (models.Student, error) {
	db := ConnectDb()
	defer db.Close()

	var existingStudent models.Student
	err := db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&existingStudent.Id, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email, &existingStudent.Class)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Failed to do existing check")
	}

	studentVal := reflect.ValueOf(&existingStudent).Elem()
	studentType := studentVal.Type()

	for k, v := range updates {
		for i := 0; i < studentVal.NumField(); i++ {
			field := studentType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if studentVal.Field(i).CanSet() {
					fieldVal := studentVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(studentVal.Field(i).Type()))
				}
			}
		}
	}

	query := "UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?"
	_, err = db.Exec(query, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email, &existingStudent.Class, id)
	return existingStudent, err
}

func PatchStudents(updates []map[string]any) error {
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
			return utils.ErrorHandler(errors.New("Invalid or missing student ID"), "Invalid or missing student ID")
		}

		id, err := strconv.Atoi(stringId)
		if err != nil {
			tx.Rollback()
			return err
		}

		var studentFromDb models.Student

		err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).
			Scan(&studentFromDb.Id, &studentFromDb.FirstName, &studentFromDb.LastName, &studentFromDb.Email, &studentFromDb.Class)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Failed to get existing students")
		}
		studentVal := reflect.ValueOf(&studentFromDb).Elem()
		studentType := studentVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}

			for i := 0; i < studentVal.NumField(); i++ {
				field := studentType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					if studentVal.Field(i).CanSet() {
						fieldVal := studentVal.Field(i)
						if fieldVal.CanSet() {
							val := reflect.ValueOf(v)
							if val.Type().ConvertibleTo(fieldVal.Type()) {
								fieldVal.Set(val.Convert(fieldVal.Type()))
							} else {
								tx.Rollback()
								return utils.ErrorHandler(errors.New("Error updating student in database"), "Error updating student in database")
							}
						}
					}
				}
			}
		}
		_, err = tx.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", studentFromDb.FirstName, studentFromDb.LastName, studentFromDb.Email, studentFromDb.Class, id)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error updating student in database")
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "Error updating transaction")
	}
	return nil
}

func DeleteStudents(ids []int) ([]int, error) {
	db := ConnectDb()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	deletedIds := []int{}
	for _, id := range ids {
		execStmt, err := tx.Prepare("DELETE FROM students WHERE id = ?")
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

func DeleteStudentById(id int) error {
	db := ConnectDb()
	defer db.Close()

	query := "DELETE FROM students WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to delete student")
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
