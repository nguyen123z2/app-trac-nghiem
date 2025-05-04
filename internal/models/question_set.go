package models

import (
	"exam-system/internal/db"
)

// Cấu trúc bộ câu hỏi
type QuestionSet struct {
	ID        int
	TeacherID int
	Name      string
}

// Thêm bộ câu hỏi vào database
func InsertQuestionSet(teacherID int, name string) error {
	database := db.GetDB()
	_, err := database.Exec(
		"INSERT INTO question_sets (teacher_id, name) VALUES (?, ?)",
		teacherID, name,
	)
	return err
}

// Lấy danh sách bộ câu hỏi từ database
func GetQuestionSets() ([]QuestionSet, error) {
	database := db.GetDB()
	rows, err := database.Query("SELECT id, teacher_id, name FROM question_sets ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questionSets []QuestionSet
	for rows.Next() {
		var q QuestionSet
		if err := rows.Scan(&q.ID, &q.TeacherID, &q.Name); err != nil {
			return nil, err
		}
		questionSets = append(questionSets, q)
	}
	return questionSets, nil
}

// Cập nhật bộ câu hỏi
func UpdateQuestionSet(id int, teacherID int, name string) error {
	database := db.GetDB()
	_, err := database.Exec(
		"UPDATE question_sets SET teacher_id=?, name=? WHERE id=?",
		teacherID, name, id,
	)
	return err
}

// Xóa bộ câu hỏi
func DeleteQuestionSet(id int) error {
	database := db.GetDB()
	_, err := database.Exec("DELETE FROM question_sets WHERE id=?", id)
	return err
}
