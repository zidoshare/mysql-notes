package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
)

func panicWhenError(err error) {
	if err != nil {
		panic(err)
	}
}

func waitingDb(db *sql.DB) {
	crt := time.Now()

	for {
		err := db.Ping()
		if err != nil {
			if time.Now().Sub(crt) > time.Second*10 {
				panic(err)
			}
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

type qb interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func printfForQuery(db qb, queryStr string) {
	fmt.Println("执行sql:" + queryStr)
	rows, err := db.Query(queryStr)
	if rows != nil {
		defer rows.Close()
	}

	panicWhenError(err)
	columns, err := rows.Columns()
	panicWhenError(err)
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
}

func printfForExec(db qb, queryStr string) {
	fmt.Print("执行sql:" + queryStr)
	result, err := db.Exec(queryStr)
	panicWhenError(err)
	rows, err := result.RowsAffected()
	if rows > 0 {
		fmt.Printf(" 执行成功,影响了%d行\n", rows)
	} else {
		fmt.Println(" 执行失败,没有更改任何数据")
	}
	panicWhenError(err)
}

func main() {

	//初始化数据
	//连接1
	db1, err := sql.Open("mysql", "root:123456@tcp(mysql:3306)/")
	db1.SetMaxOpenConns(1)
	panicWhenError(err)
	//等待数据库启动完成
	waitingDb(db1)

	initSQL, err := ioutil.ReadFile("./data.sql")
	panicWhenError(err)
	requests := strings.Split(string(initSQL), ";")
	for _, request := range requests {
		if request == "" {
			continue
		}
		_, err := db1.Exec(request)
		panicWhenError(err)
		// do whatever you need with result and error
	}
	db1.Exec("use test")
	defer db1.Close()

	//代码结束清除数据
	defer db1.Query("drop database test")

	//设置连接1事务隔离级别
	_, err = db1.Exec("set session transaction isolation level read committed")
	panicWhenError(err)

	//连接2
	db2, err := sql.Open("mysql", "root:123456@tcp(mysql:3306)/test")

	panicWhenError(err)
	db2.SetMaxOpenConns(1)
	defer db2.Close()

	//设置连接2事务隔离级别
	_, err = db2.Exec("set session transaction isolation level read committed")
	panicWhenError(err)
	//开始事务逻辑
	fmt.Println("1.cmd1 开启事务")
	tx1, err := db1.Begin()
	panicWhenError(err)
	fmt.Print("2.cmd1")
	printfForQuery(tx1, "select * from t")
	fmt.Println("3.cmd2 开启事务:")
	tx2, err := db2.Begin()
	panicWhenError(err)
	fmt.Print("4.cmd2")
	printfForQuery(tx2, "select * from t")
	fmt.Print("5.cmd2")
	printfForExec(tx2, "update t set c = 2 where c = 1")
	fmt.Print("6.cmd1")
	printfForQuery(tx1, "select * from t")
	fmt.Println("7.cmd2 提交事务:")
	tx2.Commit()
	fmt.Print("8.cmd1")
	printfForQuery(tx1, "select * from t")
	fmt.Println("9.cmd1 提交事务:")
	tx1.Commit()
	fmt.Print("10.cmd1")
	printfForQuery(db1, "select * from t")
}
