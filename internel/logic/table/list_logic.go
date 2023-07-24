package table

import (
	"context"
	"errors"
	"github.com/7134g/viewAdmin/db"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/7134g/viewAdmin/internel/view"
	"github.com/Masterminds/squirrel"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"log"
)

type List struct {
	cfg *view.Config

	TableName string `json:"table_name"`
	DbType    string `json:"db_type,default=mysql"`

	db.MysqlQueryParams
}

func NewListLogic(c *view.Config) List {
	return List{cfg: c}
}

func (h *List) List(ctx *serve.BaseContext) (resp *db.ListResponse, err error) {
	if err := ctx.ShouldBindJSON(h); err != nil {
		return nil, err
	}

	var list []map[string]interface{}
	var count int64
	switch h.DbType {
	case db.MysqlType, db.SqliteType:
		list, count, err = h.getListByGorm(h.DbType)
		if err != nil {
			return nil, err
		}
	case db.MongoType:
		list, count, err = h.getListByMongo(h.DbType)
		if err != nil {
			return nil, err
		}
	}
	resp = &db.ListResponse{
		List:  list,
		Total: count,
	}

	return resp, nil
}

func (h *List) getListByGorm(dbType string) ([]map[string]interface{}, int64, error) {
	idb, ok := h.cfg.DBS[dbType]
	if !ok {
		return nil, 0, errors.New("cannot find " + dbType)
	}

	// count
	whereCondition, whereData := h.GetWhereSql()
	countScript, values, err := squirrel.Select("count(*)").From(h.TableName).Where(
		whereCondition, whereData...).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var count int64
	if err := idb.Conn.(*gorm.DB).Raw(countScript, values...).Scan(&count).Error; err != nil {
		return nil, 0, err
	}

	// list
	rowBuilder := squirrel.Select("*").From(h.TableName)
	orderBy := h.GetOrderBy()
	if orderBy != "" {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}
	offset := h.GetOffset()
	limit := h.GetLimit()
	sqlScript, values, err := rowBuilder.Where(whereCondition, whereData...).Offset(offset).Limit(limit).ToSql()
	if err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	if values != nil && len(values) > 0 {
		if err := idb.Conn.(*gorm.DB).Raw(sqlScript, values...).Scan(&result).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := idb.Conn.(*gorm.DB).Raw(sqlScript).Scan(&result).Error; err != nil {
			return nil, 0, err
		}
	}

	return result, count, nil
}

func (h *List) getListByMongo(dbType string) ([]map[string]interface{}, int64, error) {
	idb, ok := h.cfg.DBS[dbType]
	if !ok {
		return nil, 0, errors.New("cannot find " + dbType)
	}

	client := idb.Conn.(*mongo.Client)
	_db := client.Database(idb.DBName)
	collection := _db.Collection(h.TableName)

	filter := h.GetWhereBson()
	opts := options.Find().SetSkip(int64(h.GetOffset())).SetLimit(int64(h.GetLimit()))

	ctx := context.Background()

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	result := make([]map[string]interface{}, 0)
	for cur.Next(ctx) {
		var m map[string]interface{}
		if err = cur.Decode(&m); err != nil {
			log.Println(err)
		}
		result = append(result, m)

	}

	return result, count, err
}
