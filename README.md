# MySQL 学习笔记

* 致力于每一步都能够回放
* 致力于每一步都有迹可循
* 致力于最清晰

### 复现

**若无特殊说明，命令及代码可在 windows 10（以下不保证）、Linux、Macos 下运行。**

开发语言使用 **go**，使用 golang 仅仅为了能更好的复现场景，golang 环境不必备，只需要保证安装 **docker** 即可。

若无特殊说明，复现步骤可通过此命令一键执行（cd 到该目录下）：`docker-compose up --abort-on-container-exit --build; docker-compose rm -f`

### 为什么使用 golang？

> 因为懒

# 目录

* [1.事务](./1.%E4%BA%8B%E5%8A%A1)
  * [1.1.原子性](./1.%E4%BA%8B%E5%8A%A1/1.%E5%8E%9F%E5%AD%90%E6%80%A7)
  * [1.2.隔离性](./1.%E4%BA%8B%E5%8A%A1/2.%E9%9A%94%E7%A6%BB%E6%80%A7)
    * [1.2.1.读未提交执行报告](./1.%E4%BA%8B%E5%8A%A1/2.%E9%9A%94%E7%A6%BB%E6%80%A7/1.%E8%AF%BB%E6%9C%AA%E6%8F%90%E4%BA%A4)
    * [1.2.2.读已提交执行报告](./1.%E4%BA%8B%E5%8A%A1/2.%E9%9A%94%E7%A6%BB%E6%80%A7/2.%E8%AF%BB%E5%B7%B2%E6%8F%90%E4%BA%A4)
    * [1.2.3.可重复读执行报告](./1.%E4%BA%8B%E5%8A%A1/2.%E9%9A%94%E7%A6%BB%E6%80%A7/3.%E5%8F%AF%E9%87%8D%E5%A4%8D%E8%AF%BB)
    * [1.2.4.串行化执行报告](./1.%E4%BA%8B%E5%8A%A1/2.%E9%9A%94%E7%A6%BB%E6%80%A7/4.%E4%B8%B2%E8%A1%8C%E5%8C%96)
* [2.mysql 基本原理](./2.mysql%E5%9F%BA%E6%9C%AC%E5%8E%9F%E7%90%86)
  * [2.1.sql 语句执行流程](./2.mysql%E5%9F%BA%E6%9C%AC%E5%8E%9F%E7%90%86/1.sql%E8%AF%AD%E5%8F%A5%E6%89%A7%E8%A1%8C%E6%B5%81%E7%A8%8B)
  * [2.2.MySQL 主从架构](./2.mysql%E5%9F%BA%E6%9C%AC%E5%8E%9F%E7%90%86/2.MySQL%E4%B8%BB%E4%BB%8E%E6%9E%B6%E6%9E%84)
* [3.mysql 实践](./3.mysql%E5%AE%9E%E8%B7%B5)
  * [3.1.修改 redolog 大小](./3.mysql%E5%AE%9E%E8%B7%B5/%E4%BF%AE%E6%94%B9redolog%E5%A4%A7%E5%B0%8F)
  * [2.2.设置主从同步](./3.mysql%E5%AE%9E%E8%B7%B5/%E8%AE%BE%E7%BD%AE%E4%B8%BB%E4%BB%8E%E5%90%8C%E6%AD%A5)
