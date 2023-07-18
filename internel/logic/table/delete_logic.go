package table

import (
	"context"
	"errors"
	"github.com/7134g/viewAdmin/db"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/7134g/viewAdmin/internel/view"
	"github.com/Masterminds/squirrel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Delete struct {
	cfg *view.Config

	ID        int    `query:"id"`
	TableName string `query:"table_name"`
	DbType    string `query:"db_type"`
}

func NewDeleteLogic(c *view.Config) Delete {
	return Delete{cfg: c}
}

func (h *Delete) Delete(ctx *serve.BaseContext) (interface{}, error) {
	if err := ctx.ShouldBindQuery(h); err != nil {
		return nil, err
	}

	switch h.DbType {
	case db.MysqlType, db.SqliteType:
		err := h.deleteByGorm(h.ID)
		if err != nil {
			return nil, err
		}
	case db.MongoType:
		err := h.deleteByMongo(h.ID)
		if err != nil {
			return nil, err
		}
	}

	return "ok", nil
}

func (h *Delete) deleteByGorm(id int) error {
	sqlScript, values, err := squirrel.Delete(h.TableName).Where("id = ?", id).ToSql()
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
		return errors.New("删除失败影响 0 列")
	}
	return nil
}

func (h *Delete) deleteByMongo(id int) error {
	idb, ok := h.cfg.DBS[h.DbType]
	if !ok {
		return errors.New("cannot find " + h.DbType)
	}

	client := idb.Conn.(*mongo.Client)
	_db := client.Database(idb.DBName)
	collection := _db.Collection(h.TableName)

	del, err := collection.DeleteOne(context.Background(), bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	if del.DeletedCount == 0 {
		return errors.New("删除失败影响 0 列")
	}

	return nil
}
