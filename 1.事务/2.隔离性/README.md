# 隔离性

### 复现步骤

* 启动 MySQL
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

* 使用两个控制台开启两个连接,**使用 cmd1、cmd2 分别代表这两个终端执行**
* 设置隔离级别
* cmd1 连接：`mysql -uroot -p123456 test`
* cmd2 连接：`mysql -uroot -p123456 test`
* 执行图解：
  ![网图](./%E6%89%A7%E8%A1%8C%E9%A1%BA%E5%BA%8F.png)

1. cmd1 开启事务:`begin;`
2. cmd2 开启事务:`begin;`
3. cmd1 查询数据： `select * from t;`
4. cmd2 查询数据： `select * from t;`
5. cmd2 更改数据： `update t set c = 2 from t where c = 1;`
6. cmd1 查询数据： `select * from;`
7. cmd2 提交事务：`commit;`
8. cmd1 查询数据：`select * from t;`
9. cmd1 提交事务：`commit;`
10. cmd1 查询数据：`select * from t;`
    具体执行可以参考执行报告及同目录下的代码

* [读未提交](./1.%E8%AF%BB%E6%9C%AA%E6%8F%90%E4%BA%A4)：事务 1 可以读取到事务 2 修改过但未提交的数据（产生脏读，幻读，不可重复度）
* [读已提交](./2.%E8%AF%BB%E5%B7%B2%E6%8F%90%E4%BA%A4)：事务 1 只能在事务 2 修改过并且已提交后才能读取到事务 2 修改的数据（产生幻读，不可重复度）
* [可重复读](./3.%E5%8F%AF%E9%87%8D%E5%A4%8D%E8%AF%BB)：事务 1 只能在事务 2 修改过数据并提交后，自己也提交事务后，才能读取到事务 2 修改的数据(产生幻读)
* [串行化](./4.%E4%B8%B2%E8%A1%8C%E5%8C%96)：事务 1 在执行过程中完全看不到事务 2 对数据库所做的更新。当两个事务同时操作数据库中相同数据时，如果事务 1 已经在访问该数据，事务 2 只能停下来等待，必须等到事务 1 结束后才能恢复运行

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
