// internal/ui/exam_students_view.go
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

	"github.com/xuri/excelize/v2"
)

// ShowExamStudentsAdmin dành cho Admin, quay lại ShowExamListAdmin
func ShowExamStudentsAdmin(w fyne.Window, examID int) {
	ShowExamStudentsGeneric(w, examID, ShowExamListAdmin)
}

// ShowExamStudentsTeacher dành cho Teacher, quay lại ShowExamListTeacher
func ShowExamStudentsTeacher(w fyne.Window, examID int) {
	if _, err := auth.GetCurrentUserID(); err != nil {
		dialog.ShowInformation("Lỗi", "Không xác định được giáo viên!", w)
		return
	}
	ShowExamStudentsGeneric(w, examID, ShowExamListTeacher)
}

// ShowExamStudentsGeneric hiển thị danh sách sinh viên tham gia kỳ thi,
// và sử dụng back(w) để quay về đúng màn hình (Admin hoặc Teacher).
func ShowExamStudentsGeneric(w fyne.Window, examID int, back func(fyne.Window)) {
	topBar := container.NewHBox(
		widget.NewLabelWithStyle("📋 Danh sách Sinh Viên Tham Gia Kỳ Thi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewButton("⬅ Quay lại", func() { back(w) }),
	)

	addBtn := widget.NewButton("➕ Thêm Sinh Viên", func() {
		ShowAddExamStudentForm(w, examID, 2, back)
	})

	content := container.NewVBox()
	updateExamStudentList(w, examID, content, back)

	w.SetContent(container.NewVBox(
		topBar,
		addBtn,
		content,
	))
}

// updateExamStudentList vẽ lại bảng theo danh sách lấy từ models.GetExamStudentsByExam
func updateExamStudentList(w fyne.Window, examID int, content *fyne.Container, back func(fyne.Window)) {
	students, err := models.GetExamStudentsByExam(examID)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	content.Objects = nil
	// Tiêu đề
	headers := container.NewHBox(
		widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Mã SV", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Nhóm", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Trạng thái", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Hành động", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	content.Add(headers)
	content.Add(widget.NewSeparator())

	// Nội dung
	list := container.NewVBox()
	for _, es := range students {
		e := es // capture

		editBtn := widget.NewButton("✏️", func() {
			ShowUpdateExamStudentStatusForm(w, examID, e.ID, e.Status, back)
		})
		deleteBtn := widget.NewButton("🗑️", func() {
			dialog.ShowConfirm("Xác nhận xóa", "Bạn chắc chắn muốn xóa sinh viên này?", func(ok bool) {
				if ok {
					models.DeleteExamStudent(e.ID)
					ShowExamStudentsGeneric(w, examID, back)
				}
			}, w)
		})
		exportBtn := widget.NewButton("📄 Xuất PDF", func() {
			go func() {
				file, err := models.ExportPDF(e.ID)
				if err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("Hoàn tất", "File đã lưu: "+file, w)
				}
			}()
		})

		row := container.NewHBox(
			widget.NewLabel(strconv.Itoa(e.ID)),
			layout.NewSpacer(),
			widget.NewLabel(e.StudentCode),
			layout.NewSpacer(),
			widget.NewLabel(e.GroupName),
			layout.NewSpacer(),
			widget.NewLabel(e.Status),
			layout.NewSpacer(),
			container.NewHBox(editBtn, deleteBtn, exportBtn),
		)
		list.Add(row)
		list.Add(widget.NewSeparator())
	}

	scroll := container.NewVScroll(list)
	scroll.SetMinSize(fyne.NewSize(700, 400))
	content.Add(scroll)
	content.Refresh()
}

// ShowAddExamStudentForm thêm sinh viên thủ công hoặc từ Excel,
// rồi quay lại bằng back(w).
func ShowAddExamStudentForm(w fyne.Window, examID, groupID int, back func(fyne.Window)) {
	filePicker := widget.NewButton("📄 Chọn File Excel", func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			defer uc.Close()
			if err := processExcelFile(uc.URI().Path(), examID, groupID); err != nil {
				dialog.ShowError(err, w)
			} else {
				ShowExamStudentsGeneric(w, examID, back)
			}
		}, w)
	})

	codeEntry := widget.NewEntry()
	codeEntry.SetPlaceHolder("Mã Sinh Viên")
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Họ và Tên")
	groupEntry := widget.NewEntry()
	groupEntry.SetPlaceHolder("ID Nhóm")
	groupEntry.SetText(strconv.Itoa(groupID))

	saveBtn := widget.NewButton("💾 Thêm Sinh Viên", func() {
		code, name := codeEntry.Text, nameEntry.Text
		gid, err := strconv.Atoi(groupEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("ID Nhóm không hợp lệ"), w)
			return
		}
		stu, err := models.GetStudentByCode(code)
		if err != nil {
			if err := models.AddStudent(code, name, "", "", ""); err != nil {
				dialog.ShowError(err, w)
				return
			}
			stu, err = models.GetStudentByCode(code)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
		}
		if err := models.InsertExamStudent(examID, stu.ID, gid); err != nil {
			dialog.ShowError(err, w)
		} else {
			ShowExamStudentsGeneric(w, examID, back)
		}
	})

	backBtn := widget.NewButton("⬅ Quay lại", func() {
		ShowExamStudentsGeneric(w, examID, back)
	})

	form := container.NewVBox(
		widget.NewLabel("Thêm Sinh Viên Tham Gia Kỳ Thi"),
		widget.NewLabel("Mã Sinh Viên"), codeEntry,
		widget.NewLabel("Họ và Tên"), nameEntry,
		widget.NewLabel("ID Nhóm"), groupEntry,
		saveBtn, backBtn, filePicker,
	)
	w.SetContent(container.NewCenter(form))
}

// ShowUpdateExamStudentStatusForm cập nhật trạng thái rồi back(w)
func ShowUpdateExamStudentStatusForm(w fyne.Window, examID, esID int, currentStatus string, back func(fyne.Window)) {
	options := []string{"pending", "in_progress", "completed"}
	selectStatus := widget.NewSelect(options, func(sel string) {
		if err := models.UpdateExamStudentStatus(esID, sel); err != nil {
			dialog.ShowError(err, w)
			return
		}
		ShowExamStudentsGeneric(w, examID, back)
	})
	selectStatus.SetSelected(currentStatus)

	backBtn := widget.NewButton("⬅ Quay lại", func() {
		ShowExamStudentsGeneric(w, examID, back)
	})

	w.SetContent(container.NewCenter(container.NewVBox(
		widget.NewLabel("Cập nhật trạng thái"), selectStatus, backBtn,
	)))
}

// processExcelFile giống cũ
func processExcelFile(path string, examID, groupID int) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}
	for i, row := range rows {
		if i == 0 || len(row) < 7 {
			continue
		}
		code := row[1]
		name := row[2] + " " + row[3]
		gender, dob, place := row[4], row[5], row[6]
		stu, err := models.GetStudentByCode(code)
		if err != nil {
			if err := models.AddStudent(code, name, gender, dob, place); err != nil {
				return err
			}
			stu, err = models.GetStudentByCode(code)
			if err != nil {
				return err
			}
		}
		if err := models.InsertExamStudent(examID, stu.ID, groupID); err != nil {
			return err
		}
	}
	return nil
}
