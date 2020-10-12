读未提交是最低的事务隔离级别

执行结果如下

```
1.cmd1 开启事务
2.cmd1执行sql:select * from t
+---+
| C |
+---+
| 1 |
+---+
3.cmd2 开启事务:
4.cmd2执行sql:select * from t
+---+
| C |
+---+
| 1 |
+---+
5.cmd2执行sql:update t set c = 2 where c = 1 执行成功,影响了1行
6.cmd1执行sql:select * from t
+---+
| C |
+---+
| 2 |
+---+
7.cmd2 提交事务:
8.cmd1执行sql:select * from t
+---+
| C |
+---+
| 2 |
+---+
9.cmd1 提交事务:
10.cmd1执行sql:select * from t
+---+
| C |
+---+
| 2 |
+---+
```
