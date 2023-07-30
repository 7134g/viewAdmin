package table

import (
	"context"
	"fmt"
	"github.com/7134g/viewAdmin/common/db"
	"github.com/7134g/viewAdmin/config"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"log"
	"strings"
)

type ViewTable struct {
	cfg    *config.Config
	DbType string `form:"db_type,default=mysql"`
}

func NewViewTableLogic(c *config.Config) ViewTable {
	return ViewTable{cfg: c}
}

func (h *ViewTable) ViewTable(ctx *serve.BaseContext) (interface{}, error) {
	if err := ctx.ShouldBindQuery(h); err != nil {
		return nil, err
	}

	table := make(map[string]map[string]interface{})
	switch h.DbType {
	case db.SqliteType, db.MysqlType:
		tbn, err := h.getTableNameByGorm()
		if err != nil {
			return nil, err
		}
		table = h.getTableStructByGorm(tbn)
	case db.MongoType:
		tbn, err := h.getTableNameByMongo()
		if err != nil {
			return nil, err
		}
		table = h.getTableStructByMongo(tbn)
	case db.RedisType:
		tbn, err := h.getRedisKeys()
		if err != nil {
			return nil, err
		}
		table = h.getTableStructByRedis(tbn)
	}

	return table, nil
}

// gorm
func (h *ViewTable) getTableNameByGorm() ([]string, error) {
	idb := h.cfg.DBS[h.DbType]
	_db, ok := idb.Conn.(*gorm.DB)
	if !ok {
		return nil, db.DbNotConnect
	}

	tablesName := make([]string, 0)
	err := _db.Raw(idb.Script).Scan(&tablesName).Error
	if err != nil {
		return nil, err
	}

	return tablesName, nil
}

func (h *ViewTable) getTableStructByGorm(tbn []string) map[string]map[string]interface{} {
	table := make(map[string]map[string]interface{})
	for _, name := range tbn {
		var m map[string]interface{}
		sqlScript := fmt.Sprintf(`SHOW CREATE TABLE %s`, name)
		idb := h.cfg.DBS[h.DbType]
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

// mongo
func (h *ViewTable) getTableNameByMongo() ([]string, error) {
	idb := h.cfg.DBS[h.DbType]
	_db, ok := idb.Conn.(*mongo.Client)
	if !ok {
		return nil, db.DbNotConnect
	}

	collections, err := _db.Database(idb.DBName).ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (h *ViewTable) getTableStructByMongo(tbn []string) map[string]map[string]interface{} {
	idb := h.cfg.DBS[h.DbType]
	client, ok := idb.Conn.(*mongo.Client)
	if !ok {
		return nil
	}

	_db := client.Database(idb.DBName)
	ctx := context.TODO()
	table := make(map[string]map[string]interface{})
	for _, name := range tbn {
		collection := _db.Collection(name)
		cur, err := collection.Find(ctx, bson.D{}, options.Find().SetLimit(10))
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
					vType := fmt.Sprintf("%T", v)
					if vType == saveValue {
						continue
					}
					result[k] = fmt.Sprintf("%s,%s", saveValue, vType)
				}

			}
		}
		table[name] = result

	}

	return table
}

// redis
func (h *ViewTable) getRedisKeys() ([]string, error) {
	idb := h.cfg.DBS[h.DbType]
	_db, ok := idb.Conn.(*redis.Client)
	if !ok {
		return nil, db.DbNotConnect
	}

	keysCmd := _db.Eval(context.Background(), idb.Script, nil)
	if err := keysCmd.Err(); err != nil {
		log.Println("执行 redis EVAL命令失败：", err)
		return nil, err
	}

	keys, err := keysCmd.Result()
	if err != nil {
		log.Println("执行 redis EVAL命令失败：", err)
		return nil, err
	}

	result := make([]string, 0)
	for _, key := range keys.([]interface{}) {
		result = append(result, key.(string))
	}
	return result, err
}

func (h *ViewTable) getTableStructByRedis(tbn []string) map[string]map[string]interface{} {
	idb := h.cfg.DBS[h.DbType]
	_db, _ := idb.Conn.(*redis.Client)

	// todo
	ctx := context.Background()
	for _, key := range tbn {
		keyType, err := _db.Type(ctx, key).Result()
		if err != nil {
			log.Println("redis get type: ", err)
			continue
		}

		var value interface{}
		switch keyType {
		case "string":
			valueStringCmd := _db.Get(ctx, key)
			if result, err := valueStringCmd.Result(); err != nil {
				log.Printf("无法获取字符串键 %s 的值：%v\n", key, err)
				continue
			} else {
				value = result
			}

		case "list":
			valueListCmd := _db.LRange(ctx, key, 0, -1)
			if err := valueListCmd.Err(); err != nil {
				log.Printf("无法获取列表键 %s 的值：%v\n", key, err)
				continue
			}

			result, _ := valueListCmd.Result()
			value = result

		case "set":
			valueSetCmd := _db.SMembers(ctx, key)
			if err := valueSetCmd.Err(); err != nil {
				log.Printf("无法获取集合键 %s 的值")
				continue
			}
			result, _ := valueSetCmd.Result()
			value = result
		case "hash":
			result, err := _db.HGetAll(ctx, key).Result()
			if err != nil {
				log.Println("redis get:", err)
				continue
			}
			value = result
		}

		fmt.Println(value)

	}

	return nil
}
