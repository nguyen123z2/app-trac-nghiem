package ui

import (
	"exam-system/internal/models"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
)

// Hiển thị danh sách Student
func ShowStudentManagement(w fyne.Window) {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("🔍 Nhập mã sinh viên...")

	content := container.NewVBox()
	updateStudentList(w, content, searchEntry.Text)

	// Bắt sự kiện khi nhập vào ô tìm kiếm
	searchEntry.OnChanged = func(text string) {
		updateStudentList(w, content, text)
	}

	title := widget.NewLabel("Quản lý Sinh viên")

	// Nút quay lại trang admin
	backButton := widget.NewButton("⬅ Quay lại", func() {
		ShowAdminManagement(w)
	})

	// Nút thêm Student
	addButton := widget.NewButton("➕ Thêm Sinh viên", func() {
		ShowStudentForm(w, nil) // Form trống để thêm mới
	})

	// Nút nhập sinh viên từ Excel
	importButton := widget.NewButton("📥 Nhập từ Excel", func() {
		showImportExcelDialog(w)
	})

	// Giao diện chính
	layout := container.NewVBox(
		container.NewHBox(title, layout.NewSpacer(), backButton),
		searchEntry,
		addButton,
		importButton, // Thêm nút nhập từ Excel vào giao diện
		content,
	)

	w.SetContent(layout)
}
func ShowStudentManagement1(w fyne.Window) {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("🔍 Nhập mã sinh viên...")

	content := container.NewVBox()
	updateStudentList(w, content, searchEntry.Text)

	// Bắt sự kiện khi nhập vào ô tìm kiếm
	searchEntry.OnChanged = func(text string) {
		updateStudentList(w, content, text)
	}

	title := widget.NewLabel("Quản lý Sinh viên")

	// Nút quay lại trang admin
	backButton := widget.NewButton("⬅ Quay lại", func() {
		ShowTeacherDashboard(w)
	})

	// Nút thêm Student
	addButton := widget.NewButton("➕ Thêm Sinh viên", func() {
		ShowStudentForm1(w, nil) // Form trống để thêm mới
	})

	// Nút nhập sinh viên từ Excel
	importButton := widget.NewButton("📥 Nhập từ Excel", func() {
		showImportExcelDialog(w)
	})

	// Giao diện chính
	layout := container.NewVBox(
		container.NewHBox(title, layout.NewSpacer(), backButton),
		searchEntry,
		addButton,
		importButton, // Thêm nút nhập từ Excel vào giao diện
		content,
	)

	w.SetContent(layout)
}

