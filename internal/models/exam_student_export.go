// internal/models/exam_student_export.go
package models

import (
	"fmt"

	"exam-system/internal/db"

	"github.com/phpdave11/gofpdf"
)

type StudentResult struct {
	Question      string
	Selected      string
	CorrectAnswer string
	IsCorrect     bool
	AnsweredAt    string
}

type ExportData struct {
	StudentName string
	StudentCode string
	ExamName    string
	GroupName   string
	Results     []StudentResult
}

// GetExportData giữ nguyên như trước
func GetExportData(examStudentID int) (*ExportData, error) {
	database := db.GetDB()
	var ed ExportData

	err := database.QueryRow(`
		SELECT s.name, s.student_code, e.name, eg.group_name
		FROM exam_students es
		JOIN students s ON es.student_id=s.id
		JOIN exam_groups eg ON es.group_id=eg.id
		JOIN exams e ON eg.exam_id=e.id
		WHERE es.id = ?
	`, examStudentID).Scan(
		&ed.StudentName,
		&ed.StudentCode,
		&ed.ExamName,
		&ed.GroupName,
	)
	if err != nil {
		return nil, err
	}

	rows, err := database.Query(`
		SELECT q.content, esa.selected_answer, qa.correct_answer, esa.is_correct,
		       DATE_FORMAT(esa.answered_at, '%Y-%m-%d %H:%i')
		FROM exam_student_answers esa
		JOIN questions q ON esa.question_id = q.id
		JOIN question_answers qa ON qa.question_id = q.id
		WHERE esa.exam_student_id = ?
		ORDER BY esa.answered_at ASC
	`, examStudentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r StudentResult
		var ok int
		if err := rows.Scan(
			&r.Question,
			&r.Selected,
			&r.CorrectAnswer,
			&ok,
			&r.AnsweredAt,
		); err != nil {
			return nil, err
		}
		r.IsCorrect = ok == 1
		ed.Results = append(ed.Results, r)
	}

	return &ed, nil
}

// ExportPDF tạo file PDF và trả về tên file
func ExportPDF(examStudentID int) (string, error) {
	data, err := GetExportData(examStudentID)
	if err != nil {
		return "", err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	// Đăng ký cả Regular và Bold
	pdf.AddUTF8Font("DejaVu", "", "internal/assets/fonts/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVu", "B", "internal/assets/fonts/DejaVuSans-Bold.ttf")

	// Title font
	pdf.SetFont("DejaVu", "B", 16)
	title := fmt.Sprintf("%s - %s", data.StudentCode, data.ExamName)
	pdf.SetTitle(title, false)
	pdf.AddPage()

	// Header
	pdf.Cell(0, 10, "KẾT QUẢ BÀI THI")
	pdf.Ln(12)

	// Nội dung chung
	pdf.SetFont("DejaVu", "", 12)
	pdf.Cell(40, 8, "Sinh viên:")
	pdf.Cell(0, 8, fmt.Sprintf("%s (%s)", data.StudentName, data.StudentCode))
	pdf.Ln(6)
	pdf.Cell(40, 8, "Kỳ thi:")
	pdf.Cell(0, 8, data.ExamName)
	pdf.Ln(6)
	pdf.Cell(40, 8, "Nhóm thi:")
	pdf.Cell(0, 8, data.GroupName)
	pdf.Ln(10)

	// Table header
	headers := []string{"STT", "Câu hỏi", "Đáp án chọn", "Đúng/Sai", "Thời gian"}
	widths := []float64{10, 80, 40, 20, 40}
	pdf.SetFont("DejaVu", "B", 12)
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Table body
	pdf.SetFont("DejaVu", "", 11)
	for i, r := range data.Results {
		pdf.CellFormat(widths[0], 8, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[1], 8, r.Question, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[2], 8, r.Selected, "1", 0, "", false, 0, "")
		status := "Sai"
		if r.IsCorrect {
			status = "Đúng"
		}
		pdf.CellFormat(widths[3], 8, status, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 8, r.AnsweredAt, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	filename := fmt.Sprintf("export_%s_%s.pdf", data.StudentCode, data.ExamName)
	if err := pdf.OutputFileAndClose(filename); err != nil {
		return "", err
	}
	return filename, nil
}
