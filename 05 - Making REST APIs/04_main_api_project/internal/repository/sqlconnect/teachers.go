package sqlconnect

import (
	"restapi/internal/models"
)

func GetTeachers(query string, args []any, teacherList []models.Teacher) ([]models.Teacher, error) {
	db := ConnectDb()
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, err
		}
		teacherList = append(teacherList, teacher)
	}
	return teacherList, nil
}
