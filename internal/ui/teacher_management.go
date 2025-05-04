package ui

import (
	"exam-system/internal/models"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowTeacherManagement(w fyne.Window, back func(fyne.Window)) {
	teachers, err := models.GetTeachers()
	if err != nil {
		fmt.Println("❌ Lỗi khi lấy danh sách teacher:", err)
		return
	}

	// Header với nút quay lại
	header := container.NewHBox(
		widget.NewLabel("👩‍🏫 Quản lý Giáo viên"),
		layout.NewSpacer(),
		widget.NewButton("⬅ Quay lại", func() { back(w) }),
	)

	// Nút thêm
	addBtn := widget.NewButton("➕ Thêm Giáo viên", func() {
		ShowTeacherForm(w, nil, back)
	})

	// Tiêu đề bảng
	headers := container.NewHBox(
		widget.NewLabelWithStyle("Họ tên", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Email", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Vai trò", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Hành động", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	// Container chứa các hàng
	rows := container.NewVBox(headers, widget.NewSeparator())

	for _, t := range teachers {
		teacher := t // capture
		editBtn := widget.NewButton("✏ Sửa", func() {
			ShowTeacherForm(w, &teacher, back)
		})
		delBtn := widget.NewButton("🗑 Xóa", func() {
			dialog.ShowConfirm("Xác nhận xóa",
				"Bạn có chắc chắn muốn xóa giáo viên này?",
				func(ok bool) {
					if ok {
						if err := models.DeleteTeacher(teacher.ID); err != nil {
							dialog.ShowInformation("Lỗi", "Xóa thất bại!", w)
						}
						ShowTeacherManagement(w, back)
					}
				}, w)
		})

		row := container.NewHBox(
			widget.NewLabel(teacher.Name), layout.NewSpacer(),
			widget.NewLabel(teacher.Email), layout.NewSpacer(),
			widget.NewLabel(teacher.Role), layout.NewSpacer(),
			container.NewHBox(editBtn, delBtn),
		)

		// Thêm row rồi mới thêm separator
		rows.Add(row)
		rows.Add(widget.NewSeparator())
	}

	content := container.NewVBox(header, addBtn, rows)
	w.SetContent(content)
}

func ShowTeacherForm(w fyne.Window, teacher *models.Teacher, back func(fyne.Window)) {
	isEdit := teacher != nil

	// Header và nút quay lại
	header := container.NewHBox(
		widget.NewLabel("📝 Thông tin Giáo viên"),
		layout.NewSpacer(),
		widget.NewButton("⬅ Quay lại", func() {
			ShowTeacherManagement(w, back)
		}),
	)

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Họ tên")
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")
	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("Mật khẩu")
	roleSelect := widget.NewSelect([]string{"admin", "teacher"}, nil)
	roleSelect.SetSelected("teacher")

	if isEdit {
		nameEntry.SetText(teacher.Name)
		emailEntry.SetText(teacher.Email)
		roleSelect.SetSelected(teacher.Role)
	}

	saveBtn := widget.NewButton("💾 Lưu", func() {
		name, email, pwd, role := nameEntry.Text, emailEntry.Text, passEntry.Text, roleSelect.Selected

		if isEdit {
			// Cập nhật
			if pwd != "" {
				_ = models.UpdateTeacherWithPassword(teacher.ID, name, email, pwd, role)
			} else {
				_ = models.UpdateTeacherWithoutPassword(teacher.ID, name, email, role)
			}
		} else {
			// Thêm mới
			if pwd == "" {
				dialog.ShowInformation("Lỗi", "Vui lòng nhập mật khẩu!", w)
				return
			}
			hashed := models.HashMD5(pwd)
			_ = models.AddTeacher(name, email, hashed, role)
		}
		ShowTeacherManagement(w, back)
	})

	form := container.NewVBox(
		header,
		widget.NewLabel("Họ tên"), nameEntry,
		widget.NewLabel("Email"), emailEntry,
		widget.NewLabel("Mật khẩu"), passEntry,
		widget.NewLabel("Vai trò"), roleSelect,
		saveBtn,
	)

	w.SetContent(container.NewCenter(form))
}
