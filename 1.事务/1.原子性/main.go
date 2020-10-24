package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zidoshare/desc"
)

func main() {

	//初始化数据
	db, err := sql.Open("mysql", "root:123456@tcp(mysql:3306)/")
	defer func() {
		_, err2 := db.Exec("drop database test_2_1")
		if err2 != nil {
			panic(err2)
		}
	}()
	db.SetMaxOpenConns(1)
	if err != nil {
		panic(err)
	}
	//等待数据库启动完成
	desc.WaitingDb(db)

	initSQL, err := ioutil.ReadFile("./data.sql")
	if err != nil {
		panic(err)
	}
	requests := strings.Split(string(initSQL), ";")
	for _, request := range requests {
		if request == "" {
			continue
		}
		_, err := db.Exec(request)
		if err != nil {
			panic(err)
		}
		// do whatever you need with result and error
	}
	db.Exec("use test")

	//代码结束清除数据
	defer db.Query("drop database test")

	//开始事务逻辑
	fmt.Print("1.")
	desc.PrintfForQuery(db, "select * from t")
	fmt.Println("2.开启事务")
	tx, err := db.Begin()
	defer func() {
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				panic(err2)
			}
		}
	}()
	fmt.Print("3.")
	desc.PrintfForExec(tx, "insert into t value (2)")
	fmt.Print("4.")
	desc.PrintfForQuery(tx, "select * from t")
	fmt.Print("5.")
	desc.PrintfForExec(tx, "insert into t value (3)")
	fmt.Print("6.")
	desc.PrintfForQuery(tx, "select * from t")
	fmt.Println("7.提交事务")
	tx.Commit()
	fmt.Print("8.")
	desc.PrintfForQuery(db, "select * from t")
	fmt.Println("9.开启事务")
	tx, err = db.Begin()
	defer func() {
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				panic(err2)
			}
		}
	}()
	fmt.Print("10.")
	desc.PrintfForExec(tx, "insert into t value (2)")
	fmt.Print("11.")
	desc.PrintfForQuery(tx, "select * from t")
	fmt.Print("12.")
	desc.PrintfForExec(tx, "insert into t value (3)")
	fmt.Print("13.")
	desc.PrintfForQuery(tx, "select * from t")
	fmt.Println("14.回滚事务")
	tx.Rollback()
	fmt.Print("15.")
	desc.PrintfForQuery(db, "select * from t")
}
