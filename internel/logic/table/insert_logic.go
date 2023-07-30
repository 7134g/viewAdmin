package table

import (
	"context"
	"errors"
	db2 "github.com/7134g/viewAdmin/common/db"
	"github.com/7134g/viewAdmin/config"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/Masterminds/squirrel"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Insert struct {
	cfg *config.Config

	TableName string `json:"table_name"`
	DbType    string `json:"db_type,default=mysql"`

	InsertData map[string]interface{} `json:"insert_data"`
}

func NewInsertLogic(c *config.Config) Insert {
	return Insert{cfg: c, InsertData: map[string]interface{}{}}
}

func (h *Insert) Insert(ctx *serve.BaseContext) (interface{}, error) {
	if err := ctx.ShouldBindJSON(h); err != nil {
		return nil, err
	}
	insertData, err := db2.FixJsonData(h.InsertData)
	if err != nil {
		return nil, err
	}

	switch h.DbType {
	case db2.MysqlType, db2.SqliteType:
		err := h.insertByGorm(insertData)
		if err != nil {
			return nil, err
		}
	case db2.MongoType:
		err := h.insertByMongo(insertData)
		if err != nil {
			return nil, err
		}
	}

	return "ok", nil
}

func (h *Insert) insertByGorm(insertData map[string]interface{}) error {
	sqlScript, values, err := squirrel.Insert(h.TableName).SetMap(insertData).ToSql()
	if err != nil {
		return err
	}

	idb, ok := h.cfg.DBS[h.DbType]
	if !ok {
		return errors.New("cannot find " + h.DbType)
	}

	rowBuilder := idb.Conn.(*gorm.DB).Exec(sqlScript, values...)
	if err := rowBuilder.Error; err != nil {
		return err
	}

	if rowBuilder.RowsAffected == 0 {
		return errors.New("插入失败影响 0 列")
	}

	return nil
}

func (h *Insert) insertByMongo(insertData map[string]interface{}) error {
	idb, ok := h.cfg.DBS[h.DbType]
	if !ok {
		return errors.New("cannot find " + h.DbType)
	}

	client := idb.Conn.(*mongo.Client)
	_db := client.Database(idb.DBName)
	collection := _db.Collection(h.TableName)

	_, err := collection.InsertOne(context.Background(), insertData)
	if err != nil {
		return err
	}

	return nil
}
