package models

import (
	"exam-system/internal/db"
)

// Cấu trúc nhóm kỳ thi
type ExamGroup struct {
	ID        int
	ExamID    int
	GroupName string
	ExamName  string // Thêm trường ExamName để lưu tên kỳ thi
}

// Cấu trúc kỳ thi (đổi tên thành `ExamDetail` để tránh trùng lặp)
type ExamDetail struct {
	ID   int
	Name string
}

// Thêm nhóm kỳ thi vào database
func InsertExamGroup(examID int, groupName string) error {
	database := db.GetDB()
	_, err := database.Exec(
		"INSERT INTO exam_groups (exam_id, group_name) VALUES (?, ?)",
		examID, groupName,
	)
	return err
}

// Lấy danh sách nhóm kỳ thi từ database và bao gồm tên kỳ thi
func GetExamGroups() ([]ExamGroup, error) {
	database := db.GetDB()
	rows, err := database.Query(`
		SELECT eg.id, eg.exam_id, eg.group_name, e.name 
		FROM exam_groups eg 
		JOIN exams e ON eg.exam_id = e.id
		ORDER BY eg.id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []ExamGroup
	for rows.Next() {
		var g ExamGroup
		if err := rows.Scan(&g.ID, &g.ExamID, &g.GroupName, &g.ExamName); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// Lấy tất cả kỳ thi từ database (đổi tên thành `GetExamDetails`)
func GetExamDetails() ([]ExamDetail, error) {
	database := db.GetDB()
	rows, err := database.Query("SELECT id, name FROM exams ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []ExamDetail
	for rows.Next() {
		var e ExamDetail
		if err := rows.Scan(&e.ID, &e.Name); err != nil {
			return nil, err
		}
		exams = append(exams, e)
	}
	return exams, nil
}

// Cập nhật nhóm kỳ thi
func UpdateExamGroup(id int, examID int, groupName string) error {
	database := db.GetDB()
	_, err := database.Exec(
		"UPDATE exam_groups SET exam_id=?, group_name=? WHERE id=?",
		examID, groupName, id,
	)
	return err
}

// Xóa nhóm kỳ thi
func DeleteExamGroup(id int) error {
	database := db.GetDB()
	_, err := database.Exec("DELETE FROM exam_groups WHERE id=?", id)
	return err
}
func GetExamGroupsByTeacher(teacherID int) ([]ExamGroup, error) {
	database := db.GetDB()
	rows, err := database.Query(`
		SELECT eg.id, eg.exam_id, eg.group_name, e.name
		FROM exam_groups eg
		JOIN exams e ON eg.exam_id = e.id
		WHERE e.teacher_id = ?
		ORDER BY eg.id ASC
	`, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []ExamGroup
	for rows.Next() {
		var g ExamGroup
		if err := rows.Scan(&g.ID, &g.ExamID, &g.GroupName, &g.ExamName); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
