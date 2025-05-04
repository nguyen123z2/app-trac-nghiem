// internal/ui/question_bank.go
package ui

import (
	"exam-system/internal/models"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const questionSetsPerPage = 10

// ShowQuestionBank hiển thị danh sách bộ câu hỏi với phân trang và tìm kiếm.
// back là hàm gọi để quay về dashboard (Admin hoặc Teacher).
func ShowQuestionBank(w fyne.Window, back func(fyne.Window)) {
	allSets, err := models.GetQuestionSets()
	if err != nil {
		fmt.Println("❌ Lỗi khi lấy danh sách bộ câu hỏi:", err)
		return
	}

	currentPage := 0
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("🔍 Tìm kiếm bộ câu hỏi...")

	pageLabel := widget.NewLabel("")
	listContainer := container.NewVBox()

	updateList := func() {
		listContainer.Objects = nil
		keyword := strings.ToLower(searchEntry.Text)

		// Lọc theo từ khóa
		var filtered []models.QuestionSet
		for _, s := range allSets {
			if keyword == "" || strings.Contains(strings.ToLower(s.Name), keyword) {
				filtered = append(filtered, s)
			}
		}

		// Tính tổng trang
		totalPages := (len(filtered) + questionSetsPerPage - 1) / questionSetsPerPage
		if totalPages == 0 {
			pageLabel.SetText("Không có bộ câu hỏi.")
			listContainer.Refresh()
			return
		}
		if currentPage >= totalPages {
			currentPage = totalPages - 1
		}
		start := currentPage * questionSetsPerPage
		end := start + questionSetsPerPage
		if end > len(filtered) {
			end = len(filtered)
		}

		// Hiển thị từng phần
		for _, set := range filtered[start:end] {
			setCopy := set // capture
			btn := widget.NewButton(set.Name, func() {
				ShowQuestionList(w, setCopy.ID, back)
			})
			listContainer.Add(
				container.NewHBox(
					btn,
					layout.NewSpacer(),
					widget.NewButton("✏️", func() {
						ShowEditQuestionSetForm(w, setCopy, back)
					}),
					widget.NewButton("🗑️", func() {
						deleteQuestionSet(w, setCopy.ID, back)
					}),
				),
			)
		}

		pageLabel.SetText(fmt.Sprintf("Trang %d/%d", currentPage+1, totalPages))
		listContainer.Refresh()
	}

	// Nút phân trang
	prevBtn := widget.NewButton("⬅ Trang trước", func() {
		if currentPage > 0 {
			currentPage--
			updateList()
		}
	})
	nextBtn := widget.NewButton("Trang sau ➡", func() {
		// Lấy số lượng filtered để phân trang chính xác
		keyword := strings.ToLower(searchEntry.Text)
		count := 0
		for _, s := range allSets {
			if keyword == "" || strings.Contains(strings.ToLower(s.Name), keyword) {
				count++
			}
		}
		if (currentPage+1)*questionSetsPerPage < count {
			currentPage++
			updateList()
		}
	})

	searchEntry.OnChanged = func(_ string) {
		currentPage = 0
		updateList()
	}

	updateList()

	// Các nút chức năng trên cùng
	topBar := container.NewHBox(
		widget.NewButton("➕ Thêm bộ câu hỏi", func() {
			ShowAddQuestionSetForm(w, back)
		}),
		layout.NewSpacer(),
		widget.NewButton("🔙 Quay lại", func() { back(w) }),
	)

	w.SetContent(container.NewVBox(
		topBar,
		widget.NewLabelWithStyle("📖 Danh sách bộ câu hỏi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		searchEntry,
		listContainer,
		container.NewHBox(prevBtn, pageLabel, nextBtn),
	))
}

// ShowAddQuestionSetForm hiển thị form thêm mới bộ câu hỏi.
// back để quay lại danh sách sau khi lưu hoặc huỷ.
func ShowAddQuestionSetForm(w fyne.Window, back func(fyne.Window)) {
	// TODO: thay bằng teacherID thực tế từ session
	teacherID := 1

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Tên bộ câu hỏi")

	form := widget.NewForm(
		widget.NewFormItem("Tên bộ câu hỏi", nameEntry),
	)
	form.OnSubmit = func() {
		if err := models.InsertQuestionSet(teacherID, nameEntry.Text); err != nil {
			dialog.ShowError(err, w)
		} else {
			ShowQuestionBank(w, back)
		}
	}

	w.SetContent(container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("➕ Thêm bộ câu hỏi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewButton("🔙 Quay lại", func() { ShowQuestionBank(w, back) }),
		),
		form,
	))
}

// ShowEditQuestionSetForm hiển thị form sửa bộ câu hỏi.
// Nhận thêm back để quay lại danh sách sau khi lưu.
func ShowEditQuestionSetForm(w fyne.Window, set models.QuestionSet, back func(fyne.Window)) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(set.Name)

	form := widget.NewForm(
		widget.NewFormItem("Tên bộ câu hỏi", nameEntry),
	)
	form.OnSubmit = func() {
		models.UpdateQuestionSet(set.ID, set.TeacherID, nameEntry.Text)
		ShowQuestionBank(w, back)
	}

	w.SetContent(container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("✏️ Sửa bộ câu hỏi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewButton("🔙 Quay lại", func() { ShowQuestionBank(w, back) }),
		),
		form,
	))
}

// deleteQuestionSet hỏi xác nhận rồi xóa, sau đó gọi lại back.
func deleteQuestionSet(w fyne.Window, id int, back func(fyne.Window)) {
	dialog.ShowConfirm("Xác nhận xóa",
		"Bạn có chắc chắn muốn xóa bộ câu hỏi này?",
		func(ok bool) {
			if ok {
				if err := models.DeleteQuestionSet(id); err != nil {
					dialog.ShowError(err, w)
				}
				ShowQuestionBank(w, back)
			}
		},
		w)
}
