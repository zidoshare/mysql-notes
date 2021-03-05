串行化时

执行结果如下：

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
5.cmd2执行sql:update t set c = 2 where c = 1
panic: Error 1205: Lock wait timeout exceeded; try restarting transaction
```

与可重复读相比，cmd2 执行 update t set c = 2 where c = 1 时会死锁
