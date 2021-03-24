# MySQL 基础架构

MySQL 的学习应该先看全貌，了解 MySQL 的基础架构：

![](./%E9%80%BB%E8%BE%91%E6%9E%B6%E6%9E%84%E5%9B%BE.png)

MySQL 架构大体上分为 Server 层和存储引擎层两部分。

Server 层涵盖 MySQL 的大多数核心服务功能，以及所有的内置函数（如日期、时间、数学和加密函数等），所有跨存储引擎的功能都在这一层实现，比如存储过程、触发器、视图等。Server 层包括：

* 连接器：负责跟客户端建立连接、获取权限、维持和管理连接 show processlist; Sleep: 连接完成后，如果没有后续动作，这个连接就处于空闲状态,客户端如果太长时间没动静，连接器就会自动将它断开。这个时间是由参数 wait_timeout 控制的，默认值 8 小时 MySQL 在执行过程中临时使用的内存是管理在连接对象里面的。这些资源会在连接断开的时候才释放。 所以如果长连接累积下来，可能导致内存占用太大，被系统强行杀掉（OOM），从现象看就是 MySQL 异常重启了。
* 查询缓存：8.0 版本已将该模块下线，因为缓存失效非常频繁，如修改某张表，表上的所有查询都会失效；因此，除非业务实例只有一张静态表，不然不建议开启查询缓存 query_cache_type
* 分析器： 如果查询语句没有在查询缓存中命中，就开始到了分析器
  * 词法分析：分析 SQL 中每个单词代表什
是基于磁盘的数据结构，在崩溃恢复期间用于纠正不完整事务写入的数据。 在正常操作期间，重做日志对更改表数据的请求进行编码，这些请求是由SQL语句或低级API调用产生的。 在初始化过程中以及接受连接之前，将自动重播在意外关闭之前未完成更新数据文件的修改。

默认情况下，重做日志在磁盘上由名为ib_logfile0和ib_logfile1的两个文件物理表示。

InnoDB 的 redo log 从数据结构上来看是固定大小的循环链表，如果写满了则要擦除最开始记录的数据。使用双指针遍历，**write pos**为当前记录的位置，一边写一边往后移。**checkpoint**为当前要擦除的位置，往后推移，当擦出前要把数据写进磁盘。

# 最后

Q: 如果表 T 中没有字段 k，而你执行了这个语句 select * from T where k=1, 那肯定是会报“不存在这个列”的错误： “Unknown column ‘k’ in ‘where clause’”。你觉得这个错误是在我们上面提到的哪个阶段报出来的呢？

A: 分析器，它对语句中的表、字段是否存在进行判断。

## 连接器

连接器负责跟客户端建立连接、获取权限、维持和管理连接。

连接命令一般如下：

```
mysql -h$ip -P$port -u$user -p
```

1. 与服务器完成 TCP 握手
2. 连接器使用输入的用户名密码进行身份认证
3. 如果用户名密码通过，到权限表查出拥有的权限，并缓存下来，此后**这个连接**里的所有权限判断逻辑都依赖语此时读到的权限，即使管理员对权限做了更改也不会刷新，除非再新建连接
4. 连接成功后处理空闲状态，可以使用 `show processlist` 命令看到 `Command` 列为 `Sleep`，表示现在系统里面有一个空闲连接。如果客户端太长时间没东京，连接器会自动断开，参数由 `wait_timeout` 控制，默认 8 小时

数据库里面，长连接是指连接成功后，如果客户端持续有请求，则一直使用同一个连接。短连接则是指每次执行完很少的几次查询就断开连接，下次查询再重新建立一个。

建议尽量使用长连接，因为建立连接的过程比较复杂

当全部使用长连接后，MySQL 占用内存涨得特别快，因为 MySQL 在执行过程中临时使用的内存时管理再连接对象里面的。这些资源会在连接断开的时候才释放。所以长时间累积下来，可能导致内存占用太大，被系统强行杀掉（OOM），从现象上看就是 MySQL 异常重启了。

解决方案：

