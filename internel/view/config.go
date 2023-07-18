package view

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

type Config struct {
	Mysql  MysqlClient
	Sqlite SqliteClient
	Mongo  MongoClient

	Listen string

	DBS map[string]*DBConnect
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

func (c *Config) Init() {
	c.DBS = map[string]*DBConnect{}
	switch {
	case c.Mysql.Address != "":
		key := "mysql"
		link := fmt.Sprintf("%s:%s@tcp(%s)/%s?",
			c.Mysql.Username, c.Mysql.Password, c.Mysql.Address, c.Mysql.DBName,
		) + "charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
		db, err := gorm.Open(mysql.Open(link), dbConf)
		if err != nil || db == nil {
			panic(err)
		}
		sqlScript := fmt.Sprintf(
			`SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'`,
			c.Mysql.DBName)

		c.DBS[key] = &DBConnect{
			DBName: c.Mysql.DBName,
			Script: sqlScript,
			Conn:   db,
		}
		fallthrough
	case c.Sqlite.Link != "":
		key := "sqlite"
		db, err := gorm.Open(sqlite.Open(c.Sqlite.Link), dbConf)
		if err != nil || db == nil {
			panic(err)
		}

		c.DBS[key] = &DBConnect{
			DBName: c.Sqlite.DBName,
			Script: `SELECT name FROM sqlite_master WHERE type='table'`,
			Conn:   db,
		}
		fallthrough
	case c.Mongo.Uri != "":
		key := "mongo"
		client, err := mongo.Connect(context.TODO(),
			options.Client().ApplyURI(c.Mongo.Uri).SetConnectTimeout(time.Minute))
		if err != nil {
			log.Fatal(err)
		}
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
		c.DBS[key] = &DBConnect{
			DBName: c.Mongo.DBName,
			Script: "",
			Conn:   client,
		}

	}

	if len(c.DBS) == 0 {
		log.Fatal("没有连接到任何一个数据库")
	}
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
