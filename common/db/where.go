package db

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
)

type OrderType uint

const (
	Desc = "desc"
	Asc  = "asc"
)

type GormQueryParams struct {
	//squirrel.SelectBuilder

	Page      int                    `json:"page"`
	Size      int                    `json:"size"`
	OrderKey  string                 `json:"order_key"`
	OrderType OrderType              `json:"order_type"`
	WhereSql  map[string]interface{} // key: name = ?  \  value: data

}

type ListResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

func (m *GormQueryParams) GetOrderBy() string {
	var orderType string
	switch m.OrderType {
	case 1:
		orderType = Desc
	case 2:
		orderType = Asc
	default:
		return ""
	}

	return fmt.Sprintf("%s %s", m.OrderKey, orderType)
}

func (m *GormQueryParams) GetWhereSql() (string, []interface{}) {
	link := " AND "
	sqlConditions := ""
	sqlData := make([]interface{}, 0)
	for k, v := range m.WhereSql {
		if sqlConditions == "" {
			sqlConditions = k
			sqlData = append(sqlData, v)
			continue
		}

		sqlConditions = sqlConditions + link + k
		sqlData = append(sqlData, v)
	}

	return sqlConditions, sqlData
}

func (m *GormQueryParams) GetWhereBson() bson.D {
	filter := bson.D{}
	for k, v := range m.WhereSql {
		filter = append(filter, bson.E{k, v})
	}
	return filter
}

func (m *GormQueryParams) GetOffset() uint64 {
	if m.Page < 1 {
		m.Page = 1
	}
	if m.Size < 1 {
		m.Size = 1
	}

	return uint64((m.Page - 1) * m.Size)
}

func (m *GormQueryParams) GetLimit() uint64 {
	return uint64(m.Size)
}

func FixJsonData(data map[string]interface{}) (map[string]interface{}, error) {
	newData := map[string]interface{}{}
	for k, v := range data {

		v, err := judgeType(v)
		if err != nil {
			return nil, err
		}
		newData[k] = v
	}
	return newData, nil
}

func judgeType(v interface{}) (interface{}, error) {
	switch v.(type) {
	case []interface{}:
		values := v.([]interface{})
		if len(values) == 0 {
			//(*data)[k] = v
			return v, nil
		}

		result := make([]interface{}, 0)
		for _, value := range values {
			v, err := judgeType(value)
			if err != nil {
				return nil, err
			}
			result = append(result, v)
		}
		return result, nil

	case float64:
		valueString := strconv.FormatFloat(v.(float64), 'f', 64, 64)
		valueString = strings.TrimRight(valueString, "0")
		if strings.HasSuffix(valueString, ".") {
			value, err := strconv.Atoi(valueString[:len(valueString)-1])
			if err != nil {
				return nil, err
			}
			return value, nil
		} else {
			return v, nil
		}
	case string, bool:
		return v, nil
	default:
		return nil, errors.New("数据类型不正确")

	}

}
