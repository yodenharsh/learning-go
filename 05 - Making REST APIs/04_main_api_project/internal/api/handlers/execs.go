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

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {

	dbParams := map[string]string{
		"firstName":             "first_name",
		"lastName":              "last_name",
		"email":                 "email",
		"username":              "username",
		"password":              "password",
		"inactiveStatus":        "inactive_status",
		"role":                  "role",
		"passwordResetCode":     "password_reset_code",
		"passwordCodeExpiresAt": "password_code_expires_at",
		"passwordChangedAt":     "password_changed_at",
		"userCreatedAt":         "user_created_at",
	}

	execsList, err := sqlconnect.GetExecs(dbParams, r)
	if err != nil {
		http.Error(w, "Error retrieving execs", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(execsList),
		Data:   execsList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func GetExecByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusUnprocessableEntity)
		return
	}

	exec, err := sqlconnect.GetExecById(id)
	if err == sql.ErrNoRows {
		http.Error(w, "Exec not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exec)
}

func PostExecsHandler(w http.ResponseWriter, r *http.Request) {

	var execs []models.Exec
	err := json.NewDecoder(r.Body).Decode(&execs)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	for idx, exec := range execs {
		val := reflect.ValueOf(exec)
		err = utils.CheckStringFieldsNotEmpty(val)
		if err != nil {
			http.Error(w, "One or more required fields is empty/not provided", http.StatusBadRequest)
			return
		}

		execs[idx].Password, err = utils.HashAndEncodePassword(exec.Password)
		if err != nil {
			http.Error(w, "Error adding data", http.StatusInternalServerError)
			return
		}
	}

	addedExecs, err := sqlconnect.AddExec(execs)
	if err != nil {
		http.Error(w, "Error adding execs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}

	json.NewEncoder(w).Encode(response)
}

func PatchExecByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exec ID", http.StatusUnprocessableEntity)
		return
	}

	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	existingExec, err := sqlconnect.PatchExecById(id, updates)
	if err == sql.ErrNoRows {
		http.Error(w, "Exec not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error updating exec in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingExec)
}

func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	err = sqlconnect.PatchExecs(updates)
	if err == sql.ErrNoRows {
		http.Error(w, "One or more execs not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error updating execs in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteExecByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exec ID", http.StatusUnprocessableEntity)
		return
	}

	err = sqlconnect.DeleteExecById(id)
	if err == sql.ErrNoRows {
		http.Error(w, "Exec not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Fatalln(err)
		http.Error(w, "Error deleting exec", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Exec

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	execInDb, err := sqlconnect.GetExecByUsername(req.Username)
	if err == sql.ErrNoRows {
		http.Error(w, "Incorrect credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	if execInDb.InactiveStatus {
		http.Error(w, "Account is inactive", http.StatusForbidden)
	}

	passwordMatch, err := utils.CompareHashedEncodedPassword(execInDb.Password, req.Password)
	if err != nil {
		http.Error(w, "Error comparing passwords", http.StatusInternalServerError)
	}

	if !passwordMatch {
		http.Error(w, "Incorrect credentials", http.StatusUnauthorized)
	}

	token, err := utils.SignToken(strconv.Itoa(execInDb.Id), execInDb.Username, execInDb.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	loginResponse := struct {
		Token string `json:"token"`
	}{
		Token: *token,
	}

	json.NewEncoder(w).Encode(&loginResponse)
	w.WriteHeader(http.StatusOK)
}
