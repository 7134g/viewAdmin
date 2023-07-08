package view

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

type Config struct {
	Mysql  MysqlClient
	Sqlite SqliteClient

	Listen         string
	DB             *gorm.DB
	SqlTableScript string
}

func (c *Config) Init() {
	switch {
	case c.Mysql.Address != "":
		link := fmt.Sprintf("%s:%s@tcp(%s)/%s?",
			c.Mysql.Username, c.Mysql.Password, c.Mysql.Address, c.Mysql.DBName,
		) + "charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
		db, err := gorm.Open(mysql.Open(link), dbConf)
		if err != nil || db == nil {
			panic(err)
		}
		c.DB = db
		sqlScript := fmt.Sprintf(
			`SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'`,
			c.Mysql.DBName)

		c.SqlTableScript = sqlScript
	case c.Sqlite.Link != "":
		db, err := gorm.Open(sqlite.Open(c.Sqlite.Link), dbConf)
		if err != nil || db == nil {
			panic(err)
		}
		c.DB = db
		c.SqlTableScript = `SELECT name FROM sqlite_master WHERE type='table'`
	}

	if c.DB == nil {
		log.Fatal("没有连接到任何一个数据库")
	}
}

type MysqlClient struct {
	DBName   string
	Username string
	Password string
	Address  string
}

type SqliteClient struct {
	Link string
}

var dbConf = &gorm.Config{
	PrepareStmt: true,
	Logger:      logger.Default.LogMode(logger.Silent),
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
}

func OpenLogDB() {
	dbConf.Logger = logger.Default.LogMode(logger.Info)
}

func InitConfig(path string) *Config {
	vp := viper.New()
	vp.SetConfigFile(path)
	if err := vp.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	c := &Config{}
	if err := vp.Unmarshal(c); err != nil {
		log.Fatal(err)
	}

	c.Init()

	return c
}
