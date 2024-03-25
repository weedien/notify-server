package models

import (
	"database/sql"
	"time"
)

// 倒计时结构体
type Countdown struct {
	ID        string
	QueryCode string
	ExpireAt  time.Time
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Remark    sql.NullString
	Message   sql.NullString
}
