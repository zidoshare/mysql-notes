drop database if exists test_2_3;
create database test_2_3;
use test_2_3;
create table t(c int) engine=InnoDB;
insert into t(c) values(1);