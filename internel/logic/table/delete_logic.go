package table

import (
	"errors"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/7134g/viewAdmin/internel/view"
	"github.com/Masterminds/squirrel"
)

type Delete struct {
	cfg *view.Config
}

func NewDeleteLogic(c *view.Config) Delete {
	return Delete{cfg: c}
}

func (h *Delete) Delete(ctx *serve.BaseContext) (interface{}, error) {
	tableName, exist := ctx.GetQuery("table_name")
	if !exist {
		return nil, errors.New("缺失表名")
	}
	id := ctx.GetQueryInt("id")
	if id == 0 {
		return nil, errors.New("缺失id")
	}
	if err := ctx.ShouldBind(h); err != nil {
		return nil, err
	}

	sqlScript, values, err := squirrel.Delete(tableName).Where("id = ?", id).ToSql()
	if err != nil {
		return nil, err
	}

	rowBuilder := h.cfg.DB.Exec(sqlScript, values...)
	if err := rowBuilder.Error; err != nil {
		return nil, err
	}

	if rowBuilder.RowsAffected == 0 {
		return nil, errors.New("删除失败影响 0 列")
	}

	return "ok", nil
}