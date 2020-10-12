# 隔离性

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
* 使用两个控制台开启两个连接,**使用cmd1、cmd2分别代表这两个终端执行**
* 设置隔离级别
* cmd1连接：`mysql -uroot -p123456 test`
* cmd2连接：`mysql -uroot -p123456 test`
* 执行图解：
![网图](./执行顺序.png)
1. cmd1开启事务:`begin;`
2. cmd2开启事务:`begin;`
3. cmd1查询数据： `select * from t;`
4. cmd2查询数据： `select * from t;`
5. cmd2更改数据： `update t set c = 2 from t where c = 1;`
6. cmd1查询数据： `select * from;`
7. cmd2提交事务：`commit;`
8. cmd1查询数据：`select * from t;`
9. cmd1提交事务：`commit;`
10. cmd1查询数据：`select * from t;`
具体执行可以参考执行报告及同目录下的代码
* [读未提交](./1.读未提交)：事务1可以读取到事务2修改过但未提交的数据（产生脏读，幻读，不可重复度）
* [读已提交](./2.读已提交)：事务1只能在事务2修改过并且已提交后才能读取到事务2修改的数据（产生幻读，不可重复度）
* [可重复读](./3.可重复读)：事务1只能在事务2修改过数据并提交后，自己也提交事务后，才能读取到事务2修改的数据(产生幻读)
* [串行化](./4.串行化)：事务1在执行过程中完全看不到事务2对数据库所做的更新。当两个事务同时操作数据库中相同数据时，如果事务1已经在访问该数据，事务2只能停下来等待，必须等到事务1结束后才能恢复运行
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
