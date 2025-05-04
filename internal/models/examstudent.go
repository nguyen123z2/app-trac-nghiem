package models

import (
	"exam-system/internal/db"
	"fmt"
)

// Cấu trúc sinh viên tham gia kỳ thi
type ExamStudent struct {
	ID            int
	ExamID        int
	StudentID     int
	GroupID       int    `json:"group_id"`
	StudentCode   string `json:"student_code"`
	GroupName     string `json:"group_name"`
	Status        string
	TestStartedAt string
	RemainingTime int
}

// Thêm sinh viên vào kỳ thi
func InsertExamStudent(examID, studentID, groupID int) error {
	database := db.GetDB()

	// Kiểm tra sinh viên đã có trong kỳ thi chưa
	var existingStudentID int
	err := database.QueryRow("SELECT student_id FROM exam_students WHERE exam_id = ? AND student_id = ?", examID, studentID).Scan(&existingStudentID)
	if err == nil {
		// Nếu sinh viên đã có, chỉ cần cập nhật trạng thái
		return fmt.Errorf("Sinh viên đã tham gia kỳ thi này rồi.")
	}

	// Thêm sinh viên vào kỳ thi với trạng thái "pending"
	_, err = database.Exec(
		"INSERT INTO exam_students (exam_id, student_id, group_id, status ) VALUES (?, ?, ?, ?)",
		examID, studentID, groupID, "pending",
	)
	return err
}

// Cập nhật trạng thái sinh viên trong kỳ thi (VD: từ "pending" thành "in_progress" hoặc "completed")
func UpdateExamStudentStatus(id int, status string) error {
	database := db.GetDB()
	_, err := database.Exec(
		"UPDATE exam_students SET status=? WHERE id=?",
		status, id,
	)
	return err
}

// Lấy danh sách sinh viên tham gia kỳ thi
func GetExamStudentsByExam(examID int) ([]ExamStudent, error) {
	database := db.GetDB()
	rows, err := database.Query(
		"SELECT es.id, es.exam_id, s.student_code, eg.group_name, es.status, es.group_id "+ // Thêm group_id vào SELECT
			"FROM exam_students es "+
			"JOIN students s ON es.student_id = s.id "+
			"JOIN exam_groups eg ON es.group_id = eg.id "+
			"WHERE es.exam_id = ? "+
			"ORDER BY es.id ASC",
		examID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []ExamStudent
	for rows.Next() {
		var es ExamStudent
		// Lưu ý: Thêm trường group_id vào Scan
		if err := rows.Scan(&es.ID, &es.ExamID, &es.StudentCode, &es.GroupName, &es.Status, &es.GroupID); err != nil {
			return nil, err
		}
		students = append(students, es)
	}
	return students, nil
}

// Xóa sinh viên khỏi kỳ thi
func DeleteExamStudent(id int) error {
	database := db.GetDB()
	_, err := database.Exec("DELETE FROM exam_students WHERE id=?", id)
	return err
}
