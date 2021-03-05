可重复读的场景下

执行结果为：

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
| 1 |
+---+
7.cmd2 提交事务:
8.cmd1执行sql:select * from t
+---+
| C |
+---+
| 1 |
+---+
9.cmd1 提交事务:
10.cmd1执行sql:select * from t
+---+
| C |
+---+
| 2 |
+---+
```

注意这里与读已提交相比，第 8 步仍然能读到 `1`，即使是 cmd2 的修改事务已经提交。

也就是说，可重复读场景下，即使 cmd2 提交了数据，但是 cmd1 仍然能够在没有提交事务的情况下保证读取的一致性
