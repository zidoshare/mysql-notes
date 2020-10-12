# Mysql 学习笔记

* 致力于每一步都能够回放

* 致力于每一步都有迹可循

* 致力于最清晰

### 复现

**若无特殊说明，命令及代码可在windows 10（以下不保证）、linux、Macos下运行。**

开发语言使用**go**，使用golang仅仅为了能更好的复现场景，golang环境不必备，只需要保证安装**docker**即可。

若无特殊说明，复现步骤可通过此命令一键执行（cd到该目录下）：`docker-compose up --abort-on-container-exit --build; docker-compose rm -f`

### 为什么使用golang？

> 因为懒

# 目录

* [1.事务](./1.事务)
  * [1.1.原子性](./1.事务/1.原子性)
  * [1.2.隔离性](./1.事务/2.隔离性)
    * [1.2.1.读未提交执行报告](./1.事务/2.隔离性/1.读未提交)
    * [1.2.2.读已提交执行报告](./1.事务/2.隔离性/2.读已提交)
    * [1.2.3.可重复读执行报告](./1.事务/2.隔离性/3.可重复读)
    * [1.2.4.串行化执行报告](./1.事务/2.隔离性/4.串行化)