package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"strconv"
	"strings"
	"sync"
)

var (
	teachers = make(map[int]models.Teacher)
	nextId   = 1
	mutex    = &sync.Mutex{}
)

func init() {
	teachers[nextId] = models.Teacher{
		Id:        nextId,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "10A",
		Subject:   "Math",
	}
	nextId++
	teachers[nextId] = models.Teacher{
		Id:        nextId,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10B",
		Subject:   "Science",
	}
	nextId++
	teachers[nextId] = models.Teacher{
		Id:        nextId,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "4C",
		Subject:   "English",
	}
	nextId++
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		postTeachersHandle(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	db := sqlconnect.ConnectDb()
	defer db.Close()

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id := strings.TrimPrefix(path, "/")

	if id == "" {
		query, args := buildQueryWithFilters(r)

		teacherList := make([]models.Teacher, 0, len(teachers))

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

func buildQueryWithFilters(r *http.Request) (string, []any) {
	var query strings.Builder
	query.WriteString("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1")
	var args []any

	params := map[string]string{
		"firstName": "first_name",
		"lastName":  "last_name",
		"email":     "email",
		"class":     "class",
		"subject":   "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query.WriteString(" AND " + dbField + " = ?")
			args = append(args, value)
		}
	}
	return query.String(), args
}
