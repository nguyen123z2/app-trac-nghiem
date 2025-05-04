package models

import (
	"exam-system/internal/db"
)

// Cấu trúc kỳ thi
type Exam struct {
	ID            int    `json:"id"`
	TeacherID     int    `json:"teacher_id"`
	Name          string `json:"name"`
	QuestionSetID int    `json:"question_set_id"`
}

// Thêm kỳ thi vào database
func InsertExam(teacherID int, name string, questionSetID int) error {
	database := db.GetDB()
	_, err := database.Exec(
		"INSERT INTO exams (teacher_id, name, question_set_id) VALUES (?, ?, ?)",
		teacherID, name, questionSetID,
	)
	return err
}

// Lấy danh sách kỳ thi từ database
func GetExams() ([]Exam, error) {
	database := db.GetDB()
	rows, err := database.Query("SELECT id, teacher_id, name, question_set_id FROM exams ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []Exam
	for rows.Next() {
		var e Exam
		if err := rows.Scan(&e.ID, &e.TeacherID, &e.Name, &e.QuestionSetID); err != nil {
			return nil, err
		}
		exams = append(exams, e)
	}
	return exams, nil
}

// Cập nhật kỳ thi
func UpdateExam(id int, teacherID int, name string, questionSetID int) error {
	database := db.GetDB()
	_, err := database.Exec(
		"UPDATE exams SET teacher_id=?, name=?, question_set_id=? WHERE id=?",
		teacherID, name, questionSetID, id,
	)
	return err
}

// Xóa kỳ thi
func DeleteExam(id int) error {
	database := db.GetDB()
	_, err := database.Exec("DELETE FROM exams WHERE id=?", id)
	return err
}

// Lấy danh sách kỳ thi của giáo viên
func GetExamsByTeacher(teacherID int) ([]Exam, error) {
	database := db.GetDB()
	rows, err := database.Query(
		"SELECT id, teacher_id, name, question_set_id FROM exams WHERE teacher_id = ? ORDER BY id ASC",
		teacherID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []Exam
	for rows.Next() {
		var e Exam
		if err := rows.Scan(&e.ID, &e.TeacherID, &e.Name, &e.QuestionSetID); err != nil {
			return nil, err
		}
		exams = append(exams, e)
	}
	return exams, nil
}
