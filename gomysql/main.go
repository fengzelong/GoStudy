package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Person struct {
	UserId   int    `db:"user_id"`
	UserName string `db:"username"`
	Sex      string `db:"sex"`
	Email    string `db:"email"`
}

type Place struct {
	Country string `db:"country"`
	City    string `db:"city"`
	TelCode int    `db:"telcode"`
}

var Db *sqlx.DB

func openDB() (*sqlx.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(127.0.0.1:3306)/go_test"
		fmt.Println("未设置 MYSQL_DSN，使用本地默认连接示例")
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("打开 MySQL 失败: %v", err)
	}
	Db = db
	defer Db.Close()

	// 新增数据到 MySQL
	// insertFunc()

	// 查询 MySQL 数据
	// selectFunc()

	// 更新 MySQL 数据
	// updateFunc("张三3", 3)

	// 删除 MySQL 数据
	delFunc(6)
}

// insertFunc 新增 MySQL 记录。
func insertFunc() {
	res, err := Db.NamedExec("insert into person(username, sex, email)values(:username, :sex, :email)",
		map[string]interface{}{
			"username": "张三",
			"sex":      "男",
			"email":    "zhangsan@163.com",
		})
	if err != nil {
		fmt.Println("执行新增失败: ", err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("获取新增 ID 失败: ", err)
		return
	}
	fmt.Println("新增成功:", id)
}

// selectFunc 查询 MySQL 记录。
func selectFunc() {
	p1 := Person{}
	err := Db.Get(&p1, "SELECT * FROM person WHERE user_id = ?", 4)
	if err != nil {
		fmt.Printf("查询失败 = %v\n", err)
		return
	}
	fmt.Printf("%#v\n", p1)
}

// updateFunc 更新 MySQL 记录。
func updateFunc(username string, id int) bool {
	res, err := Db.Exec("update person set username = ? where user_id = ?", username, id)
	if err != nil {
		fmt.Printf("更新失败 = %v\n", err)
		return false
	}

	row, err := res.RowsAffected()
	if err != nil {
		fmt.Println("获取影响行数失败: ", err)
		return false
	}
	fmt.Println("更新成功 = ", row)
	return true
}

// delFunc 删除 MySQL 记录。
func delFunc(id int) bool {
	res, err := Db.Exec("delete from person where user_id = ?", id)
	if err != nil {
		fmt.Println("执行删除失败: ", err)
		return false
	}

	row, err := res.RowsAffected()
	if err != nil {
		fmt.Println("获取影响行数失败: ", err)
		return false
	}

	fmt.Println("删除成功: ", row)
	return true
}
