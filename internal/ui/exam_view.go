// internal/ui/exam_view.go
package ui

import (
	"fmt"
	"strconv"

	"exam-system/internal/auth"
	"exam-system/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// updateExamList lấy danh sách kỳ thi. Nếu teacherID > 0 sẽ chỉ lấy kỳ thi của giáo viên đó,
// ngược lại lấy toàn bộ (dành cho Admin).
func updateExamList(w fyne.Window, content *fyne.Container, teacherID int) {
	var exams []models.Exam
	var err error

	if teacherID > 0 {
		exams, err = models.GetExamsByTeacher(teacherID)
	} else {
		exams, err = models.GetExams()
	}
	if err != nil {
		fmt.Println("❌ Lỗi khi lấy danh sách kỳ thi:", err)
		return
	}

	// Xóa nội dung cũ
	content.Objects = nil

	// Tiêu đề bảng
	headers := container.NewHBox(
		widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Tên Kỳ Thi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Giáo Viên", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Bộ Câu Hỏi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Thao Tác", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	content.Add(headers)
	content.Add(widget.NewSeparator())

	// Duyệt qua từng kỳ thi
	for _, ex := range exams {
		e := ex // capture cho closure

		// Tạo nút SV tuỳ theo role
		var svBtn *widget.Button
		if teacherID > 0 {
			// Teacher chỉ xem SV của kỳ thi mình
			svBtn = widget.NewButton("📋 SV", func() {
				ShowExamStudentsTeacher(w, e.ID)
			})
		} else {
			// Admin xem SV của mọi kỳ thi
			svBtn = widget.NewButton("📋 SV", func() {
				ShowExamStudentsAdmin(w, e.ID)
			})
		}

		row := container.NewHBox(
			widget.NewLabel(strconv.Itoa(e.ID)),
			layout.NewSpacer(),
			widget.NewLabel(e.Name),
			layout.NewSpacer(),
			widget.NewLabel(strconv.Itoa(e.TeacherID)),
			layout.NewSpacer(),
			widget.NewLabel(strconv.Itoa(e.QuestionSetID)),
			layout.NewSpacer(),
			widget.NewButton("✏️", func() { ShowExamForm(w, &e) }),
			widget.NewButton("🗑️", func() { deleteExam(w, e.ID) }),
			svBtn,
		)
		content.Add(row)
		content.Add(widget.NewSeparator())
	}

	content.Refresh()
}

// ShowExamListAdmin hiển thị toàn bộ kỳ thi cho Admin.
func ShowExamListAdmin(w fyne.Window) {
	content := container.NewVBox()
	updateExamList(w, content, 0)

	w.SetContent(container.NewVBox(
		widget.NewLabelWithStyle("📅 Danh sách Kỳ Thi (Admin)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("➕ Thêm Kỳ Thi", func() { ShowExamForm(w, nil) }),
		content,
		widget.NewButton("⬅ Quay lại", func() { ShowAdminManagement(w) }),
	))
}

// ShowExamListTeacher hiển thị kỳ thi của giáo viên đang đăng nhập.
func ShowExamListTeacher(w fyne.Window) {
	teacherID, err := auth.GetCurrentUserID()
	if err != nil {
		dialog.ShowInformation("Lỗi", "Không xác định được giáo viên!", w)
		return
	}

	content := container.NewVBox()
	updateExamList(w, content, teacherID)

	w.SetContent(container.NewVBox(
		widget.NewLabelWithStyle("📅 Danh sách Kỳ Thi (Teacher)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		content,
		widget.NewButton("⬅ Quay lại", func() { ShowTeacherDashboard(w) }),
	))
}

// ShowExamForm hiển thị form thêm/sửa kỳ thi chung cho Admin và Teacher.
func ShowExamForm(w fyne.Window, exam *models.Exam) {
	var isEdit bool
	if exam != nil {
		isEdit = true
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nhập tên kỳ thi")
	teacherEntry := widget.NewEntry()
	teacherEntry.SetPlaceHolder("Nhập ID giáo viên")
	qsEntry := widget.NewEntry()
	qsEntry.SetPlaceHolder("Nhập ID bộ câu hỏi")

	if isEdit {
		nameEntry.SetText(exam.Name)
		teacherEntry.SetText(strconv.Itoa(exam.TeacherID))
		qsEntry.SetText(strconv.Itoa(exam.QuestionSetID))
	}

	saveBtn := widget.NewButton("💾 Lưu", func() {
		name := nameEntry.Text
		tid, _ := strconv.Atoi(teacherEntry.Text)
		qsid, _ := strconv.Atoi(qsEntry.Text)

		if isEdit {
			models.UpdateExam(exam.ID, tid, name, qsid)
		} else {
			models.InsertExam(tid, name, qsid)
		}

		// Quay lại tuỳ role
		if role, _ := auth.GetUserRole(""); role == "teacher" {
			ShowExamListTeacher(w)
		} else {
			ShowExamListAdmin(w)
		}
	})

	backBtn := widget.NewButton("⬅ Quay lại", func() {
		if role, _ := auth.GetUserRole(""); role == "teacher" {
			ShowExamListTeacher(w)
		} else {
			ShowExamListAdmin(w)
		}
	})

	form := container.NewVBox(
		widget.NewLabel("Thông tin Kỳ Thi"),
		widget.NewLabel("Tên kỳ thi"), nameEntry,
		widget.NewLabel("ID giáo viên"), teacherEntry,
		widget.NewLabel("ID bộ câu hỏi"), qsEntry,
		saveBtn,
		backBtn,
	)
	w.SetContent(container.NewCenter(form))
}

func deleteExam(w fyne.Window, id int) {
	dialog.ShowConfirm("Xác nhận xóa", "Bạn có chắc chắn muốn xóa kỳ thi này?", func(ok bool) {
		if ok {
			models.DeleteExam(id)
			// refresh lại theo role
			if role, _ := auth.GetUserRole(""); role == "teacher" {
				ShowExamListTeacher(w)
			} else {
				ShowExamListAdmin(w)
			}
		}
	}, w)
}
