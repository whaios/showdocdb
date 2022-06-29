# showdocdb

自动化生成 [ShowDoc](https://github.com/star7th/showdoc)  数据字典文档，支持 mysql、postgres、sqlserver、sqlite3、oracle。

- [安装](#安装)
    - [下载可执行文件](#下载可执行文件)
    - [源码编译](#源码编译)
- [设置变量](#设置变量)
- [生成数据字典](#生成数据字典)
    - [MySQL](#MySQL)
    - [PostgreSQL](#PostgreSQL)
    - [SQLServer](#SQLServer)
    - [SQlite](#SQlite)
    - [Oracle](#Oracle)

## 安装

### 下载可执行文件

已经编译好的平台有： [点击下载](https://github.com/whaios/showdocdb/releases)
- windows/amd64
- linux/amd64
- darwin/amd64

如果需要连接 SQlite 或 Oracle 数据库生成数据字典，请下载对应的cgo版本可执行程序。

## 通用参数

### 1. apihost 

ShowDoc 服务器地址

工具中默认配置的官方线上地址 `https://www.showdoc.com.cn` ，如果你也使用的该服务则不需要修改。

如果你使用的是私有版ShowDoc，则使用时需要通过 `--apihost` 参数指定为自己的地址。

为了避免每次都手动输入该参数，建议将该地址配置为环境变量 `GOSHOWDOC_HOST`。

### 2. apikey 和 apitoken

开放 API 认证凭证

工具生成的文档最终会通过开放API同步到 ShowDoc 的项目中，所以需要配置认证凭证（api_key 和 api_token） 。

登录showdoc > 进入具体项目 > 点击右上角的”项目设置” > “开放API” 便可看到。

为了避免每次生成时都输入这两个参数，建议将该参数配置为环境变量 `GOSHOWDOC_APIKEY` 和 `GOSHOWDOC_APITOKEN`。

### 3. debug

设置 `--debug` 参数开启调试模式，输出更详细的日志。

## 生成数据字典

| 参数       | 简写      | 说明                                                                  |
|:---------|:--------|:--------------------------------------------------------------------|
| --cat    | -cat    | 文档所在目录，如果需要多层目录请用斜杠隔开，例如：“一层/二层/三层”                                 |
| --driver | -d      | 数据库类型，支持：mysql、postgres、sqlserver、sqlite3、oracle (default: "mysql") |
| --host   | -h      | 数据库地址和端口，如果是SQlite数据库则为文件 (default: "127.0.0.1:3306")               |
| --user   | -u      | 数据库用户名                                                              |
| --pwd    | -p      | 数据库密码                                                               |
| --db     | -db     | 要同步的数据库名                                                            |
| --schema | -schema | PostgreSQL 数据库模式 (default: "public")                                |

### MySQL

```shell
.\showdocdb-windows-amd64.exe -cat 数据字典演示/MySQL -d mysql -h 127.0.0.1:3306 -u root -p 123456 --db demo
```

### PostgreSQL

```shell
.\showdocdb-windows-amd64.exe -d postgres -h 127.0.0.1:5432 -u postgres -p 123456 -db postgres
```

### SQLServer

```shell
.\showdocdb-windows-amd64.exe -d sqlserver -h 127.0.0.1:1433 -u sa -p 123456 -db testdb
```

### SQlite

因为 go-sqlite3 库是一个 cgo 库，编译代码时需要 gcc 环境。

```shell
.\showdocdb-windows-amd64-cgo.exe -driver sqlite3 -h .\test.db
```

### Oracle

因为 godror 库是一个 cgo 库，编译代码时需要 gcc 环境。

**注意：**
连接 Oracle 需要安装 Oracle 客户端库，可查看 [ODPI-C](https://oracle.github.io/odpi/doc/installation.html) 文档，
从 https://www.oracle.com/database/technologies/instant-client/downloads.html 下载免费的Basic或Basic Light软件包。

```shell
.\showdocdb-windows-amd64-cgo.exe -d oracle -h 127.0.0.1:1521 -u scott -p tiger -db orclpdb1
```