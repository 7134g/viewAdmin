package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

type DBConnect struct {
	DBName string // 数据库名称
	Script string // 获取表信息的sql语句
	Conn   interface{}
}

type MysqlClient struct {
	DBName   string
	Username string
	Password string
	Address  string

	Script string
}

func (c *MysqlClient) getConn() *DBConnect {
	link := fmt.Sprintf("%s:%s@tcp(%s)/%s?",
		c.Username, c.Password, c.Address, c.DBName,
	) + "charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	db, err := gorm.Open(mysql.Open(link), dbConf)
	if err != nil || db == nil {
		panic(err)
	}
	sqlScript := fmt.Sprintf(
		c.Script,
		c.DBName)

	return &DBConnect{
		DBName: c.DBName,
		Script: sqlScript,
		Conn:   db,
	}
}

type SqliteClient struct {
	Link   string
	DBName string

	Script string
}

func (c *SqliteClient) getConn() *DBConnect {
	db, err := gorm.Open(sqlite.Open(c.Link), dbConf)
	if err != nil || db == nil {
		panic(err)
	}

	return &DBConnect{
		DBName: c.DBName,
		Script: c.Script,
		Conn:   db,
	}
}

type MongoClient struct {
	Uri    string
	DBName string
}

func (c *MongoClient) getConn() *DBConnect {
	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI(c.Uri).SetConnectTimeout(time.Minute))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return &DBConnect{
		DBName: c.DBName,
		Script: "",
		Conn:   client,
	}
}

type RedisClient struct {
	Host     string
	Password string
	DBName   string

	Script string
}

func (c *RedisClient) getConn() *DBConnect {
	if c.DBName == "" {
		c.DBName = "0"
	}
	rds := redis.NewClient(&redis.Options{
		Addr:     c.Host,     // Redis服务器地址和端口号
		Password: c.Password, // 密码（如果有）
		DB:       0,          // 使用默认数据库
	})

	return &DBConnect{
		DBName: c.DBName,
		Script: c.Script,
		Conn:   rds,
	}
}
