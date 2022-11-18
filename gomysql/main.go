package main

import (
	"fmt"
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

func init() {
	db, err := sqlx.Open("mysql", "root:fl666@2022@tcp(124.223.8.183:3306)/go_test")
	if err != nil {
		fmt.Println("open mysql failed", err)
		return
	} else {
		fmt.Println("open mysql success")
	}
	Db = db
}

func main() {
	//新增数据到mysql
	//insertFunc()

	//查询mysql数据
	//selectFunc()

	//更新mysql数据
	//updateFunc("张三3", 3)

	//删除mysql数据
	delFunc(6)

	defer Db.Close()
}

// insertFunc mysql新增
func insertFunc() {
	res, err := Db.NamedExec("insert into person(username, sex, email)values(:username, :sex, :email)",
		map[string]interface{}{
			"username": "张三",
			"sex":      "男",
			"email":    "zhangsan@163.com",
		})
	if err != nil {
		fmt.Println("exec failed1, ", err)
		return
	}
	id, err1 := res.LastInsertId()
	if err1 != nil {
		fmt.Println("exec failed2, ", err1)
		return
	}
	fmt.Println("insert success:", id)
}

// selectFunc mysql查询
func selectFunc() {
	p1 := Person{}
	err := Db.Get(&p1, "SELECT * FROM person WHERE user_id = ?", 4)
	if err != nil {
		fmt.Printf("query error = %v\n", err)
	} else {
		fmt.Printf("%#v\n", p1)
	}
}

// updateFunc mysql更新
func updateFunc(username string, id int) bool {
	res, err := Db.Exec("update person set username = ? where user_id = ?", username, id)
	if err != nil {
		fmt.Printf("update error = %v\n", err)
		return false
	}
	row, err1 := res.RowsAffected()
	if err1 != nil {
		fmt.Println("rows failed, ", err1)
		return false
	}
	fmt.Println("update success = ", row)
	return true
}

// delFunc mysql删除
func delFunc(id int) bool {
	res, err := Db.Exec("delete from person where user_id = ?", id)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return false
	}

	row, err1 := res.RowsAffected()
	if err1 != nil {
		fmt.Println("rows failed, ", err1)
		return false
	}

	fmt.Println("delete success: ", row)
	return true
}
