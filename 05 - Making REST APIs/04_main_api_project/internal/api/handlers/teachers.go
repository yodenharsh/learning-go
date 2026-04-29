package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"strconv"
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	dbParams := map[string]string{
		"firstName": "first_name",
		"lastName":  "last_name",
		"email":     "email",
		"class":     "class",
		"subject":   "subject",
	}

	teacherList, err := sqlconnect.GetTeachers(dbParams, r)
	if err != nil {
		http.Error(w, "Error retrieving teachers", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teacherList),
		Data:   teacherList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func GetTeacherByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusUnprocessableEntity)
		return
	}

	teacher, err := sqlconnect.GetTeacherById(id)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func PostTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	addedTeachers, err := sqlconnect.AddTeacher(newTeachers)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error adding teachers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateTeachersHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusUnprocessableEntity)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	err = sqlconnect.UpdateTeacherById(id, updatedTeacher)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

func PatchTeacherByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusUnprocessableEntity)
		return
	}

	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	db := sqlconnect.ConnectDb()
	defer db.Close()

	type IdHolder struct {
		id int
	}

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.Id, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error checking teacher existence", http.StatusInternalServerError)
		return
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

	if err != nil {
		http.Error(w, "Error updating teacher in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)
}

func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlconnect.ConnectDb()
	defer db.Close()

	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}

	for _, update := range updates {
		stringId, ok := update["id"].(string)

		if !ok {
			tx.Rollback()
			http.Error(w, "Invalid or missing teacher ID in update", http.StatusUnprocessableEntity)
			return
		}

		id, err := strconv.Atoi(stringId)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Invalid teacher ID", http.StatusUnprocessableEntity)
			return
		}

		var teacherFromDb models.Teacher

		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
			Scan(&teacherFromDb.Id, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)
		if err == sql.ErrNoRows {
			tx.Rollback()
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		} else if err != nil {
			tx.Rollback()
			http.Error(w, "Error checking teacher existence", http.StatusInternalServerError)
			return
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
								log.Printf("Cannot convert %v to %v", val.Type(), fieldVal.Type())
							}
						}
					}
				}
			}
		}

		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, id)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			http.Error(w, "Error updating teacher in database", http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Something went wrong when updating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteTeacherByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusUnprocessableEntity)
		return
	}

	db := sqlconnect.ConnectDb()
	defer db.Close()

	query := "DELETE FROM teachers WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting teacher from database", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking deletion result", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlconnect.ConnectDb()
	defer db.Close()

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}

	deletedIds := []int{}
	for _, id := range ids {
		execStmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
		if err != nil {
			log.Println(err)
			tx.Rollback()
			http.Error(w, "Error preparing delete statement", http.StatusInternalServerError)
			return
		}
		defer execStmt.Close()

		result, err := execStmt.Exec(id)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			http.Error(w, "Error deleting teacher from database", http.StatusInternalServerError)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Println(err)
			tx.Rollback()
			http.Error(w, "Error checking deletion result", http.StatusInternalServerError)
			return
		}
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}

	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	if len(deletedIds) < 1 {
		http.Error(w, "No teachers found for the provided IDs", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		DeletedIds []int  `json:"deletedIds"`
	}{
		Status:     "success",
		DeletedIds: deletedIds,
	}

	json.NewEncoder(w).Encode(response)
}
