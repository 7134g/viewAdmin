package table

import (
	"context"
	"errors"
	db2 "github.com/7134g/viewAdmin/common/db"
	"github.com/7134g/viewAdmin/config"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/Masterminds/squirrel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Update struct {
	cfg *config.Config

	TableName  string                 `json:"table_name"`
	DbType     string                 `json:"db_type,default=mysql"`
	ID         interface{}            `json:"id"`
	UpdateData map[string]interface{} `json:"update_data"`
}

func NewUpdateLogic(c *config.Config) Update {
	return Update{cfg: c}
}

func (h *Update) Update(ctx *serve.BaseContext) (resp interface{}, err error) {
	if err := ctx.ShouldBindJSON(h); err != nil {
		return nil, err
	}
	switch h.DbType {
	case db2.MysqlType, db2.SqliteType:
		err = h.updateByGorm()
		if err != nil {
			return nil, err
		}
	case db2.MongoType:
		err = h.updateByMongo()
		if err != nil {
			return nil, err
		}
	}

	return "ok", nil
}

func (h *Update) updateByGorm() error {

	updateData, err := db2.FixJsonData(h.UpdateData)
	if err != nil {
		return err
	}

	sqlScript, values, err := squirrel.Update(h.TableName).Where("id = ?", h.ID).SetMap(updateData).ToSql()
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
		return errors.New("更新失败影响 0 列")
	}

	return nil
}

func (h *Update) updateByMongo() error {
	idb, ok := h.cfg.DBS[h.DbType]
	if !ok {
		return errors.New("cannot find " + h.DbType)
	}

	client := idb.Conn.(*mongo.Client)
	_db := client.Database(idb.DBName)
	collection := _db.Collection(h.TableName)

	_id, err := primitive.ObjectIDFromHex(h.ID.(string))
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", _id}}
	update := bson.D{
		{"$set", h.UpdateData},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("更新失败影响 0 列")
	}

	return nil

}
