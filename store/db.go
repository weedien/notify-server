package store

import (
	"database/sql"
	"log"
)

// 数据库连接
var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "weedien:031209@tcp(remote:3306)/countdown?parseTime=true&&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	// 创建倒计时表
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS countdowns (
		id VARCHAR(20) PRIMARY KEY,
		query_code VARCHAR(6) NOT NULL,
		expire_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME DEFAULT NULL,
		message VARCHAR(255) DEFAULT NULL,
		remark VARCHAR(255) DEFAULT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
}
