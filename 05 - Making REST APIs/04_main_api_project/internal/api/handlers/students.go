package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"restapi/pkg/utils"
	"strconv"
)

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {

	dbParams := map[string]string{
		"firstName": "first_name",
		"lastName":  "last_name",
		"email":     "email",
		"class":     "class",
	}

	studentList, err := sqlconnect.GetStudents(dbParams, r)
	if err != nil {
		http.Error(w, "Error retrieving students", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(studentList),
		Data:   studentList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func GetStudentByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusUnprocessableEntity)
		return
	}

	student, err := sqlconnect.GetStudentById(id)
	if err == sql.ErrNoRows {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func PostStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var students []models.Student
	err := json.NewDecoder(r.Body).Decode(&students)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	for _, student := range students {
		val := reflect.ValueOf(student)
		err = utils.CheckStringFieldsNotEmpty(val)
		if err != nil {
			http.Error(w, "One or more required fields is empty/not provided", http.StatusBadRequest)
		}
	}

	addedStudents, err := sqlconnect.AddStudent(students)
	if err != nil {
		http.Error(w, "Error adding students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateStudentsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusUnprocessableEntity)
		return
	}

	var updatedStudent models.Student
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	updatedStudentFromDb, err := sqlconnect.UpdateStudentById(id, updatedStudent)
	if err == sql.ErrNoRows {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error updating student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudentFromDb)
}

func PatchStudentByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusUnprocessableEntity)
		return
	}

	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	existingStudent, err := sqlconnect.PatchStudentById(id, updates)
	if err == sql.ErrNoRows {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error updating student in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingStudent)
}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	err = sqlconnect.PatchStudents(updates)
	if err == sql.ErrNoRows {
		http.Error(w, "One or more students not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error updating students in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteStudentByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusUnprocessableEntity)
		return
	}

	err = sqlconnect.DeleteStudentById(id)
	if err == sql.ErrNoRows {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Fatalln(err)
		http.Error(w, "Error deleting student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	deletedIds, err := sqlconnect.DeleteStudents(ids)
	if err != nil {
		http.Error(w, "Error deleting students", http.StatusInternalServerError)
		return
	}

	if len(deletedIds) < 1 {
		http.Error(w, "No students found for the provided IDs", http.StatusNotFound)
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
