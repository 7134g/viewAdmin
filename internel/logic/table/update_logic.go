package table

import (
	"errors"
	"github.com/7134g/viewAdmin/db"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/7134g/viewAdmin/internel/view"
	"github.com/Masterminds/squirrel"
)

type Update struct {
	cfg *view.Config

	updateData map[string]interface{}
}

func NewUpdateLogic(c *view.Config) Update {
	return Update{cfg: c}
}

func (h *Update) Update(ctx *serve.BaseContext) (interface{}, error) {
	tableName, exist := ctx.GetQuery("table_name")
	if !exist {
		return nil, errors.New("缺失表名")
	}
	id := ctx.GetQueryInt("id")
	if id == 0 {
		return nil, errors.New("缺失id")
	}

	if err := ctx.ShouldBind(&h.updateData); err != nil {
		return nil, err
	}

	updateData, err := db.FixJsonData(h.updateData)
	if err != nil {
		return nil, err
	}

	sqlScript, values, err := squirrel.Update(tableName).Where("id = ?", id).SetMap(updateData).ToSql()
	if err != nil {
		return nil, err
	}

	rowBuilder := h.cfg.DB.Exec(sqlScript, values...)
	if err := rowBuilder.Error; err != nil {
		return nil, err
	}

	if rowBuilder.RowsAffected == 0 {
		return nil, errors.New("更新失败影响 0 列")
	}

	return "ok", nil
}
