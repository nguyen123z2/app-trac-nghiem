package models

import (
	"crypto/md5"
	"encoding/hex"
	"exam-system/internal/db"
	"fmt"
)

// Cấu trúc Teacher
type Teacher struct {
	ID    int
	Name  string
	Email string
	Role  string
}

// Mã hóa mật khẩu MD5
func HashMD5(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Lấy danh sách giáo viên từ database
func GetTeachers() ([]Teacher, error) {
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	query := "SELECT id, name, email, role FROM users WHERE role IN ('admin', 'teacher')"
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []Teacher
	for rows.Next() {
		var t Teacher
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Role)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, nil
}

// Thêm mới giáo viên
func AddTeacher(name, email, password, role string) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	hashedPassword := HashMD5(password) // Mã hóa mật khẩu

	query := "INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)"
	_, err := database.Exec(query, name, email, hashedPassword, role)
	return err
}

// Cập nhật thông tin giáo viên
func UpdateTeacher(id int, name, email, role string) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	query := "UPDATE users SET name = ?, email = ?, role = ? WHERE id = ?"
	_, err := database.Exec(query, name, email, role, id)
	return err
}

// Xóa giáo viên
func DeleteTeacher(id int) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	query := "DELETE FROM users WHERE id = ?"
	_, err := database.Exec(query, id)
	return err
}
func UpdateTeacherWithoutPassword(id int, name, email, role string) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	// Cập nhật thông tin mà không thay đổi mật khẩu
	query := "UPDATE users SET name = ?, email = ?, role = ? WHERE id = ?"
	_, err := database.Exec(query, name, email, role, id)
	return err
}
func UpdateTeacherWithPassword(id int, name, email, password, role string) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	// Mã hóa mật khẩu mới trước khi lưu
	hashedPassword := HashMD5(password)

	// Cập nhật thông tin cùng với mật khẩu mới
	query := "UPDATE users SET name = ?, email = ?, password = ?, role = ? WHERE id = ?"
	_, err := database.Exec(query, name, email, hashedPassword, role, id)
	return err
}
func GetTeacherPassword(id int) (string, error) {
	database := db.GetDB()
	if database == nil {
		return "", fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	var password string
	query := "SELECT password FROM users WHERE id = ?"
	err := database.QueryRow(query, id).Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}
func GetTeacherID(username string) (int, error) {
	var teacherID int
	// Truy vấn để lấy teacherID từ cơ sở dữ liệu
	err := db.GetDB().QueryRow("SELECT id FROM teachers WHERE email = ?", username).Scan(&teacherID)
	if err != nil {
		return 0, fmt.Errorf("lỗi khi lấy ID giáo viên: %v", err)
	}
	return teacherID, nil
}
