# Docker 部署数据库写入测试数据

## mysql

dockerhub: https://hub.docker.com/_/mysql
```shell
$ docker run --name showdocdb-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql:8.0.29
```

测试数据
```sql
CREATE DATABASE demo;

CREATE TABLE IF NOT EXISTS `user`(
   `id` INT UNSIGNED AUTO_INCREMENT COMMENT '主键',
   `user_name` VARCHAR(20) NOT NULL COMMENT '用户名',
   `group_id` TINYINT(2) NOT NULL DEFAULT 2 COMMENT '1为管理员，2为普通用户',
   `password` VARCHAR(50) NOT NULL COMMENT '密码',
   `email` VARCHAR(50) NOT NULL COMMENT '邮箱',
   `avatar` VARCHAR(200) NULL COMMENT '头像',
   `avatar_small` VARCHAR(200) NULL COMMENT '小头像',
   `name` VARCHAR(5) NULL COMMENT '昵称',
   `reg_time` INT(11) NOT NULL DEFAULT 0 COMMENT '注册时间',
   `last_login_time` INT(11) NOT NULL DEFAULT 0 COMMENT '最后一次登录时间',
   PRIMARY KEY ( `id` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户表，储存用户信息';

CREATE TABLE IF NOT EXISTS `page`(
    `page_id` INT UNSIGNED AUTO_INCREMENT COMMENT '页面自增id',
    `author_id` INT(10) NOT NULL DEFAULT 0 COMMENT '页面作者id',
    `author_username` VARCHAR(50) NOT NULL COMMENT '页面作者用户名',
    `item_id` INT(10) NOT NULL DEFAULT 0 COMMENT '项目id',
    `cat_id` INT(10) NOT NULL DEFAULT 0 COMMENT '父目录id',
    `page_title` VARCHAR(50) NOT NULL COMMENT '页面标题',
    `page_content` TEXT NOT NULL COMMENT '页面内容',
    `order` INT(10) NOT NULL DEFAULT 99 COMMENT '顺序号，数字越小越靠前',
    `addtime` INT(11) NOT NULL DEFAULT 0 COMMENT '该记录添加的时间',
    PRIMARY KEY ( `page_id` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '页面表，保存编辑的页面内容';


CREATE TABLE IF NOT EXISTS `item`(
    `item_id` INT UNSIGNED AUTO_INCREMENT COMMENT '项目id、自增id',
    `item_name` VARCHAR(50) NOT NULL COMMENT '项目名',
    `item_description` VARCHAR(225) NOT NULL COMMENT '项目描述',
    `user_id` INT(10) NULL DEFAULT 0 COMMENT '创建人id',
    `user_name` VARCHAR(50) NOT NULL COMMENT '创建人用户名',
    `password` VARCHAR(50) NOT NULL COMMENT '项目密码。可为空',
    `addtime` INT(11) NOT NULL DEFAULT 0 COMMENT '项目添加的时间，时间戳',
    PRIMARY KEY ( `item_id` )
    )ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '项目表，储存项目信息';
```