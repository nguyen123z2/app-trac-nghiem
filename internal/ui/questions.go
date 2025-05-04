// internal/ui/question_list.go
package ui

import (
	"fmt"
	"strings"

	"exam-system/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Số câu hỏi trên mỗi trang
const questionsPerPage = 10

// ShowQuestionList hiển thị danh sách câu hỏi của một bộ câu hỏi,
// back là hàm quay lại QuestionBank tương ứng (Admin hoặc Teacher).
func ShowQuestionList(w fyne.Window, questionSetID int, back func(fyne.Window)) {
	fmt.Println("Đang hiển thị danh sách câu hỏi với QuestionSetID:", questionSetID)

	questions, err := models.GetQuestions(questionSetID)
	if err != nil {
		fmt.Println("Lỗi lấy danh sách câu hỏi:", err)
		return
	}

	currentPage := 0
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("🔍 Tìm kiếm câu hỏi...")

	pageLabel := widget.NewLabel("Trang 1/1")
	listContainer := container.NewVBox()

	updateList := func() {
		listContainer.Objects = nil
		keyword := strings.ToLower(searchEntry.Text)

		var filtered []models.Question
		for _, q := range questions {
			if keyword == "" || strings.Contains(strings.ToLower(q.Content), keyword) {
				filtered = append(filtered, q)
			}
		}

		totalPages := (len(filtered) + questionsPerPage - 1) / questionsPerPage
		if totalPages == 0 {
			pageLabel.SetText("Không có câu hỏi nào.")
			listContainer.Refresh()
			return
		}
		if currentPage >= totalPages {
			currentPage = totalPages - 1
		}
		if currentPage < 0 {
			currentPage = 0
		}

		start := currentPage * questionsPerPage
		end := start + questionsPerPage
		if end > len(filtered) {
			end = len(filtered)
		}

		for _, q := range filtered[start:end] {
			qCopy := q // capture
			btn := widget.NewButton(qCopy.Content, func() {
				ShowEditQuestionForm(w, qCopy, questionSetID, back)
			})
			listContainer.Add(btn)
		}

		pageLabel.SetText(fmt.Sprintf("Trang %d/%d", currentPage+1, totalPages))
		listContainer.Refresh()
	}

	prevBtn := widget.NewButton("⬅ Trang trước", func() {
		if currentPage > 0 {
			currentPage--
			updateList()
		}
	})
	nextBtn := widget.NewButton("Trang sau ➡", func() {
		if (currentPage+1)*questionsPerPage < len(questions) {
			currentPage++
			updateList()
		}
	})

	searchEntry.OnChanged = func(_ string) {
		currentPage = 0
		updateList()
	}

	updateList()

	w.SetContent(container.NewVBox(
		widget.NewButton("➕ Thêm câu hỏi", func() {
			ShowAddQuestionForm(w, questionSetID, back)
		}),
		widget.NewButton("🔙 Quay lại", func() {
			back(w)
		}),
		widget.NewLabelWithStyle("📖 Danh sách câu hỏi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		searchEntry,
		listContainer,
		container.NewHBox(prevBtn, pageLabel, nextBtn),
	))
}

// ShowAddQuestionForm hiển thị form thêm câu hỏi cho bộ câu hỏi.
// back để quay lại ShowQuestionList sau khi lưu hoặc huỷ.
func ShowAddQuestionForm(w fyne.Window, questionSetID int, back func(fyne.Window)) {
	contentEntry := widget.NewEntry()
	optionA := widget.NewEntry()
	optionB := widget.NewEntry()
	optionC := widget.NewEntry()
	optionD := widget.NewEntry()
	correct := widget.NewEntry()
	difficulty := widget.NewEntry()

	form := widget.NewForm(
		widget.NewFormItem("Nội dung", contentEntry),
		widget.NewFormItem("Đáp án A", optionA),
		widget.NewFormItem("Đáp án B", optionB),
		widget.NewFormItem("Đáp án C", optionC),
		widget.NewFormItem("Đáp án D", optionD),
		widget.NewFormItem("Đáp án đúng", correct),
		widget.NewFormItem("Độ khó", difficulty),
	)
	form.OnSubmit = func() {
		// Lấy teacherID của bộ câu hỏi
		set, err := models.GetQuestionSetByID(questionSetID)
		if err != nil {
			fmt.Println("❌ Lỗi lấy thông tin bộ câu hỏi:", err)
			return
		}

		answers := []string{optionA.Text, optionB.Text, optionC.Text, optionD.Text}
		qid, err := models.InsertQuestion(
			contentEntry.Text,
			answers,
			correct.Text,
			difficulty.Text,
			questionSetID,
			set.TeacherID,
		)
		if err != nil {
			fmt.Println("❌ Lỗi thêm câu hỏi:", err)
			return
		}
		if err := models.InsertQuestionAnswer(qid, answers, correct.Text); err != nil {
			fmt.Println("❌ Lỗi thêm đáp án:", err)
			return
		}
		ShowQuestionList(w, questionSetID, back)
	}

	w.SetContent(container.NewVBox(
		form,
		widget.NewButton("🔙 Quay lại", func() {
			back(w)
		}),
	))
}

// ShowEditQuestionForm hiển thị form sửa câu hỏi.
// nhận thêm questionSetID và back để quay lại đúng chỗ.
func ShowEditQuestionForm(w fyne.Window, q models.Question, questionSetID int, back func(fyne.Window)) {
	contentEntry := widget.NewEntry()
	contentEntry.SetText(q.Content)
	qa, err := models.GetQuestionAnswer(q.ID)
	if err != nil {
		fmt.Println("Lỗi lấy đáp án:", err)
		return
	}

	optionA := widget.NewEntry()
	optionA.SetText(qa.Answers[0])
	optionB := widget.NewEntry()
	optionB.SetText(qa.Answers[1])
	optionC := widget.NewEntry()
	optionC.SetText(qa.Answers[2])
	optionD := widget.NewEntry()
	optionD.SetText(qa.Answers[3])
	correct := widget.NewEntry()
	correct.SetText(qa.CorrectAnswer)
	difficulty := widget.NewEntry()
	difficulty.SetText(q.Difficulty)

	form := widget.NewForm(
		widget.NewFormItem("Nội dung", contentEntry),
		widget.NewFormItem("Đáp án A", optionA),
		widget.NewFormItem("Đáp án B", optionB),
		widget.NewFormItem("Đáp án C", optionC),
		widget.NewFormItem("Đáp án D", optionD),
		widget.NewFormItem("Đáp án đúng", correct),
		widget.NewFormItem("Độ khó", difficulty),
	)
	form.OnSubmit = func() {
		if err := models.UpdateQuestion(q.ID, contentEntry.Text, difficulty.Text); err != nil {
			fmt.Println("❌ Lỗi cập nhật câu hỏi:", err)
		}
		answers := []string{optionA.Text, optionB.Text, optionC.Text, optionD.Text}
		if err := models.UpdateQuestionAnswer(q.ID, answers, correct.Text); err != nil {
			fmt.Println("❌ Lỗi cập nhật đáp án:", err)
		}
		ShowQuestionList(w, questionSetID, back)
	}

	w.SetContent(container.NewVBox(
		form,
		widget.NewButton("🔙 Quay lại", func() {
			back(w)
		}),
	))
}
