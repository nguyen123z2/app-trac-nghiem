package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Giao diện quản lý Teacher
func ShowTeacherDashboard(w fyne.Window) {
	menu := container.NewVBox(
		widget.NewButton("🎓 Quản lý Student", func() {
			// gọi bản dành cho teacher, có nút quay lại về dashboard
			ShowStudentManagement1(w)
		}),
		widget.NewButton("📖 Ngân hàng câu hỏi", func() {
			// truyền ShowTeacherDashboard làm back callback
			ShowQuestionBank(w, ShowTeacherDashboard)
		}),
		widget.NewButton("📅 Danh sách Kỳ Thi", func() {
			ShowExamListTeacher(w)
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
