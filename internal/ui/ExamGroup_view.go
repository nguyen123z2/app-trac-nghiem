package ui

import (
	"fmt"

	"exam-system/internal/auth"
	"exam-system/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowExamGroupList(w fyne.Window) {
	role, _ := auth.GetUserRole("")
	var groups []models.ExamGroup
	var err error
	if role == "teacher" {
		tid, _ := auth.GetCurrentUserID()
		groups, err = models.GetExamGroupsByTeacher(tid)
	} else {
		groups, err = models.GetExamGroups()
	}
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	// 1) Xây topBar (chỉ 1 lần)
	title := widget.NewLabelWithStyle("📋 Danh sách Nhóm Kỳ Thi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	backBtn := widget.NewButton("⬅ Quay lại", func() {
		if role == "teacher" {
			ShowTeacherDashboard(w)
		} else {
			ShowAdminManagement(w)
		}
	})
	topBar := container.NewHBox(title, layout.NewSpacer(), backBtn)

	// 2) Nút thêm nhóm
	addBtn := widget.NewButton("➕ Thêm Nhóm", func() {
		ShowExamGroupForm(w, nil)
	})

	// 3) Container dữ liệu chỉ chứa headers + rows
	content := container.NewVBox()
	// Tiêu đề cột
	headers := container.NewHBox(
		widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Tên Nhóm", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Tên Kỳ Thi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Thao Tác", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	content.Add(headers)
	content.Add(widget.NewSeparator())

	for _, g := range groups {
		grp := g
		editBtn := widget.NewButton("✏️", func() { ShowExamGroupForm(w, &grp) })
		delBtn := widget.NewButton("🗑️", func() {
			dialog.ShowConfirm("Xác nhận xóa", "Bạn có chắc chắn muốn xóa nhóm này?", func(ok bool) {
				if ok {
					models.DeleteExamGroup(grp.ID)
					ShowExamGroupList(w)
				}
			}, w)
		})

		row := container.NewHBox(
			widget.NewLabel(fmt.Sprintf("%d", grp.ID)),
			layout.NewSpacer(),
			widget.NewLabel(grp.GroupName),
			layout.NewSpacer(),
			widget.NewLabel(grp.ExamName),
			layout.NewSpacer(),
			container.NewHBox(editBtn, delBtn),
		)
		content.Add(row)
		content.Add(widget.NewSeparator())
	}

	// 4) Cuối cùng SetContent chỉ cần topBar, addBtn, content
	w.SetContent(container.NewVBox(
		topBar,
		addBtn,
		content,
	))
}
func ShowExamGroupForm(w fyne.Window, group *models.ExamGroup) {
	role, _ := auth.GetUserRole("")

	// Lấy danh sách tất cả kỳ thi (Admin) hoặc chỉ kỳ thi của Teacher
	var exams []models.Exam
	var err error
	if role == "teacher" {
		teacherID, err2 := auth.GetCurrentUserID()
		if err2 != nil {
			dialog.ShowInformation("Lỗi", "Không xác định được giáo viên!", w)
			return
		}
		exams, err = models.GetExamsByTeacher(teacherID)
	} else {
		exams, err = models.GetExams()
	}
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	// chuẩn bị form
	var isEdit bool
	if group != nil {
		isEdit = true
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Tên nhóm kỳ thi")
	examNames := make([]string, len(exams))
	for i, ex := range exams {
		examNames[i] = ex.Name
	}

	// Nếu edit thì set trước
	if isEdit {
		nameEntry.SetText(group.GroupName)
	}

	examSelect := widget.NewSelect(examNames, func(sel string) {
		// tìm ID tương ứng
		for _, ex := range exams {
			if ex.Name == sel {
				group.ExamID = ex.ID
				break
			}
		}
	})
	if isEdit {
		// chọn sẵn
		for _, ex := range exams {
			if ex.ID == group.ExamID {
				examSelect.SetSelected(ex.Name)
				break
			}
		}
	}

	saveBtn := widget.NewButton("💾 Lưu", func() {
		if nameEntry.Text == "" || examSelect.Selected == "" {
			dialog.ShowInformation("Lỗi", "Vui lòng nhập đủ thông tin", w)
			return
		}
		if isEdit {
			models.UpdateExamGroup(group.ID, group.ExamID, nameEntry.Text)
		} else {
			models.InsertExamGroup(group.ExamID, nameEntry.Text)
		}
		ShowExamGroupList(w)
	})

	backBtn := widget.NewButton("⬅ Quay lại", func() {
		ShowExamGroupList(w)
	})

	w.SetContent(container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("Thông tin Nhóm Kỳ Thi", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel("Tên nhóm"), nameEntry,
			widget.NewLabel("Chọn Kỳ Thi"), examSelect,
			saveBtn, backBtn,
		),
	))
}
