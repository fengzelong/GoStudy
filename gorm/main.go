package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type User struct {
	// gorm.Model
	UserId     int       `gorm:"user_id;PRIMARY_KEY"`
	UserName   string    `gorm:"user_name"`
	Sex        string    `gorm:"sex"`
	Email      string    `gorm:"email"`
	CreateTime time.Time `gorm:"create_time"`
}

// TableName 修改默认表名
func (User) TableName() string {
	return "person"
}

var Db *gorm.DB

func init() {
	//打印慢 SQL 和错误
	newLogger := logger.New(log.New(os.Stdin, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Silent,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:fl666@2022@tcp(124.223.8.183:3306)/go_test?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                                     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                                    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                                    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                                   // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("open failed")
		return
	}
	Db = db
}

var users []User

func main() {

	// 新增记录
	//user := User{UserName: "123", Sex: "男", Email: "123@163.com"}
	//AddRecordFunc(&user)

	// 批量新增
	//users := []User{
	//	{UserName: "123", Sex: "男", Email: "123@163.com", CreateTime: time.Now()},
	//	{UserName: "456", Sex: "女", Email: "456@163.com", CreateTime: time.Now()},
	//	{UserName: "789", Sex: "女", Email: "789@163.com", CreateTime: time.Now()},
	//}
	//BatchAddFunc(&users)

	// 查询
	u, err := QueryFunc(2)
	if err != nil {
		fmt.Printf("query failed, reason: %v\n", err)
		return
	} else {
		// 更新
		// UpdateFunc(u)

		// 删除
		DeleteFunc(u)
	}
}

// CreateDb 创建数据库
func CreateDb() {

}

// BeforeCreate 钩子方法
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	println("add record before, u.UserName: ", u.UserName)
	return
}

// AfterCreate 钩子方法
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	//println("add record after")
	return
}

// AddRecordFunc 新增记录
func AddRecordFunc(user *User) bool {
	tx := Db.Select("UserName", "Sex", "Email").Create(&user)
	if tx.Error != nil {
		fmt.Printf("add record err, %v", tx.Error)
		return false
	}
	fmt.Println("add record success")
	return true
}

// BatchAddFunc 批量新增
func BatchAddFunc(users *[]User) bool {
	tx := Db.Select("UserName", "Sex", "Email", "CreateTime").CreateInBatches(&users, 2)
	if tx.Error != nil {
		fmt.Printf("batch add err %v", tx.Error)
		return false
	}
	fmt.Println("batch add success")
	return true
}

// QueryFunc 查询
func QueryFunc(index int) (*User, error) {
	//tx := Db.First(&User{})
	//tx := Db.Last(&User{})
	//fmt.Printf("record = %v\n", tx.Statement.Dest)

	//tx := Db.Where("user_name like ?", "%张三%").Find(&users)

	// 迭代查询
	tx := Db.Table("person").Where([]int{2, 3, 7}).Find(&users)
	rows, _ := tx.Rows()
	defer rows.Close()

	var u1 User
	var i int
	for rows.Next() {
		var user User
		i++
		Db.ScanRows(rows, &user)

		//fmt.Println("scan row:", i)

		if i == index {
			u1 = user
		}

		fmt.Printf("record %d = %v, create at %v\n", i, user.UserName, user.CreateTime)
	}
	return &u1, nil
}

// UpdateFunc 更新
func UpdateFunc(user *User) bool {
	// 保存所有字段
	//user.Email = "test01@gmail.com"
	//user.CreateTime = time.Now()
	//tx := Db.Save(&user)

	// 保存部分字段
	tx := Db.Model(&user).Updates(User{Email: "test02@163.com", CreateTime: time.Now()})

	if tx.Error != nil {
		fmt.Println("save failed")
		return false
	}
	return true
}

// DeleteFunc 删除
func DeleteFunc(user *User) bool {
	tx := Db.Delete(&user)
	if tx.Error != nil {
		fmt.Println("del failed")
		return false
	}
	return true
}
