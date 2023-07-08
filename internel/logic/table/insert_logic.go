package table

import (
	"errors"
	"github.com/7134g/viewAdmin/db"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/7134g/viewAdmin/internel/view"
	"github.com/Masterminds/squirrel"
)

type Insert struct {
	cfg *view.Config

	insertData map[string]interface{}
}

func NewInsertLogic(c *view.Config) Insert {
	return Insert{cfg: c, insertData: map[string]interface{}{}}
}

func (h *Insert) Insert(ctx *serve.BaseContext) (interface{}, error) {
	tableName, exist := ctx.GetQuery("table_name")
	if !exist {
		return nil, errors.New("缺失表名")
	}

	if err := ctx.ShouldBind(&h.insertData); err != nil {
		return nil, err
	}
	insertData, err := db.FixJsonData(h.insertData)
	if err != nil {
		return nil, err
	}

	sqlScript, values, err := squirrel.Insert(tableName).SetMap(insertData).ToSql()
	if err != nil {
		return nil, err
	}

	rowBuilder := h.cfg.DB.Exec(sqlScript, values...)
	if err := rowBuilder.Error; err != nil {
		return nil, err
	}

	if rowBuilder.RowsAffected == 0 {
		return nil, errors.New("插入失败影响 0 列")
	}

	return "ok", nil
}
