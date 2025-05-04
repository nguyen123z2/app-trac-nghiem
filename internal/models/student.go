package models

import (
	"database/sql"
	"exam-system/internal/db"
	"fmt"
)

// Student struct
type Student struct {
	ID           int
	StudentCode  string
	Name         string
	Gender       string // Giới tính
	DateOfBirth  string // Ngày sinh
	PlaceOfBirth string // Nơi sinh
}

// Lấy danh sách Student từ database theo student_code
func GetStudents(search string) ([]Student, error) {
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	// Cập nhật truy vấn để lấy thêm thông tin về giới tính, ngày sinh và nơi sinh
	query := "SELECT id, student_code, name, gender, date_of_birth, place_of_birth FROM students"
	var rows *sql.Rows
	var err error

	// Thêm điều kiện tìm kiếm theo mã sinh viên
	if search != "" {
		query += " WHERE student_code LIKE ?"
		rows, err = database.Query(query, "%"+search+"%")
	} else {
		rows, err = database.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		err := rows.Scan(&s.ID, &s.StudentCode, &s.Name, &s.Gender, &s.DateOfBirth, &s.PlaceOfBirth)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

// Thêm Student vào cơ sở dữ liệu
func AddStudent(studentCode, name, gender, dateOfBirth, placeOfBirth string) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	// Cập nhật truy vấn để thêm thông tin giới tính, ngày sinh và nơi sinh
	query := "INSERT INTO students (student_code, name, gender, date_of_birth, place_of_birth) VALUES (?, ?, ?, ?, ?)"
	_, err := database.Exec(query, studentCode, name, gender, dateOfBirth, placeOfBirth)
	return err
}

// Cập nhật Student
func UpdateStudent(id int, studentCode, name, gender, dateOfBirth, placeOfBirth string) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	// Cập nhật truy vấn để sửa thông tin sinh viên, bao gồm giới tính, ngày sinh và nơi sinh
	query := "UPDATE students SET student_code = ?, name = ?, gender = ?, date_of_birth = ?, place_of_birth = ? WHERE id = ?"
	_, err := database.Exec(query, studentCode, name, gender, dateOfBirth, placeOfBirth, id)
	return err
}

// Xóa Student
func DeleteStudent(id int) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	query := "DELETE FROM students WHERE id = ?"
	_, err := database.Exec(query, id)
	return err
}

// Thêm sinh viên vào kỳ thi
func InsertStudentToExam(examID int, studentID int) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	// Thêm sinh viên vào kỳ thi với trạng thái "pending"
	query := "INSERT INTO exam_students (exam_id, student_id, status) VALUES (?, ?, ?)"
	_, err := database.Exec(query, examID, studentID, "pending")
	return err
}
func GetStudentByCode(studentCode string) (Student, error) {
	database := db.GetDB()
	if database == nil {
		return Student{}, fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	var student Student
	query := "SELECT id, student_code, name, gender, date_of_birth, place_of_birth FROM students WHERE student_code = ?"
	err := database.QueryRow(query, studentCode).Scan(&student.ID, &student.StudentCode, &student.Name, &student.Gender, &student.DateOfBirth, &student.PlaceOfBirth)
	if err != nil {
		return student, err
	}
	return student, nil
}
