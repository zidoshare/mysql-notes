drop database if exists test_4;
create database test_4;
use test_4;
create table t(c int) engine=InnoDB;
insert into t(c) values(1);