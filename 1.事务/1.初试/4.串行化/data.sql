drop database if exists test;
create database test;
use test;
create table t(c int) engine=InnoDB;
insert into t(c) values(1);