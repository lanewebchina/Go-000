package dao

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	UserTable = "t_user"
)

var DB *gorm.DB
var err error

func init() {
	connect := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", "develop", "123456", "127.0.0.1", 3306, "test")
	DB, err = gorm.Open("mysql", connect)
	if err != nil {
		log.Fatal("Mysql Connect Error:", err.Error())
	}
	DB.SingularTable(true)
	logPath := os.Getenv("CONFIG")
	DBLogFile, err := os.OpenFile(logPath+"mysql.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666) //mysql日志
	if err != nil {
		log.Println("gorm db log path init error:", err.Error())
	}
	DB.LogMode(true)
	DB.SetLogger(log.New(DBLogFile, "", log.Ldate|log.Ltime|log.Lshortfile))
}
