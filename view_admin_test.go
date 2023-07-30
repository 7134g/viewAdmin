package viewAdmin

import (
	"github.com/7134g/viewAdmin/config"
	"testing"
)

func TestRun(t *testing.T) {
	c := &config.Config{
		Mysql: config.MysqlClient{
			DBName:   "blog",
			Username: "root",
			Password: "mysql",
			Address:  "127.0.0.1:3306",
		},
		Sqlite: config.SqliteClient{
			Link: "./admin.db",
		},
		Listen: "127.0.0.1:10086",
	}
	Run(c)
}

func TestRunDebugYaml(t *testing.T) {
	RunDebugYaml()
}

func TestCat(t *testing.T) {
	Cat()
}
