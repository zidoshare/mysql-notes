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
			if time.Now().Sub(crt) > time.Second*30 {
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
	db, err := sql.Open("mysql", "root:123456@tcp(mysql:3306)/")
	db.SetMaxOpenConns(1)
	panicWhenError(err)
	//等待数据库启动完成
	waitingDb(db)

	initSQL, err := ioutil.ReadFile("./data.sql")
	panicWhenError(err)
	requests := strings.Split(string(initSQL), ";")
	for _, request := range requests {
		if request == "" {
			continue
		}
		_, err := db.Exec(request)
		panicWhenError(err)
		// do whatever you need with result and error
	}
	db.Exec("use test")
	defer db.Close()

	//代码结束清除数据
	defer db.Query("drop database test")

	//开始事务逻辑
	fmt.Print("1.")
	printfForQuery(db, "select * from t")
	fmt.Println("2.开启事务")
	tx, err := db.Begin()
	panicWhenError(err)
	fmt.Print("3.")
	printfForExec(tx, "insert into t value (2)")
	fmt.Print("4.")
	printfForQuery(tx, "select * from t")
	fmt.Print("5.")
	printfForExec(tx, "insert into t value (3)")
	fmt.Print("6.")
	printfForQuery(tx, "select * from t")
	fmt.Println("7.提交事务")
	tx.Commit()
	fmt.Print("8.")
	printfForQuery(db, "select * from t")
	fmt.Println("9.开启事务")
	tx, err = db.Begin()
	panicWhenError(err)
	fmt.Print("10.")
	printfForExec(tx, "insert into t value (2)")
	fmt.Print("11.")
	printfForQuery(tx, "select * from t")
	fmt.Print("12.")
	printfForExec(tx, "insert into t value (3)")
	fmt.Print("13.")
	printfForQuery(tx, "select * from t")
	fmt.Println("14.回滚事务")
	tx.Rollback()
	fmt.Print("15.")
	printfForQuery(db, "select * from t")
}
