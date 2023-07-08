package table

import (
	"errors"
	"github.com/7134g/viewAdmin/db"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/7134g/viewAdmin/internel/view"
	"github.com/Masterminds/squirrel"
)

type List struct {
	cfg *view.Config

	db.MysqlQueryParams
}

func NewListLogic(c *view.Config) List {
	return List{cfg: c}
}

func (h *List) List(ctx *serve.BaseContext) (interface{}, error) {
	tableName, exist := ctx.GetQuery("table_name")
	if !exist {
		return nil, errors.New("缺失表名")
	}

	if err := ctx.ShouldBind(h); err != nil {
		return nil, err
	}

	rowBuilder := squirrel.Select("*").From(tableName)
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

	var resp []map[string]interface{}
	if values != nil && len(values) > 0 {
		if err := h.cfg.DB.Raw(sqlScript, values...).Scan(&resp).Error; err != nil {
			return nil, err
		}
	} else {
		if err := h.cfg.DB.Raw(sqlScript).Scan(&resp).Error; err != nil {
			return nil, err
		}
	}

	return resp, nil
}
