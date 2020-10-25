drop database if exists test_3_3;
create database test_3_3;
use test_3_3;
create table t(c int) engine=InnoDB;
insert into t(c) values(1);