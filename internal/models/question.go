package models

import (
	"encoding/json"
	"exam-system/internal/db"
)

// Question đại diện cho một câu hỏi
type Question struct {
	ID            int    `json:"id"`
	TeacherID     int    `json:"teacher_id"`
	Content       string `json:"content"`
	Difficulty    string `json:"difficulty"`
	QuestionSetID int    `json:"question_set_id"` // để UI biết quay về bộ nào
}

// QuestionAnswer đại diện cho câu trả lời của một câu hỏi
type QuestionAnswer struct {
	ID            int      `json:"id"`
	QuestionID    int      `json:"question_id"`
	Answers       []string `json:"answers"`
	CorrectAnswer string   `json:"correct_answer"`
}

// InsertQuestion thêm mới câu hỏi và liên kết vào question_set_questions
func InsertQuestion(content string, answers []string, correctAnswer, difficulty string, questionSetID int, teacherID int) (int, error) {
	database := db.GetDB()

	// 1) Thêm câu hỏi vào bảng questions
	res, err := database.Exec(
		"INSERT INTO questions (teacher_id, content, difficulty) VALUES (?, ?, ?)",
		teacherID, content, difficulty,
	)
	if err != nil {
		return 0, err
	}
	qID64, _ := res.LastInsertId()
	qID := int(qID64)

	// 2) Thêm vào bảng liên kết question_set_questions
	_, err = database.Exec(
		"INSERT INTO question_set_questions (question_set_id, question_id) VALUES (?, ?)",
		questionSetID, qID,
	)
	if err != nil {
		return 0, err
	}

	return qID, nil
}

// InsertQuestionAnswer thêm bộ đáp án cho câu hỏi
func InsertQuestionAnswer(questionID int, answers []string, correctAnswer string) error {
	database := db.GetDB()

	answersJSON, err := json.Marshal(answers)
	if err != nil {
		return err
	}
	_, err = database.Exec(
		"INSERT INTO question_answers (question_id, answers, correct_answer) VALUES (?, ?, ?)",
		questionID, answersJSON, correctAnswer,
	)
	return err
}

// GetQuestionAnswer lấy đúng/đủ 4 phương án
func GetQuestionAnswer(questionID int) (QuestionAnswer, error) {
	var qa QuestionAnswer
	var answersJSON string
	database := db.GetDB()

	err := database.QueryRow(`
		SELECT id, question_id, answers, correct_answer
		FROM question_answers
		WHERE question_id = ?
	`, questionID).Scan(&qa.ID, &qa.QuestionID, &answersJSON, &qa.CorrectAnswer)
	if err != nil {
		return qa, err
	}

	if err := json.Unmarshal([]byte(answersJSON), &qa.Answers); err != nil {
		return qa, err
	}
	for len(qa.Answers) < 4 {
		qa.Answers = append(qa.Answers, "")
	}
	return qa, nil
}

// UpdateQuestion cập nhật nội dung và độ khó
func UpdateQuestion(id int, content, difficulty string) error {
	database := db.GetDB()
	_, err := database.Exec(
		"UPDATE questions SET content=?, difficulty=? WHERE id=?",
		content, difficulty, id,
	)
	return err
}

// UpdateQuestionAnswer cập nhật đáp án
func UpdateQuestionAnswer(questionID int, answers []string, correctAnswer string) error {
	database := db.GetDB()
	answersJSON, err := json.Marshal(answers)
	if err != nil {
		return err
	}
	_, err = database.Exec(
		"UPDATE question_answers SET answers=?, correct_answer=? WHERE question_id=?",
		answersJSON, correctAnswer, questionID,
	)
	return err
}

// DeleteQuestion xóa câu hỏi và xóa liên kết trong question_set_questions
func DeleteQuestion(id int) error {
	database := db.GetDB()
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	// 1) xóa trong question_set_questions
	if _, err := tx.Exec("DELETE FROM question_set_questions WHERE question_id=?", id); err != nil {
		tx.Rollback()
		return err
	}
	// 2) xóa đáp án
	if _, err := tx.Exec("DELETE FROM question_answers WHERE question_id=?", id); err != nil {
		tx.Rollback()
		return err
	}
	// 3) xóa câu hỏi
	if _, err := tx.Exec("DELETE FROM questions WHERE id=?", id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// GetQuestions lấy toàn bộ câu hỏi của một bộ thông qua bảng trung gian
func GetQuestions(questionSetID int) ([]Question, error) {
	database := db.GetDB()
	rows, err := database.Query(`
		SELECT q.id, q.teacher_id, q.content, q.difficulty
		FROM question_set_questions qs
		JOIN questions q ON qs.question_id = q.id
		WHERE qs.question_set_id = ?
		ORDER BY q.id ASC
	`, questionSetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Question
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.ID, &q.TeacherID, &q.Content, &q.Difficulty); err != nil {
			return nil, err
		}
		q.QuestionSetID = questionSetID
		list = append(list, q)
	}
	return list, nil
}
func GetQuestionSetByID(id int) (QuestionSet, error) {
	var qs QuestionSet
	err := db.GetDB().QueryRow(
		"SELECT id, teacher_id, name FROM question_sets WHERE id = ?", id,
	).Scan(&qs.ID, &qs.TeacherID, &qs.Name)
	return qs, err
}
