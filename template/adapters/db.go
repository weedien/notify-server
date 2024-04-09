package adapters

import (
	"fmt"
	"github.com/weedien/notify-server/common/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库连接
var db *gorm.DB

func InitDB() {
	var err error
	dsn := buildDSN()
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}
	//sqlDB, err := db.DB()
	//if err != nil {
	//	panic(err)
	//}
	//sqlDB.SetMaxIdleConns(10)
	//sqlDB.SetMaxOpenConns(100)
	//sqlDB.SetConnMaxLifetime(0)
}

func buildDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		config.Config().Database.Username,
		config.Config().Database.Password,
		config.Config().Database.Host,
		config.Config().Database.Port,
		config.Config().Database.Database,
		config.Config().Database.Params)
}

func DB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}
