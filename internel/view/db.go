package view

type MysqlClient struct {
	DBName   string
	Username string
	Password string
	Address  string
}

type SqliteClient struct {
	Link   string
	DBName string
}

type MongoClient struct {
	Uri    string
	DBName string
}

type DBConnect struct {
	DBName string
	Script string
	Conn   interface{}
}
