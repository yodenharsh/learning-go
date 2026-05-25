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

func GetExecs(dbParams map[string]string, r *http.Request) ([]models.Exec, error) {
	db := ConnectDb()
	defer db.Close()
	execs := make([]models.Exec, 0)
	query := "SELECT id, first_name, last_name, email, username, password, inactive_status, role, password_reset_code, password_code_expires_at, password_changed_at, user_created_at FROM execs WHERE 1=1"
	query, args := utils.BuildQueryWithFilters(r, query, dbParams)

	query = utils.BuildQueryWithSorting(r, query, dbParams)

	fmt.Println(query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error querying DB")
	}
	defer rows.Close()

	for rows.Next() {
		var exec models.Exec
		err := rows.Scan(&exec.Id, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.InactiveStatus, &exec.Role, &exec.PasswordResetCode, &exec.PasswordCodeExpiresAt, &exec.PasswordChangedAt, &exec.UserCreatedAt)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Scan row failed")
		}
		execs = append(execs, exec)
	}
	return execs, nil
}

func GetExecById(id int) (models.Exec, error) {
	db := ConnectDb()
	defer db.Close()

	var exec models.Exec
	err := db.QueryRow("SELECT id, first_name, last_name, email, username, password, inactive_status, role, password_reset_code, password_code_expires_at, password_changed_at, user_created_at FROM execs WHERE id = ?", id).Scan(&exec.Id, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.InactiveStatus, &exec.Role, &exec.PasswordResetCode, &exec.PasswordCodeExpiresAt, &exec.PasswordChangedAt, &exec.UserCreatedAt)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Querying execs failed")
	}
	return exec, nil
}

func GetExecByUsername(username string) (models.Exec, error) {
	db := ConnectDb()
	defer db.Close()

	var exec models.Exec
	err := db.QueryRow("SELECT id, first_name, last_name, email, username, password, inactive_status, role, password_reset_code, password_code_expires_at, password_changed_at, user_created_at FROM execs WHERE username = ?", username).Scan(&exec.Id, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.InactiveStatus, &exec.Role, &exec.PasswordResetCode, &exec.PasswordCodeExpiresAt, &exec.PasswordChangedAt, &exec.UserCreatedAt)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Querying execs failed")
	}
	return exec, nil
}

func GetExecByEmail(email string) (models.Exec, error) {
	db := ConnectDb()
	defer db.Close()

	var exec models.Exec
	err := db.QueryRow("SELECT id, first_name, last_name, email, username, password, inactive_status, role, password_reset_code, password_code_expires_at, password_changed_at, user_created_at FROM execs WHERE email = ?", email).Scan(&exec.Id, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.InactiveStatus, &exec.Role, &exec.PasswordResetCode, &exec.PasswordCodeExpiresAt, &exec.PasswordChangedAt, &exec.UserCreatedAt)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Querying execs failed")
	}
	return exec, nil
}

func AddExec(newExecs []models.Exec) ([]models.Exec, error) {
	db := ConnectDb()
	defer db.Close()

	fmt.Println(utils.GenerateInsertQuery("execs", models.Exec{}))
	stmt, err := db.Prepare(utils.GenerateInsertQuery("execs", models.Exec{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Preparing insert stmt failde")
	}
	defer stmt.Close()

	addedExecs := make([]models.Exec, len(newExecs))
	for i, newExec := range newExecs {
		values := utils.GetStructValues(newExec)
		res, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to add exec")
		}

		res.LastInsertId()
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to check last inserted ID")
		}

		newExec.Id = int(lastId)
		addedExecs[i] = newExec
	}
	return addedExecs, nil
}

func PatchExecById(id int, updates map[string]any) (models.Exec, error) {
	db := ConnectDb()
	defer db.Close()

	var existingExec models.Exec
	err := db.QueryRow("SELECT id, first_name, last_name, email, username, password, inactive_status, role, password_reset_code, password_code_expires_at, password_changed_at, user_created_at FROM execs WHERE id = ?", id).Scan(&existingExec.Id, &existingExec.FirstName, &existingExec.LastName, &existingExec.Email, &existingExec.Username, &existingExec.Password, &existingExec.InactiveStatus, &existingExec.Role, &existingExec.PasswordResetCode, &existingExec.PasswordCodeExpiresAt, &existingExec.PasswordChangedAt, &existingExec.UserCreatedAt)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Failed to do existing check")
	}

	execVal := reflect.ValueOf(&existingExec).Elem()
	execType := execVal.Type()

	for k, v := range updates {
		for i := 0; i < execVal.NumField(); i++ {
			field := execType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if execVal.Field(i).CanSet() {
					fieldVal := execVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(execVal.Field(i).Type()))
				}
			}
		}
	}

	query := "UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, password = ?, inactive_status = ?, role = ?, password_reset_code = ?, password_code_expires_at = ?, password_changed_at = ?, user_created_at = ? WHERE id = ?"
	_, err = db.Exec(query, &existingExec.FirstName, &existingExec.LastName, &existingExec.Email, &existingExec.Username, &existingExec.Password, &existingExec.InactiveStatus, &existingExec.Role, &existingExec.PasswordResetCode, &existingExec.PasswordCodeExpiresAt, &existingExec.PasswordChangedAt, &existingExec.UserCreatedAt, id)
	return existingExec, err
}

