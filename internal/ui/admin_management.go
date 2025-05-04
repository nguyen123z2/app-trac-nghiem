package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Giao diện quản lý Admin
func ShowAdminManagement(w fyne.Window) {
	// Tạo menu cho quản lý Admin
	menu := container.NewVBox(
		widget.NewButton("📚 Quản lý Teacher", func() {
			// back = ShowAdminManagement để quay lại chính nó
			ShowTeacherManagement(w, ShowAdminManagement)
		}),
		widget.NewButton("🎓 Quản lý Student", func() {
			ShowStudentManagement(w)
		}),
		widget.NewButton("📖 Ngân hàng câu hỏi", func() {
			ShowQuestionBank(w, ShowAdminManagement)
		}),

		widget.NewButton("📅 Danh sách Kỳ Thi", func() {
			ShowExamListAdmin(w)
		}),
		widget.NewButton("Danh sách Nhóm", func() {
			ShowExamGroupList(w)
		}),
		widget.NewButton("⬅ Đăng xuất", func() {
			ShowLogin(w)
		}),
	)

	w.SetContent(menu)
}
