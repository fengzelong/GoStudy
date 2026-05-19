package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	UserId     int       `gorm:"user_id;PRIMARY_KEY"`
	UserName   string    `gorm:"user_name"`
	Sex        string    `gorm:"sex"`
	Email      string    `gorm:"email"`
	CreateTime time.Time `gorm:"create_time"`
}

// TableName 修改默认表名。
func (User) TableName() string {
	return "person"
}

var Db *gorm.DB

func openDB() (*gorm.DB, error) {
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Silent,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})

	dsn := os.Getenv("GORM_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(127.0.0.1:3306)/go_test?charset=utf8&parseTime=True&loc=Local"
		fmt.Println("未设置 GORM_DSN，使用本地默认连接示例")
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

var users []User

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("打开 Gorm 连接失败: %v", err)
	}
	Db = db

	// 新增记录
	// user := User{UserName: "123", Sex: "男", Email: "123@163.com"}
	// AddRecordFunc(&user)

	// 批量新增
	// users := []User{
	// 	{UserName: "123", Sex: "男", Email: "123@163.com", CreateTime: time.Now()},
	// 	{UserName: "456", Sex: "女", Email: "456@163.com", CreateTime: time.Now()},
	// 	{UserName: "789", Sex: "女", Email: "789@163.com", CreateTime: time.Now()},
	// }
	// BatchAddFunc(&users)

	// 查询
	u, err := QueryFunc(2)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}

	// 更新
	// UpdateFunc(u)

	// 删除
	DeleteFunc(u)
}

// CreateDb 创建数据库示例。
func CreateDb() {
}

// BeforeCreate 是新增记录前的钩子方法。
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	println("新增记录前，用户名: ", u.UserName)
	return
}

// AfterCreate 是新增记录后的钩子方法。
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	return
}

// AddRecordFunc 新增记录。
func AddRecordFunc(user *User) bool {
	tx := Db.Select("UserName", "Sex", "Email").Create(&user)
	if tx.Error != nil {
		fmt.Printf("新增记录失败: %v", tx.Error)
		return false
	}
	fmt.Println("新增记录成功")
	return true
}

// BatchAddFunc 批量新增记录。
func BatchAddFunc(users *[]User) bool {
	tx := Db.Select("UserName", "Sex", "Email", "CreateTime").CreateInBatches(&users, 2)
	if tx.Error != nil {
		fmt.Printf("批量新增失败: %v", tx.Error)
		return false
	}
	fmt.Println("批量新增成功")
	return true
}

// QueryFunc 查询记录。
func QueryFunc(index int) (*User, error) {
	tx := Db.Table("person").Where([]int{2, 3, 7}).Find(&users)
	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u1 User
	var i int
	for rows.Next() {
		var user User
		i++
		Db.ScanRows(rows, &user)

		if i == index {
			u1 = user
		}

		fmt.Printf("记录 %d = %v，创建时间 %v\n", i, user.UserName, user.CreateTime)
	}
	return &u1, nil
}

// UpdateFunc 更新记录。
func UpdateFunc(user *User) bool {
	tx := Db.Model(&user).Updates(User{Email: "test02@163.com", CreateTime: time.Now()})
	if tx.Error != nil {
		fmt.Println("保存失败")
		return false
	}
	return true
}

// DeleteFunc 删除记录。
func DeleteFunc(user *User) bool {
	tx := Db.Delete(&user)
	if tx.Error != nil {
		fmt.Println("删除失败")
		return false
	}
	return true
}
