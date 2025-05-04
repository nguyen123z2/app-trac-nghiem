package main

import (
	"exam-system/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Exam System")

	// Đặt cửa sổ luôn ở kích thước Full HD (1920x1080)
	w.Resize(fyne.NewSize(1180, 620))
	w.SetFixedSize(true) // Ngăn người dùng thay đổi kích thước cửa sổ

	ui.ShowAdminManagement(w)
	ui.ShowTeacherDashboard(w)

	ui.ShowLogin(w)
	w.ShowAndRun()
}
