package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/whaios/showdocdb/database"
	"github.com/whaios/showdocdb/log"
	"github.com/whaios/showdocdb/showdoc"
	"os"
	"strconv"
)

const Version = "1.1.0"

const (
	SHOWDOC_HOST     = "SHOWDOC_HOST"     // 环境变量：ShowDoc 地址。
	SHOWDOC_APIKEY   = "SHOWDOC_APIKEY"   // 环境变量：ShowDoc 开放 API 认证凭证。
	SHOWDOC_APITOKEN = "SHOWDOC_APITOKEN" // 环境变量：ShowDoc 开放 API 认证凭证。
)

const (
	flagDriver = "driver"
	flagHost   = "host"
	flagUser   = "user"
	flagPwd    = "pwd"
	flagDb     = "db"
	flagSchema = "schema"
	flagCat    = "cat"
)

func main() {
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "显示帮助",
	}
	app := &cli.App{
		Name:    "showdocdb",
		Usage:   "自动化生成 ShowDoc 数据字典文档，支持 mysql、postgres、sqlserver、sqlite3、oracle。",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "apihost",
				Usage:       "ShowDoc 地址。",
				Value:       showdoc.Host,
				Destination: &showdoc.Host,
				EnvVars:     []string{SHOWDOC_HOST},
			},
			&cli.StringFlag{
				Name:        "apikey",
				Usage:       "ShowDoc 开放 API 认证凭证。",
				Value:       showdoc.ApiKey,
				Destination: &showdoc.ApiKey,
				EnvVars:     []string{SHOWDOC_APIKEY},
			},
			&cli.StringFlag{
				Name:        "apitoken",
				Usage:       "ShowDoc 开放 API 认证凭证。",
				Value:       showdoc.ApiToken,
				Destination: &showdoc.ApiToken,
				EnvVars:     []string{SHOWDOC_APITOKEN},
			},

			&cli.StringFlag{
				Name:  flagCat,
				Usage: "文档所在目录，如果需要多层目录请用斜杠隔开，例如：“一层/二层/三层”",
			},
			&cli.StringFlag{
				Name:    flagDriver,
				Aliases: []string{"d"},
				Usage: fmt.Sprintf("数据库类型，支持：%s、%s、%s、%s、%s",
					database.MySQL, database.PostgreSQL, database.SQLServer, database.SQlite, database.Oracle),
				Value: database.MySQL,
			},
			&cli.StringFlag{
				Name:    flagHost,
				Aliases: []string{"h"},
				Usage:   "数据库地址和端口，如果是SQlite数据库则为文件",
				Value:   "127.0.0.1:3306",
			},
			&cli.StringFlag{
				Name:    flagUser,
				Aliases: []string{"u"},
				Usage:   "数据库用户名",
			},
			&cli.StringFlag{
				Name:    flagPwd,
				Aliases: []string{"p"},
				Usage:   "数据库密码",
			},
			&cli.StringFlag{
				Name:  flagDb,
				Usage: "要同步的数据库名",
			},
			&cli.StringFlag{
				Name:  flagSchema,
				Usage: "PostgreSQL 数据库模式",
				Value: "public",
			},

			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "开启调试模式。",
				Value:       log.IsDebug,
				Destination: &log.IsDebug,
			},
		},
		Action: func(c *cli.Context) error {
			UpdateDataDict(
				c.String(flagDriver),
				c.String(flagHost),
				c.String(flagUser),
				c.String(flagPwd),
				c.String(flagDb),
				c.String(flagSchema),
				c.String(flagCat),
			)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// UpdateDataDict 生成数据字典
func UpdateDataDict(driver, host, user, pwd, db, schema, cat string) {
	query, err := database.NewQuery(driver, host, user, pwd, db, schema)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if err = query.Open(); err != nil {
		log.Error("无法连接到数据库: %s", err.Error())
		return
	}
	defer query.Close()

	tbs, err := query.Query()
	if err != nil {
		log.Error("查询数据库出错: %s", err.Error())
		return
	}
	max := len(tbs)
	for i, tb := range tbs {
		if err = showdoc.UpdateByApi(cat, tb.Name, strconv.FormatInt(int64(i+1), 10), tb.Markdown()); err != nil {
			log.Error("更新文档[%s/%s]失败: %s", cat, tb.Name, err.Error())
			return
		}
		log.DrawProgressBar("更新文档", i+1, max)
	}
	log.Success("更新完成")
}