1. 定期断开长连接，使用一段时间后，或者程序判断执行过一个大查询后，断开连接，之后需要时再重连
2. 如果是 MySQL5.7 及以上版本，再执行一个比较大的操作后，执行 `mysql_reset_connection` 重新初始化连接资源，这个过程不需要重连和权限校验，但是连接会恢复到刚刚创建的状态

## 查询缓存

MySQL 拿到一个查询请求后，先到查询缓存中查找之前是否执行过这条语句。之前执行的结果会直接缓存在内存中。

缺点：

* 缓存的 Key 为查询的语句，也就是说 SQL 语句和参数必须完全相同
* 失效非常频繁，只要对一个表更新，那么所有的查询缓存都会被清空。缓存命中率非常低

可以按需开启查询缓存，但是 MySQL8.0 完全删除了此功能

## 分析器

如果没有命中查询缓存，就开始真正执行语句。这需要先对 SQL 语句进行解析。

此阶段先做词法分析识别表名、列名等基本信息。后做语法分析，判断 SQL 语句是否合法。

此时若语法错误，会报：

> you have an error in your SQL syntax

## 优化器

优化器是再表里面有多个索引的时候决定使用哪个索引，或者有多表关联时决定连接顺序。从而确定执行方法。

## 执行器

MySQL 通过分析器分析了语句目的，优化器确定该如何做，之后进入执行器开始执行。

开始执行的时候，先**判断对表 T 是否有相关权限**，如果没有，则返回权限错误。

> 权限验证不仅仅在执行器这部分会做，在分析器之后，也就是知道了该语句要“干什么”之后，也会先做一次权限验证。叫做 precheck。而 precheck 是无法对运行时涉及到的表进行权限验证的，比如使用了触发器的情况。因此在执行器这里也要做一次执行时的权限验证。

> MySQL 的权限验证阶段，在执行器执行环节为什么还要进行验证，因为除了 SQL 还可能有存储引擎，触发器等，在这些对象中，也可能需要调用其它表去获取数据，也需要权限验证，如果在分析节点之前就进行对象的验证，那么对于触发器，存储引擎这种对象的执行是做不到的。

接着执行器打开表调用的存储引擎的 API 来操作数据的。也就是说此时才会用上存储引擎，从而实现插件化。

1. 调用 InnoDB 引擎接口取这个表的第一行，判断 ID 值是不是 10，如果不是则跳过，如果是则将这行存在结果集中；
2. 调用引擎接口取“下一行”，重复相同的判断逻辑，直到取到这个表的最后一行。
3. 执行器将上述遍历过程中所有满足条件的行组成的记录集作为结果集返回给客户端。

而对于有索引的表，就执行"取满足条件的下一行"，直到遍历结束。

# 更新语句的执行

执行更新语句跟查询语句一样会将上述流程走一次，但是其中会设计到两个重要的日志模块**redo log**和**binlog**

## redo log

redo log是InnoDB引擎特有的日志，在MySQL中如果每一次的更新操作都需要写进磁盘，然后磁盘也要找到对应的那条记录，然后再更新，整个过程 IO 成本、查找成本都很高。MySQL 使用 WAL 技术来解决这个问题，WAL 的全称是 Write-Ahead Logging，它的关键点就是先写日志，再写磁盘。

具体来说，当有一条记录需要更新的时候，InnoDB 引擎就会先把记录写到 redo log里面，并更新内存，这个时候更新就算完成了。同时，InnoDB 引擎会在适当的时候，将这个操作记录更新到磁盘里面，而这个更新往往是在系统比较空闲的时候做。

是基于磁盘的数据结构，在崩溃恢复期间用于纠正不完整事务写入的数据。 在正常操作期间，重做日志对更改表数据的请求进行编码，这些请求是由SQL语句或低级API调用产生的。 在初始化过程中以及接受连接之前，将自动重播在意外关闭之前未完成更新数据文件的修改。

默认情况下，重做日志在磁盘上由名为ib_logfile0和ib_logfile1的两个文件物理表示。

InnoDB 的 redo log 从数据结构上来看是固定大小的循环链表，如果写满了则要擦除最开始记录的数据。使用双指针遍历，**write pos**为当前记录的位置，一边写一边往后移。**checkpoint**为当前要擦出的位置，往后推移，当擦出前要把数据写进磁盘。

