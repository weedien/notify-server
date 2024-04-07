package adapters

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库连接
var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "weedien:031209@tcp(remote:3306)/countdown?parseTime=true&&loc=Local",
		DefaultStringSize:         255,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}
