package viewAdmin

import (
	"github.com/7134g/viewAdmin/internel/view"
	"testing"
)

func TestRun(t *testing.T) {
	c := &view.Config{
		Mysql: view.MysqlClient{
			DBName:   "blog",
			Username: "root",
			Password: "mysql",
			Address:  "127.0.0.1:3306",
		},
		Sqlite: view.SqliteClient{
			Link: "./admin.db",
		},
		Listen: "127.0.0.1:10086",
	}
	Run(c)
}

func TestRunDebugYaml(t *testing.T) {
	RunDebugYaml("etc/admin.yaml")
}

func TestCat(t *testing.T) {
	Cat()
}