有了 redo log，InnoDB 就可以保证即使数据库发生异常重启，之前提交的记录都不会丢失，这个能力称为 crash-safe。因此，如果MySQL异常退出，仅需要重新启动MySQL即可。InnoDB自动检查日志并执行数据库到当前的前滚。 InnoDB自动回退崩溃时存在的未提交的事务。

在恢复期间的日志示例如下：

```
InnoDB: Log scan progressed past the checkpoint lsn 369163704
InnoDB: Doing recovery: scanned up to log sequence number 374340608
InnoDB: Doing recovery: scanned up to log sequence number 379583488
InnoDB: Doing recovery: scanned up to log sequence number 384826368
InnoDB: Doing recovery: scanned up to log sequence number 390069248
InnoDB: Doing recovery: scanned up to log sequence number 395312128
InnoDB: Doing recovery: scanned up to log sequence number 400555008
InnoDB: Doing recovery: scanned up to log sequence number 405797888
InnoDB: Doing recovery: scanned up to log sequence number 411040768
InnoDB: Doing recovery: scanned up to log sequence number 414724794
InnoDB: Database was not shutdown normally!
InnoDB: Starting crash recovery.
InnoDB: 1 transaction(s) which must be rolled back or cleaned up in
total 518425 row operations to undo
InnoDB: Trx id counter is 1792
InnoDB: Starting an apply batch of log records to the database...
InnoDB: Progress in percent: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15
16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37
38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59
60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81
82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99
InnoDB: Apply batch completed
...
InnoDB: Starting in background the rollback of uncommitted transactions
InnoDB: Rolling back trx with id 1511, 518425 rows to undo
...
InnoDB: Waiting for purge to start
InnoDB: 5.7.18 started; log sequence number 414724794
...
./mysqld: ready for connections.
```
## binlog

在MySQL的Server层也有自己的日志，称为binlog。

binlog不具有CrashSafe的能力，因为数据落盘成功与否是引擎层保证的。而binlog主要用于归档。他们的区别主要是：

* redo log 是 InnoDB 引擎特有的；binlog 是 MySQL 的 Server 层实现的，所有引擎都可以使用。
* redo log 是物理日志，记录的是“在某个数据页上做了什么修改”。binlog是逻辑日志，记录的是这个语句的原始逻辑，比如“给 ID=2 这一行的 c 字段加 1 ”。
* redo log 是循环写的，空间固定会用完。binlog 是可以追加写入的。“追加写”是指 binlog 文件写到一定大小后会切换到下一个，并不会覆盖以前的日志。

### binlog与redo log的执行流程

当数据更新后，首先写入redo log，此时redo log处于prepare状态，接着，写入binlog，之后提交事务后再对redo log进行commit，这被称之为**两阶段提交 **

例如对于更新语句：`update T set c=c+1 where ID=2;`来说

* 执行器先找引擎取 ID=2 这一行。ID 是主键，引擎直接用树搜索找到这一行。如果 ID=2 这一行所在的数据页本来就在内存中，就直接返回给执行器；否则，需要先从磁盘读入内存，然后再返回。
* 执行器拿到引擎给的行数据，把这个值加上 1，比如原来是 N，现在就是 N+1，得到新的一行数据，再调用引擎接口写入这行新数据。
* 引擎将这行新数据更新到内存中，同时将这个更新操作记录到 redo log 里面，此时 redo log 处于 prepare 状态。然后告知执行器执行完成了，随时可以提交事务。
* 执行器生成这个操作的 binlog，并把 binlog 写入磁盘。
* 执行器调用引擎的提交事务接口，引擎把刚刚写入的 redo log 改成提交（commit）状态，更新完成。

执行流程图如图所示：
![更新语句执行流程](./更新语句执行流程.png)

# 最后

Q: 如果表 T 中没有字段 k，而你执行了这个语句 select * from T where k=1, 那肯定是会报“不存在这个列”的错误： “Unknown column ‘k’ in ‘where clause’”。你觉得这个错误是在我们上面提到的哪个阶段报出来的呢？

A: 分析器，它对语句中的表、字段是否存在进行判断。
