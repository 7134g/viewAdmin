package config

import (
	"github.com/7134g/viewAdmin/common/db"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

type Config struct {
	Mysql  MysqlClient
	Sqlite SqliteClient
	Mongo  MongoClient
	Redis  RedisClient

	Listen string
	Mode   string
	Front  bool

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
		c.DBS[db.MysqlType] = c.Mysql.getConn()
		fallthrough
	case c.Sqlite.Link != "":
		c.DBS[db.SqliteType] = c.Sqlite.getConn()
		fallthrough
	case c.Mongo.Uri != "":
		c.DBS[db.MongoType] = c.Mongo.getConn()
		fallthrough
	case c.Redis.Host != "":
		c.DBS[db.RedisType] = c.Redis.getConn()
		fallthrough
	default:
		if len(c.DBS) == 0 {
			log.Fatal("没有连接到任何一个数据库")
		} else {
			for k, _ := range c.DBS {
				log.Printf("成功连接：%s \n", k)
			}
		}
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
