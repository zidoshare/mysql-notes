drop database if exists test_2;
create database test_2;
use test_2;
create table t(c int) engine=InnoDB;
insert into t(c) values(1);