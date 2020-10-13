package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
)

//ErrNotDbOrTx type error
var ErrNotDbOrTx = errors.New("db必须为sql.DB或sql.Tx类型")

func waitingDb(db *sql.DB) {
	crt := time.Now()

	for {
		err := db.Ping()
		if err != nil {
			if time.Now().Sub(crt) > time.Second*30 {
				panic(err)
			}
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

func printfForQuery(db interface{}, queryStr string) error {
	fmt.Println("执行sql:" + queryStr)
	var rows *sql.Rows
	var err error
	switch db.(type) {
	case *sql.DB:
		rows, err = db.(*sql.DB).Query(queryStr)
		if err != nil {
			return err
		}
		defer rows.Close()
	case *sql.Tx:
		rows, err = db.(*sql.Tx).Query(queryStr)
		if err != nil {
			return err
		}
		defer rows.Close()
	default:
		return ErrNotDbOrTx
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)
	data := make([][]string, 0)

	for rows.Next() {
		row := make([]interface{}, len(columns))
		container := make([]string, len(columns))
		for i := range row {
			row[i] = &container[i]
		}
		rows.Scan(row...)
		data = append(data, container)
	}
	table.AppendBulk(data)
	table.Render()
	return nil
}

func printfForExec(db interface{}, queryStr string) error {
	fmt.Print("执行sql:" + queryStr)
	var result sql.Result
	var err error
	switch db.(type) {
	case *sql.DB:
		result, err = db.(*sql.DB).Exec(queryStr)
		if err != nil {
			return err
		}
	case *sql.Tx:
		result, err = db.(*sql.Tx).Exec(queryStr)
		if err != nil {
			return err
		}
	default:
		return ErrNotDbOrTx
	}
	rows, err := result.RowsAffected()
	if err != nil {
		switch db.(type) {
		case sql.Tx:
			db.(*sql.Tx).Rollback()
		}
		panic(err)
	}
	if rows > 0 {
		fmt.Printf(" 执行成功,影响了%d行\n", rows)
	} else {
		fmt.Println(" 执行失败,没有更改任何数据")
	}
	if err != nil {
		return err
	}
	return nil
}

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
	waitingDb(db1)
	//代码结束清除数据
	defer func() {
		_, err2 := db1.Exec("drop database test_4")
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
	db1.Exec("use test_4")
	db1.Exec("SET GLOBAL innodb_lock_wait_timeout=3;")

	//设置连接1事务隔离级别
	_, err = db1.Exec("set session transaction isolation level serializable")
	if err != nil {
		panic(err)
	}

	//连接2
	db2, err := sql.Open("mysql", "root:123456@tcp(mysql:3306)/test_4")
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
	_, err = db2.Exec("set session transaction isolation level serializable")
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
	err = printfForQuery(tx1, "select * from t")
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
	err = printfForQuery(tx2, "select * from t")
	if err != nil {
		panic(err)
	}
	//这里会被锁住，因为cmd1还未提交。
	fmt.Print("5.cmd2")
	err = printfForExec(tx2, "update t set c = 2 where c = 1")
	if err != nil {
		panic(err)
	}
	fmt.Print("6.cmd1")
	err = printfForQuery(tx1, "select * from t")
	if err != nil {
		panic(err)
	}
	fmt.Println("7.cmd2 提交事务:")
	err = tx2.Commit()
	if err != nil {
		panic(err)
	}
	fmt.Print("8.cmd1")
	err = printfForQuery(tx1, "select * from t")
	if err != nil {
		panic(err)
	}
	fmt.Println("9.cmd1 提交事务:")
	err = tx1.Commit()
	if err != nil {
		panic(err)
	}
	fmt.Print("10.cmd1")
	err = printfForQuery(db1, "select * from t")
	if err != nil {
		panic(err)
	}
}
