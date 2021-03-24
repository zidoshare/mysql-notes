# 如何设置主从同步

主备启动后，分别创建数据库demo并导入demo.sql

> source demo.sql

# 主库

```
#创建slave账号account，密码123456
mysql>grant replication slave on *.* to 'account'@'0.0.0.0' identified by '123456';
#更新数据库权限
mysql>flush privileges;
```

查看主服务器状态：

```
mysql>show master status\G;
***************** 1. row ****************
            File: mysql-bin.000033 #当前记录的日志
        Position: 337523 #日志中记录的位置
    Binlog_Do_DB:
Binlog_Ignore_DB:
```

# 从库

执行同步命令：

```
#执行同步命令，设置主服务器ip，同步账号密码，同步位置
mysql>change master to master_host='demo_master_mysql',master_user='account',master_password='123456',master_log_file='mysql-bin.000033',master_log_pos=337523;
#开启同步功能
mysql>start slave;
```

查看同步状态：

```
mysql> show slave status\G;
*************************** 1. row ***************************
               Slave_IO_State: Connecting to master
                  Master_Host: demo_master_mysql
                  Master_User: account
                  Master_Port: 3306
                Connect_Retry: 60
              Master_Log_File: mysql-bin.000033
          Read_Master_Log_Pos: 337523
               Relay_Log_File: 7bf8b44b32c6-relay-bin.000001
                Relay_Log_Pos: 4
        Relay_Master_Log_File: mysql-bin.000033
             Slave_IO_Running: Connecting
            Slave_SQL_Running: Yes
              Replicate_Do_DB: demo
          Replicate_Ignore_DB: 
           Replicate_Do_Table: 
       Replicate_Ignore_Table: 
      Replicate_Wild_Do_Table: 
  Replicate_Wild_Ignore_Table: 
                   Last_Errno: 0
                   Last_Error: 
                 Skip_Counter: 0
          Exec_Master_Log_Pos: 337523
              Relay_Log_Space: 154
              Until_Condition: None
               Until_Log_File: 
                Until_Log_Pos: 0
           Master_SSL_Allowed: No
           Master_SSL_CA_File: 
           Master_SSL_CA_Path: 
              Master_SSL_Cert: 
            Master_SSL_Cipher: 
               Master_SSL_Key: 
        Seconds_Behind_Master: NULL
Master_SSL_Verify_Server_Cert: No
                Last_IO_Errno: 1045
                Last_IO_Error: error connecting to master 'account@demo-mysql:3306' - retry-time: 60  retries: 1
               Last_SQL_Errno: 0
               Last_SQL_Error: 
  Replicate_Ignore_Server_Ids: 
             Master_Server_Id: 0
                  Master_UUID: 
             Master_Info_File: /var/lib/mysql/master.info
                    SQL_Delay: 0
          SQL_Remaining_Delay: NULL
      Slave_SQL_Running_State: Slave has read all relay log; waiting for more updates
           Master_Retry_Count: 86400
                  Master_Bind: 
      Last_IO_Error_Timestamp: 210324 09:46:01
     Last_SQL_Error_Timestamp: 
               Master_SSL_Crl: 
           Master_SSL_Crlpath: 
           Retrieved_Gtid_Set: 
            Executed_Gtid_Set: 
                Auto_Position: 0
         Replicate_Rewrite_DB: 
                 Channel_Name: 
           Master_TLS_Version: 
1 row in set (0.00 sec)

```