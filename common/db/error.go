package db

import "errors"

var (
	DbNotConnect = errors.New("not redis connect")
	DbNotFound   = errors.New("can not find anything")
)
