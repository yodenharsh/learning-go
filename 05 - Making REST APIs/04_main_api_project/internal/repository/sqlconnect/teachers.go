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
)

func GetTeachers(dbParams map[string]string, r *http.Request) ([]models.Teacher, error) {
	db := ConnectDb()
	defer db.Close()
	teacherList := make([]models.Teacher, 0)
	query, args := utils.BuildQueryWithFilters(r, "teachers", dbParams)

	query = utils.BuildQueryWithSorting(r, query, dbParams)

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

	fmt.Println(utils.GenerateInsertQuery("teachers", models.Teacher{}))
	stmt, err := db.Prepare(utils.GenerateInsertQuery("teachers", models.Teacher{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Preparing insert stmt failde")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		values := utils.GetStructValues(newTeacher)
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

func GetStudentsOfTeacher(teacherId int) ([]models.Student, error) {
	db := ConnectDb()
	defer db.Close()

	type IdAndClassHolder struct {
		Id    int
		Class string
	}

	var teacherIdAndClass IdAndClassHolder
	err := db.QueryRow("SELECT id, class FROM teachers where id = ?", teacherId).Scan(&teacherIdAndClass.Id, &teacherIdAndClass.Class)
	if err == sql.ErrNoRows {
		return nil, utils.ErrorHandler(err, "teacher does not exist")
	} else if err != nil {
		return nil, utils.ErrorHandler(err, "Error querying teacher")
	}

	rows, err := db.Query("SELECT id, first_name, last_name, email, class FROM students WHERE class = ?", teacherIdAndClass.Class)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error querying students")
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error scanning student row")
		}
		students = append(students, student)
	}

	return students, nil
}
