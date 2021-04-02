# 如何修改 MySQL redo log 大小

参阅 MySQL 文档中的 [Changing the Number or Size of InnoDB Redo Log Files](https://dev.mysql.com/doc/refman/5.7/en/innodb-redo-log.html)

**binlog 是磁盘数据结构**

确保 MySQL 处于正常关闭状态。

修改配置文件中的 **innodb_log_file_size**，用于配置每个 redo log 文件大小。

* 命令行：	--innodb-log-file-size=#
* 环境变量： innodb_log 各 it 各 it 各 it_file_size
* 作用范围： 全局
* 是否可动态修改： 否
* 类型:	Integer
* 默认值:	50331648
* 最小值(≥ 5.7.11)	4194304
* 最小值 (≤ 5.7.10)	1048576
* 最大值	512GB / innodb_log_files_in_group

日志组中每个日志文件的大小（以字节为单位）。 日志文件的总大小（innodb_log_file_size * innodb_log_files_in_group）不能超过略小于 512GB 的最大值。 例如，一对 255 GB 的日志文件已达到限制，但没有超过该限制。 默认值为 48MB。

通常，日志文件的总大小应足够大，以使服务器可以消除工作负载活动中的高峰和低谷，这通常意味着有足够的重做日志空间来处理一个小时以上的写活动。 值越大，缓冲池中需要的检查点刷新活动越少，从而节省了磁盘 I / O。 较大的日志文件也会使崩溃恢复变慢。

在 MySQL 5.7.11 中，最小的 innodb_log_file_size 值从 1MB 增加到 4MB。

修改 innodb_log_files_in_group 用于配置有多少个组：

* 命令行：	--innodb-log-file-size=#
* 环境变量： innodb_log_file_size
* 作用范围： 全局
* 是否可动态修改： 否
* 类型：Integer
* 默认值：2
* 最小值：2
* 最大值：100

# 计算每分钟 redo log 量

```
mysql>  pager grep -i "Log sequence number";
PAGER set to 'grep -i "Log sequence number"'
mysql> show engine innodb status \G select sleep(60); show engine innodb status \G
Log sequence number 4951647
1 row in set (0.01 sec)
1 row in set (1 min 0.00 sec)
Log sequence number 5046150
1 row in set (0.00 sec)
```

在这 60s 期间，我们业务系统处于正常的运行状态，此次为实验环境，我做了简单的业务模拟操作。

lsn 号从 4951647 增长到 5046150

一分钟 redo log 量：round((5046150-4951647)/1024)=92KB

一小时 redo log 量：92K x 60=5520KB

正常来讲，数据库 10 分钟切换一次 redo log，故对于此数据库，单个 redo log 大小 5520KB /6=920KB

由于此数据库为测试平台，业务量较小，正常来讲生产库的单个 redo log 大小在 200M-2G 之间

# 修改后结果

```
mysql> show variables like 'innodb_log%';
+-----------------------------+-----------+
| Variable_name               | Value     |
+-----------------------------+-----------+
| innodb_log_buffer_size      | 16777216  |
| innodb_log_checksums        | ON        |
| innodb_log_compressed_pages | ON        |
| innodb_log_file_size        | 209715200 |
| innodb_log_files_in_group   | 4         |
| innodb_log_group_home_dir   | ./        |
| innodb_log_write_ahead_size | 8192      |
+-----------------------------+-----------+
```

# 参考连接

[mysql 官方文档](https://dev.mysql.com/doc/refman/5.7/en/innodb-parameters.html#sysvar_innodb_log_file_size)
[MySQL 设置 redo log 大小](http://blog.itpub.net/30135314/viewspace-2222251/)
