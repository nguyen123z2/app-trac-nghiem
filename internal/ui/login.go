package ui

import (
	"exam-system/internal/auth"
	"exam-system/internal/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Hiển thị màn hình đăng nhập
func ShowLogin(w fyne.Window) {
	db.InitDB() // Khởi tạo database

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Nhập email")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Nhập mật khẩu")

	loginButton := widget.NewButton("Đăng nhập", func() {
		username := usernameEntry.Text
		password := passwordEntry.Text

		// Xác thực đăng nhập
		if auth.Authenticate(username, password) {
			// Lấy vai trò người dùng
			userRole, err := auth.GetUserRole(username)
			if err != nil {
				dialog.ShowInformation("Lỗi", "Không thể xác định vai trò người dùng!", w)
				return
			}

			// Chuyển hướng theo vai trò
			switch userRole {
			case "admin":
				ShowAdminManagement(w)
			case "teacher":
				ShowTeacherDashboard(w)
			default:
				dialog.ShowInformation("Lỗi", "Bạn không có quyền truy cập!", w)
			}
		} else {
			dialog.ShowInformation("Lỗi", "Sai thông tin đăng nhập", w)
		}
	})

	form := container.NewVBox(
		widget.NewLabel("HỆ THỐNG THI TRỰC TUYẾN"),
		usernameEntry,
		passwordEntry,
		loginButton,
	)

	w.SetContent(container.NewCenter(form))
}
