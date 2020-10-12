# 原子性

### 复现步骤

* 启动mysql

* 准备最简单的建库建表语句
```sql
create database test;
use test;
create table t(c int) engine=InnoDB;
insert into t(c) values(1);
```
* 导入数据
```sql
drop database if exists test;
create database test;
use test;
create table t(c int) engine=InnoDB;
insert into t(c) values(1);
```
* 连接：`mysql -uroot -p123456 test`
* 执行
1. 查询数据：`select * from t;`
2. 开启事务:`begin;`
3. 插入一条数据： `insert into t value  (2)`
4. 查询数据：`select * from t;`
5. 插入一条数据： `insert into t value (3)`
6. 查询数据：`select * from t;`
7. 提交事务：`commit;`
8. 查询数据：`select * from t;`
9. 开启事务:`begin;`
10. 插入一条数据： `insert into t value (2)`
11. 查询数据：`select * from t;`
12. 插入一条数据： `insert into t value (3)`
13. 查询数据：`select * from t;`
14. 回滚事务：`rollback;`
15. 查询数据：`select * from t;`


执行期望： 首先查看原数据，接着在事务中尝试添加两条数据，之后提交，能够查到事务中添加的数据，接着再开启事务，插入两条数据并回滚，此时事务中执行的插入并不生效。

### 执行结果

1.执行sql:select * from t
+---+
| C |
+---+
| 1 |
+---+
2.开启事务
3.执行sql:insert into t value (2) 执行成功,影响了1行
4.执行sql:select * from t
+---+
| C |
+---+
| 1 |
| 2 |
+---+
5.执行sql:insert into t value (3) 执行成功,影响了1行
6.执行sql:select * from t
+---+
| C |
+---+
| 1 |
| 2 |
| 3 |
+---+
7.提交事务
8.执行sql:select * from t
+---+
| C |
+---+
| 1 |
| 2 |
| 3 |
+---+
9.开启事务
10.执行sql:insert into t value (2) 执行成功,影响了1行
11.执行sql:select * from t
+---+
| C |
+---+
| 1 |
| 2 |
| 3 |
| 2 |
+---+
12.执行sql:insert into t value (3) 执行成功,影响了1行
13.执行sql:select * from t
+---+
| C |
+---+
| 1 |
| 2 |
| 3 |
| 2 |
| 3 |
+---+
14.回滚事务
15.执行sql:select * from t
+---+
| C |
+---+
| 1 |
| 2 |
| 3 |
+---+