// Cập nhật danh sách sinh viên theo tìm kiếm
func updateStudentList(w fyne.Window, content *fyne.Container, search string) {
	students, err := models.GetStudents(search)
	if err != nil {
		fmt.Println("❌ Lỗi khi lấy danh sách student:", err)
		return
	}

	// Xóa danh sách cũ
	content.Objects = nil

	// Tiêu đề bảng
	headers := container.NewHBox(
		widget.NewLabelWithStyle("Mã SV", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Họ tên", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Giới tính", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Ngày sinh", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Nơi sinh", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Thao tác", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	content.Add(headers)

	// Dòng phân cách giữa các tiêu đề và dữ liệu
	content.Add(widget.NewSeparator())

	// Tạo container cho bảng sinh viên
	studentList := container.NewVBox()

	// Danh sách sinh viên
	for _, student := range students {
		editButton := widget.NewButton("✏ Sửa", func() {
			ShowStudentForm(w, &student)
		})
		deleteButton := widget.NewButton("🗑 Xóa", func() {
			deleteStudent(w, student.ID)
		})

		// Dòng sinh viên
		row := container.NewHBox(
			widget.NewLabel(student.StudentCode),
			layout.NewSpacer(),
			widget.NewLabel(student.Name),
			layout.NewSpacer(),
			widget.NewLabel(student.Gender),
			layout.NewSpacer(),
			widget.NewLabel(student.DateOfBirth),
			layout.NewSpacer(),
			widget.NewLabel(student.PlaceOfBirth),
			layout.NewSpacer(),
			container.NewHBox(editButton, deleteButton),
		)

		// Dòng phân cách giữa các sinh viên
		studentList.Add(row)
		studentList.Add(widget.NewSeparator()) // Dòng phân cách giữa các dòng sinh viên
	}

	// Bọc danh sách sinh viên trong container cuộn và đặt kích thước cố định
	scrollContainer := container.NewVScroll(studentList)
	scrollContainer.SetMinSize(fyne.NewSize(700, 400)) // Đặt kích thước tối thiểu cho container cuộn

	// Thêm vào giao diện chính
	content.Add(scrollContainer)

	// Cập nhật UI
	content.Refresh()
}

// Hiển thị form thêm/sửa Student
func ShowStudentForm(w fyne.Window, student *models.Student) {
	var isEditMode bool
	if student != nil {
		isEditMode = true
	}

	title := widget.NewLabel("Thông tin Sinh viên")

	codeEntry := widget.NewEntry()
	codeEntry.SetPlaceHolder("Nhập mã sinh viên")

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nhập họ tên")

	genderEntry := widget.NewSelect([]string{"Nam", "Nữ"}, nil)

	dobEntry := widget.NewEntry()
	dobEntry.SetPlaceHolder("Nhập ngày sinh")

	pobEntry := widget.NewEntry()
	pobEntry.SetPlaceHolder("Nhập nơi sinh")

	// Kiểm tra nếu student không phải là nil
	if isEditMode {
		// Gán giá trị vào các trường nhập liệu nếu student không phải là nil
		codeEntry.SetText(student.StudentCode)
		nameEntry.SetText(student.Name)
		genderEntry.SetSelected(student.Gender)
		dobEntry.SetText(student.DateOfBirth)
		pobEntry.SetText(student.PlaceOfBirth)
	}

	// Nút Lưu
	saveButton := widget.NewButton("💾 Lưu", func() {
		studentCode := codeEntry.Text
		name := nameEntry.Text
		gender := genderEntry.Selected
		dob := dobEntry.Text
		pob := pobEntry.Text

		if isEditMode {
			models.UpdateStudent(student.ID, studentCode, name, gender, dob, pob)
		} else {
			models.AddStudent(studentCode, name, gender, dob, pob)
		}
		ShowStudentManagement(w)
	})

	// Nút Quay lại
	backButton := widget.NewButton("⬅ Quay lại", func() {
		ShowStudentManagement(w)
	})

	form := container.NewVBox(
		title,
		widget.NewLabel("Mã sinh viên"), codeEntry,
		widget.NewLabel("Họ tên"), nameEntry,
		widget.NewLabel("Giới tính"), genderEntry,
		widget.NewLabel("Ngày sinh"), dobEntry,
		widget.NewLabel("Nơi sinh"), pobEntry,
		saveButton,
		backButton,
	)

	w.SetContent(container.NewCenter(form))
}

func ShowStudentForm1(w fyne.Window, student *models.Student) {
	var isEditMode bool
	if student != nil {
		isEditMode = true
	}

	title := widget.NewLabel("Thông tin Sinh viên")

	codeEntry := widget.NewEntry()
	codeEntry.SetPlaceHolder("Nhập mã sinh viên")

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nhập họ tên")

	genderEntry := widget.NewSelect([]string{"Nam", "Nữ"}, nil)
	genderEntry.SetSelected(student.Gender)

	dobEntry := widget.NewEntry()
	dobEntry.SetPlaceHolder("Nhập ngày sinh")

	pobEntry := widget.NewEntry()
	pobEntry.SetPlaceHolder("Nhập nơi sinh")

	if isEditMode {
		codeEntry.SetText(student.StudentCode)
		nameEntry.SetText(student.Name)
		genderEntry.SetSelected(student.Gender)
		dobEntry.SetText(student.DateOfBirth)
		pobEntry.SetText(student.PlaceOfBirth)
	}

	// Nút Lưu
	saveButton := widget.NewButton("💾 Lưu", func() {
		studentCode := codeEntry.Text
		name := nameEntry.Text
		gender := genderEntry.Selected
		dob := dobEntry.Text
		pob := pobEntry.Text

		if isEditMode {
			models.UpdateStudent(student.ID, studentCode, name, gender, dob, pob)
		} else {
			models.AddStudent(studentCode, name, gender, dob, pob)
		}
		ShowStudentManagement1(w)
	})

	// Nút Quay lại
	backButton := widget.NewButton("⬅ Quay lại", func() {
		ShowStudentManagement1(w)
	})

	form := container.NewVBox(
		title,
		widget.NewLabel("Mã sinh viên"), codeEntry,
		widget.NewLabel("Họ tên"), nameEntry,
		widget.NewLabel("Giới tính"), genderEntry,
		widget.NewLabel("Ngày sinh"), dobEntry,
		widget.NewLabel("Nơi sinh"), pobEntry,
		saveButton,
		backButton,
	)

	w.SetContent(container.NewCenter(form))
}

// Xóa Student
func deleteStudent(w fyne.Window, id int) {
	dialog.ShowConfirm("Xác nhận xóa", "Bạn có chắc chắn muốn xóa sinh viên này?", func(confirmed bool) {
		if confirmed {
			models.DeleteStudent(id)
			ShowStudentManagement(w)
		}
	}, w)
}

// Hàm nhập sinh viên từ file Excel
func importStudentsFromExcel(filePath string) ([]models.Student, error) {
	// Mở file Excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Lỗi khi mở file Excel: %v", err)
	}

	// Đọc dữ liệu từ sheet đầu tiên
	rows, err := f.GetRows("Sheet1") // Đảm bảo bạn sử dụng đúng tên sheet trong file của bạn
	if err != nil {
		return nil, fmt.Errorf("Lỗi khi đọc dữ liệu từ sheet: %v", err)
	}

	var students []models.Student
	for _, row := range rows[1:] { // Bỏ qua dòng đầu tiên (tiêu đề)
		// Kiểm tra số lượng cột và đảm bảo đủ dữ liệu
		if len(row) >= 7 { // Đảm bảo có đủ 7 cột
			student := models.Student{
				StudentCode:  row[1],                // Cột 2: Mã sinh viên
				Name:         row[2] + " " + row[3], // Cột 3: Họ, Cột 4: Tên
				Gender:       row[4],                // Cột 5: Giới tính
				DateOfBirth:  row[5],                // Cột 6: Ngày sinh
				PlaceOfBirth: row[6],                // Cột 7: Nơi sinh
			}
			students = append(students, student)
		}
	}
	return students, nil
}

// Hiển thị hộp thoại chọn file Excel
func showImportExcelDialog(w fyne.Window) {
	dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
		if err != nil || r == nil {
			return
		}
		filePath := r.URI().Path()

		students, err := importStudentsFromExcel(filePath)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Lỗi khi nhập dữ liệu từ Excel: %v", err), w)
			return
		}

		// Sau khi đọc được danh sách sinh viên, bạn có thể thêm vào hệ thống
		for _, student := range students {
			models.AddStudent(student.StudentCode, student.Name, student.Gender, student.DateOfBirth, student.PlaceOfBirth)
		}

		// Cập nhật lại giao diện danh sách sinh viên
		ShowStudentManagement(w)
	}, w).Show()
}
