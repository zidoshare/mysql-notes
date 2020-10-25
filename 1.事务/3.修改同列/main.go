package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zidoshare/desc"
)

func main() {
	pprofHandler := http.NewServeMux()
	pprofHandler.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	server := &http.Server{Addr: ":7890", Handler: pprofHandler}
	go server.ListenAndServe()
	//初始化数据
	//连接1
	db1, err := sql.Open("mysql", "root:123456@tcp(mysql:3306)/")
	if err != nil {
		panic(err)
	}
	defer func() {
		err = db1.Close()
		if err != nil {
			panic(err)
		}
	}()
	db1.SetMaxOpenConns(1)
	//等待数据库启动完成
	desc.WaitingDb(db1)
	//代码结束清除数据
	defer func() {
		_, err2 := db1.Exec("drop database test_3_3")
		if err2 != nil {
			panic(err2)
		}
	}()
	initSQL, err := ioutil.ReadFile("./data.sql")
	if err != nil {
		panic(err)
	}
	requests := strings.Split(string(initSQL), ";")
	for _, request := range requests {
		if request == "" {
			continue
		}
		_, err := db1.Exec(request)
		if err != nil {
			panic(err)
		}
		// do whatever you need with result and error
	}
	db1.Exec("use test_3_3")
	db1.Exec("SET GLOBAL innodb_lock_wait_timeout=3;")

	//连接2
	db2, err := sql.Open("mysql", "root:123456@tcp(mysql:3306)/test_3_3")
	if err != nil {
		panic(err)
	}
	db2.SetMaxOpenConns(1)

	defer func() {
		err2 := db2.Close()
		if err2 != nil {
			panic(err2)
		}
	}()
	//设置连接2事务隔离级别
	_, err = db2.Exec("set session transaction isolation level repeatable read")
	if err != nil {
		panic(err)
	}
	//开始事务逻辑
	fmt.Println("1.cmd1 开启事务")
	tx1, err := db1.Begin()
	if err != nil {
		panic(err)
	}
	fmt.Print("2.cmd1")
	err = desc.PrintfForQuery(tx1, "select * from t")
	defer func() {
		if err != nil {
			err2 := tx1.Rollback()
			if err2 != nil {
				panic(err2)
			}
		}
	}()
	if err != nil {
		panic(err)
	}
	fmt.Println("3.cmd2 开启事务:")
	tx2, err := db2.Begin()
	defer func() {
		if err != nil {
			err2 := tx2.Rollback()
			if err2 != nil {
				panic(err2)
			}
		}
	}()
	fmt.Print("4.cmd2")
	err = desc.PrintfForQuery(tx2, "select * from t")
	if err != nil {
		panic(err)
	}
	fmt.Print("5.cmd1")
	err = desc.PrintfForExec(tx1, "update t set c = 3 where c = 1")
	if err != nil {
		panic(err)
	}
	//这里会被锁住，因为cmd1还未提交。
	fmt.Print("6.cmd2")
	err = desc.PrintfForExec(tx2, "update t set c = 2 where c = 1")
	if err != nil {
		panic(err)
	}
}
