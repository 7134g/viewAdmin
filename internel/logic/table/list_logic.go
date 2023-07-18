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
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"log"
)

type List struct {
	cfg *view.Config

	TableName string `json:"table_name"`
	DbType    string `json:"db_type"`

	db.MysqlQueryParams
}

func NewListLogic(c *view.Config) List {
	return List{cfg: c}
}

func (h *List) List(ctx *serve.BaseContext) (resp interface{}, err error) {
	if err := ctx.ShouldBind(h); err != nil {
		return nil, err
	}

	switch h.DbType {
	case db.MysqlType, db.SqliteType:
		resp, err = h.getListByGorm(h.DbType)
		if err != nil {
			return nil, err
		}
	case db.MongoType:
		resp, err = h.getListByMongo(h.DbType)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (h *List) getListByGorm(dbType string) ([]map[string]interface{}, error) {
	rowBuilder := squirrel.Select("*").From(h.TableName)
	orderBy := h.GetOrderBy()
	if orderBy != "" {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}
	whereCondition, whereData := h.GetWhereSql()
	offset := h.GetOffset()
	limit := h.GetLimit()
	sqlScript, values, err := rowBuilder.Where(whereCondition, whereData...).Offset(offset).Limit(limit).ToSql()
	if err != nil {
		return nil, err
	}

	idb, ok := h.cfg.DBS[dbType]
	if !ok {
		return nil, errors.New("cannot find " + dbType)
	}

	var resp []map[string]interface{}
	if values != nil && len(values) > 0 {
		if err := idb.Conn.(*gorm.DB).Raw(sqlScript, values...).Scan(&resp).Error; err != nil {
			return nil, err
		}
	} else {
		if err := idb.Conn.(*gorm.DB).Raw(sqlScript).Scan(&resp).Error; err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (h *List) getListByMongo(dbType string) ([]map[string]interface{}, error) {
	idb, ok := h.cfg.DBS[dbType]
	if !ok {
		return nil, errors.New("cannot find " + dbType)
	}

	client := idb.Conn.(*mongo.Client)
	_db := client.Database(idb.DBName)
	collection := _db.Collection(h.TableName)

	filter := bson.M{} // todo
	opts := options.Find().SetSkip(int64(h.GetOffset())).SetLimit(int64(h.GetLimit()))

	ctx := context.Background()
	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	resp := make([]map[string]interface{}, 0)
	for cur.Next(ctx) {
		var m map[string]interface{}
		if err = cur.Decode(&m); err != nil {
			log.Println(err)
		}
		resp = append(resp, m)

	}

	return resp, err
}
