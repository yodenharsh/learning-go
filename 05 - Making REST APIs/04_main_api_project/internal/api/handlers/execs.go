package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"restapi/internal/api/middlewares"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"restapi/pkg/utils"
	"strconv"
	"time"
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
		return
	}

	passwordMatch, err := utils.CompareHashedEncodedPassword(execInDb.Password, req.Password)
	if err != nil {
		http.Error(w, "Error comparing passwords", http.StatusInternalServerError)
		return
	}

	if !passwordMatch {
		http.Error(w, "Incorrect credentials", http.StatusUnauthorized)
		return
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

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    *token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&loginResponse)
	w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	parsedIdFromPath, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(w, "id should be an integer", http.StatusBadRequest)
		return
	}

	idFromCookie := r.Context().Value(middlewares.ContextKey("id"))
	if idFromCookie != idString {
		fmt.Println("id from cookie ", idFromCookie, "\nParsed ID: ", parsedIdFromPath)
		http.Error(w, "You can only change your own password", http.StatusUnauthorized)
		return
	}

	var updatePasswordRequestBody models.UpdatePasswordRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&updatePasswordRequestBody)

	if err != nil {
		http.Error(w, "Request body is well-formed", http.StatusBadRequest)
		return
	}

	err = utils.CheckStringFieldsNotEmpty(reflect.ValueOf(updatePasswordRequestBody))
	if err != nil {
		http.Error(w, "One or more fields is empty", http.StatusBadRequest)
		return
	}

	execInDb, err := sqlconnect.GetExecById(parsedIdFromPath)
	if err == sql.ErrNoRows {
		http.Error(w, "Exec not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	passwordMatch, err := utils.CompareHashedEncodedPassword(execInDb.Password, updatePasswordRequestBody.Password)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	} else if !passwordMatch {
		http.Error(w, "Current password doesn't match", http.StatusUnauthorized)
		return
	}

	hashedAndEncodedPassword, err := utils.HashAndEncodePassword(updatePasswordRequestBody.NewPassword)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Assigning it first to the struct ensures type-safety
	execInDb.Password = hashedAndEncodedPassword
	execInDb.PasswordChangedAt = sql.NullString{Valid: true, String: time.Now().String()}

	updatesMap := map[string]any{
		"password":          execInDb.Password,
		"passwordChangedAt": execInDb.PasswordChangedAt,
	}

	_, err = sqlconnect.PatchExecById(parsedIdFromPath, updatesMap)
	if err != nil {
		http.Error(w, "Error updating password in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Request body is not well-formed", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Email == "" {
		http.Error(w, "email field is empty", http.StatusBadRequest)
		return
	}

	execInDb, err := sqlconnect.GetExecByEmail(req.Email)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	generatedToken, err := utils.CreatePasswordResetToken()
	if err != nil {
		http.Error(w, "Something went wrong when sending email", http.StatusInternalServerError)
	}

	err = sqlconnect.UpdatePasswordResetCode(execInDb.Id, generatedToken.HashedValue, generatedToken.ExpiresAt.Format(time.RFC3339))
	if err != nil {
		http.Error(w, "Something went wrong when sending email", http.StatusInternalServerError)
		return
	}

	err = utils.SendPasswordResetEmail(generatedToken.Value, req.Email)
	if err != nil {
		http.Error(w, "Something went wrong when trying to send email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	resetCode := r.PathValue("resetCode")
	if resetCode == "" {
		http.Error(w, "Reset code is required", http.StatusBadRequest)
		return
	}

	var req struct {
		NewPassword string `json:"newPassword"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.NewPassword == "" {
		http.Error(w, "Request body has no newPassword field", http.StatusBadRequest)
		return
	}

	bytes, err := hex.DecodeString(resetCode)
	if err != nil {
		http.Error(w, "Invalid reset code", http.StatusBadRequest)
		return
	}
	hashedToken := sha256.Sum256(bytes)
	hashedTokenString := hex.EncodeToString(hashedToken[:])

	execInDb, err := sqlconnect.GetExecByResetCode(hashedTokenString)
	if err != nil {
		http.Error(w, "Invalid error code", http.StatusBadRequest)
		return
	}

	if !execInDb.PasswordCodeExpiresAt.Valid {
		http.Error(w, "Reset code has no expiry time, contact support", http.StatusInternalServerError)
		return
	}

	expiryTime, err := time.Parse(time.RFC3339, execInDb.PasswordCodeExpiresAt.String)
	if time.Now().After(expiryTime) {
		http.Error(w, "Reset code has expired", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashAndEncodePassword(req.NewPassword)
	if err != nil {
		utils.ErrorHandler(err, "Something went wrong")
		return
	}

	execInDb.Password = hashedPassword
	execInDb.PasswordResetCode = sql.NullString{Valid: false, String: ""}
	execInDb.PasswordChangedAt = sql.NullString{Valid: true, String: time.Now().String()}

	updatesMap := map[string]any{
		"password":          execInDb.Password,
		"passwordResetCode": execInDb.PasswordResetCode,
		"passwordChangedAt": execInDb.PasswordChangedAt,
	}

	_, err = sqlconnect.PatchExecById(execInDb.Id, updatesMap)
	if err != nil {
		http.Error(w, "Error updating password in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
