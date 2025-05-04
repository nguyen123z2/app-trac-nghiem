package auth

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"exam-system/internal/db"
	"fmt"
)

var currentUserID int

// hashMD5 băm MD5
func hashMD5(password string) string {
	h := md5.Sum([]byte(password))
	return hex.EncodeToString(h[:])
}

// Authenticate xác thực và lưu currentUserID
func Authenticate(username, password string) bool {
	database := db.GetDB()
	if database == nil {
		fmt.Println("❌ Database chưa được khởi tạo!")
		return false
	}

	hashed := hashMD5(password)
	var id int
	err := database.QueryRow(
		"SELECT id FROM users WHERE email = ? AND password = ?",
		username, hashed,
	).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("❌ Lỗi truy vấn:", err)
		}
		return false
	}

	currentUserID = id
	return true
}

// GetUserRole giữ nguyên
func GetUserRole(username string) (string, error) {
	database := db.GetDB()
	if database == nil {
		return "", fmt.Errorf("❌ Database chưa được khởi tạo!")
	}

	var role string
	err := database.QueryRow(
		"SELECT role FROM users WHERE email = ?",
		username,
	).Scan(&role)
	return role, err
}

// GetCurrentUserID trả về ID của user đã đăng nhập
func GetCurrentUserID() (int, error) {
	if currentUserID == 0 {
		return 0, errors.New("chưa có user đăng nhập")
	}
	return currentUserID, nil
}
