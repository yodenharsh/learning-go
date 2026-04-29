package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"slices"
	"strconv"
	"strings"
)

func TeachersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		postTeachersHandle(w, r)
	case http.MethodPut:
		updateTeachersHandle(w, r)
	case http.MethodPatch:
		patchTeachersHandle(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	db := sqlconnect.ConnectDb()
	defer db.Close()

	dbParams := map[string]string{
		"firstName": "first_name",
		"lastName":  "last_name",
		"email":     "email",
		"class":     "class",
		"subject":   "subject",
	}

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id := strings.TrimPrefix(path, "/")

	if id == "" {
		query, args := buildQueryWithFilters(r, dbParams)

		query = buildQueryWithSorting(r, query, dbParams)

		teacherList := make([]models.Teacher, 0)

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var teacher models.Teacher
			err := rows.Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
			if err != nil {
				http.Error(w, "Error scanning database result", http.StatusInternalServerError)
				return
			}
			teacherList = append(teacherList, teacher)
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
	} else {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var teacher models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

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
}

func postTeachersHandle(w http.ResponseWriter, r *http.Request) {

	db := sqlconnect.ConnectDb()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		http.Error(w, "Error preparing statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)

	if err != nil {
		http.Error(w, "error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Error inserting teacher into database", http.StatusInternalServerError)
			return
		}

		res.LastInsertId()
		lastId, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error retrieving last insert ID", http.StatusInternalServerError)
			return
		}

		newTeacher.Id = int(lastId)
		addedTeachers[i] = newTeacher
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

func updateTeachersHandle(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
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

	db := sqlconnect.ConnectDb()
	defer db.Close()

	type IdHolder struct {
		id int
	}

	var existingTeacherId IdHolder
	err = db.QueryRow("SELECT id FROM teachers WHERE id = ?", id).Scan(&existingTeacherId.id)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error checking teacher existence", http.StatusInternalServerError)
		return
	}

	query := "UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?"
	_, err = db.Exec(query, &updatedTeacher.FirstName, &updatedTeacher.LastName, &updatedTeacher.Email, &updatedTeacher.Class, &updatedTeacher.Subject, id)

	if err != nil {
		http.Error(w, "Error updating teacher in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

func patchTeachersHandle(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
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