func PatchExecs(updates []map[string]any) error {
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
			return utils.ErrorHandler(errors.New("Invalid or missing exec ID"), "Invalid or missing exec ID")
		}

		id, err := strconv.Atoi(stringId)
		if err != nil {
			tx.Rollback()
			return err
		}

		var execFromDb models.Exec

		err = db.QueryRow("SELECT id, first_name, last_name, email, username, password, inactive_status, role, password_reset_code, password_code_expires_at, password_changed_at, user_created_at FROM execs WHERE id = ?", id).
			Scan(&execFromDb.Id, &execFromDb.FirstName, &execFromDb.LastName, &execFromDb.Email, &execFromDb.Username, &execFromDb.Password, &execFromDb.InactiveStatus, &execFromDb.Role, &execFromDb.PasswordResetCode, &execFromDb.PasswordCodeExpiresAt, &execFromDb.PasswordChangedAt, &execFromDb.UserCreatedAt)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Failed to get existing execs")
		}
		execVal := reflect.ValueOf(&execFromDb).Elem()
		execType := execVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}

			for i := 0; i < execVal.NumField(); i++ {
				field := execType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					if execVal.Field(i).CanSet() {
						fieldVal := execVal.Field(i)
						if fieldVal.CanSet() {
							val := reflect.ValueOf(v)
							if val.Type().ConvertibleTo(fieldVal.Type()) {
								fieldVal.Set(val.Convert(fieldVal.Type()))
							} else {
								tx.Rollback()
								return utils.ErrorHandler(errors.New("Error updating exec in database"), "Error updating exec in database")
							}
						}
					}
				}
			}
		}
		_, err = tx.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, password = ?, inactive_status = ?, role = ?, password_reset_code = ?, password_code_expires_at = ?, password_changed_at = ?, user_created_at = ? WHERE id = ?", execFromDb.FirstName, execFromDb.LastName, execFromDb.Email, execFromDb.Username, execFromDb.Password, execFromDb.InactiveStatus, execFromDb.Role, execFromDb.PasswordResetCode, execFromDb.PasswordCodeExpiresAt, execFromDb.PasswordChangedAt, execFromDb.UserCreatedAt, id)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error updating exec in database")
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "Error updating transaction")
	}
	return nil
}

func DeleteExecs(ids []int) ([]int, error) {
	db := ConnectDb()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	deletedIds := []int{}
	for _, id := range ids {
		execStmt, err := tx.Prepare("DELETE FROM execs WHERE id = ?")
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

func DeleteExecById(id int) error {
	db := ConnectDb()
	defer db.Close()

	query := "DELETE FROM execs WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to delete exec")
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

func UpdatePasswordResetCode(id int, resetCode string, expiresAt string) error {
	db := ConnectDb()
	defer db.Close()

	query := "UPDATE execs SET password_reset_code = ?, password_code_expires_at = ? WHERE id = ?"
	result, err := db.Exec(query, resetCode, expiresAt, id)

	if err != nil {
		return utils.ErrorHandler(err, "Error when adding password reset code")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Error when adding password reset code")
	} else if rowsAffected == 0 {
		return utils.ErrorHandler(err, "Exec with id doesn't exist")
	}

	return nil
}

func GetExecByResetCode(resetCode string) (models.Exec, error) {
	db := ConnectDb()
	defer db.Close()

	var exec models.Exec
	err := db.QueryRow("SELECT id, first_name, last_name, email, username, password, inactive_status, role, password_reset_code, password_code_expires_at, password_changed_at, user_created_at FROM execs WHERE password_reset_code = ?", resetCode).Scan(&exec.Id, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.InactiveStatus, &exec.Role, &exec.PasswordResetCode, &exec.PasswordCodeExpiresAt, &exec.PasswordChangedAt, &exec.UserCreatedAt)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Querying execs failed")
	}
	return exec, nil
}
