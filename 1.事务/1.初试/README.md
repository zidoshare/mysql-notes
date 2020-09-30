# 初试事务

准备最简单的建库建表语句

```sql
create database test;
use test;
create table t(c int) engine=InnoDB;
insert into t(c) values(1);
```

启动mysql: `docker-compose up -d`

导入数据文件：`cat .\data.sql | docker exec -i test_mysql mysql -uroot -p123456`

### 事务隔离级别相关操作说明

查看事务隔离级别：

` select @@global.tx_isolation,@@tx_isolation;`

默认隔离级别为：可重复读
<pre>
+-----------------------+-----------------+
| @@global.tx_isolation | @@tx_isolation  |
+-----------------------+-----------------+
| REPEATABLE-READ       | REPEATABLE-READ |
+-----------------------+-----------------+
</pre>

设置不同的事务隔离级别：

`set global transaction isolation level read committed;` //全局设置

`set session transaction isolation level read committed;` //当前会话

四个不同的隔离级别为：
* `Read Uncommitted`:读未提交
* `Read Committed`: 读已提交
* `Repeatable Read`: 可重复读
* `Serializable`: 串行化

### 控制台相关说明：

使用两个控制台开启两个连接:

`docker exec -it test_mysql mysql -uroot -p123456`

使用cmd1、cmd2分别代表这两个终端执行

### 操作说明：

执行图解：
![网图](./执行顺序.png)

1. 设置全局隔离级别：`set global transaction isolation level 隔离级别;`
2. cmd1连接：`docker exec -it test_mysql mysql -uroot -p123456`
3. cmd2连接：`docker exec -it test_mysql mysql -uroot -p123456`
4. cmd1开启事务:`begin;`
5. cmd2开启事务:`begin;`
6. cmd1查询数据： `select * from t;`
7. cmd2查询数据： `select * from t;`
8. cmd2更改数据： `update t set c = 2 from t where c = 1;`
9. cmd1查询数据： `select * from;`
10. cmd2提交事务：`commit;`
11. cmd1查询数据：`select * from t;`
12. cmd1提交事务：`commit;`
13. cmd1查询数据：`select * from t;` 

具体执行可以参考执行报告及同目录下的代码
* [读未提交](./1.读未提交)
* [读已提交](./2.读已提交)
* [可重复读](./3.可重复读)
* [串行化](./4.串行化)