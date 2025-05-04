package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func InitDB() {
	var err error
	database, err = sql.Open("mysql", "exam_system:3Ytx6ZjxZcStENsC@tcp(servernguyen.zapto.org:3306)/exam_system?parseTime=true")
	if err != nil {
		fmt.Println("❌ Lỗi kết nối MySQL:", err)
		panic(err)
	}

	if err = database.Ping(); err != nil {
		fmt.Println("❌ Không thể kết nối đến MySQL:", err)
		panic(err)
	}

	fmt.Println("✅ Kết nối MySQL thành công!")
}

func GetDB() *sql.DB {
	return database
}
