package table

import (
	"context"
	"errors"
	"fmt"
	"github.com/7134g/viewAdmin/db"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/7134g/viewAdmin/internel/view"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"log"
	"strings"
)

type ViewTable struct {
	cfg *view.Config
}

func NewViewTableLogic(c *view.Config) ViewTable {
	return ViewTable{cfg: c}
}

func (h *ViewTable) ViewTable(ctx *serve.BaseContext) (interface{}, error) {
	dbType, ok := ctx.GetQuery("db_type")
	if !ok {
		return nil, errors.New("need db_type")
	}

	table := make(map[string]map[string]interface{})
	switch dbType {
	case db.SqliteType, db.MysqlType:
		tbn := h.getTableNameByGorm(dbType)
		if tbn == nil {
			return nil, nil
		}
		table = h.getTableStructByGorm(dbType, tbn)
	case db.MongoType:
		tbn := h.getTableNameByMongo(dbType)
		if tbn == nil {
			return nil, nil
		}
	}

	return table, nil
}

func (h *ViewTable) getTableNameByGorm(dbType string) []string {
	idb := h.cfg.DBS[dbType]
	_db, ok := idb.Conn.(*gorm.DB)
	if !ok {
		return nil
	}

	tablesName := make([]string, 0)
	err := _db.Raw(idb.Script).Scan(&tablesName).Error
	if err != nil {
		return nil
	}

	return tablesName
}

func (h *ViewTable) getTableStructByGorm(dbType string, tbn []string) map[string]map[string]interface{} {
	table := make(map[string]map[string]interface{})
	for _, name := range tbn {
		var m map[string]interface{}
		sqlScript := fmt.Sprintf(`SHOW CREATE TABLE %s`, name)
		idb := h.cfg.DBS[dbType]
		_db, ok := idb.Conn.(*gorm.DB)
		if !ok {
			return nil
		}

		err := _db.Raw(sqlScript).Scan(&m).Error
		if err != nil {
			return nil
		}
		m = h.parseTable(m)
		table[name] = m
	}
	return table
}

func (h *ViewTable) parseTable(m map[string]interface{}) map[string]interface{} {
	createSql, ok := m["Create Table"]
	if !ok {
		return nil
	}
	v, ok := createSql.(string)
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	lines := strings.Split(v, "\n")
	for i := 1; i < len(lines)-1; i++ {
		line := lines[i]

		if h.skip(line) {
			continue
		}

		line = strings.TrimLeft(line, " ")
		line = strings.TrimRight(line, " ")
		fields := strings.Split(line, " ")
		key := strings.ReplaceAll(fields[0], "`", "")
		value := fields[1]
		result[key] = value
	}
	return result
}

func (h *ViewTable) skip(line string) bool {
	switch {
	case strings.Contains(line, "UNIQUE"),
		strings.Contains(line, "PRIMARY"),
		strings.Contains(line, "CONSTRAINT"),
		strings.Contains(line, "KEY"),
		strings.Contains(line, "UNIQUE"),
		strings.Contains(line, "CONSTRAINT"):
		return true
	default:
		return false
	}
}

func (h *ViewTable) getTableNameByMongo(dbType string) []string {
	idb := h.cfg.DBS[dbType]
	_db, ok := idb.Conn.(*mongo.Client)
	if !ok {
		return nil
	}

	collections, err := _db.Database(idb.DBName).ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	return collections
}

func (h *ViewTable) getTableStructByMongo(dbType string, tbn []string) map[string]map[string]interface{} {
	idb := h.cfg.DBS[dbType]
	client, ok := idb.Conn.(*mongo.Client)
	if !ok {
		return nil
	}

	_db := client.Database(idb.DBName)
	ctx := context.TODO()
	table := make(map[string]map[string]interface{})
	for _, name := range tbn {
		collection := _db.Collection(name)
		cur, err := collection.Find(ctx, bson.D{}, options.Find())
		if err != nil {
			log.Fatal(err)
		}

		result := map[string]interface{}{}
		for cur.Next(ctx) {
			var m map[string]interface{}
			if err = cur.Decode(&m); err != nil {
				log.Println(err)
			}

			for k, v := range m {
				saveValue, ok := result[k]
				if !ok {
					result[k] = fmt.Sprintf("%T", v)
				} else {
					result[k] = fmt.Sprintf("%T,%T", saveValue, v)
				}

			}

		}
		table[name] = result

	}

	return table
}
