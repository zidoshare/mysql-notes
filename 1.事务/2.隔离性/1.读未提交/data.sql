drop database if exists test_1_1;
create database test_1;
use test_1;
create table t(c int) engine=InnoDB;
insert into t(c) values(1);