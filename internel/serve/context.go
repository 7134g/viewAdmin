package serve

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type BaseContext struct {
	*gin.Context
}

func wrap(f func(content *BaseContext)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f(&BaseContext{Context: ctx})
	}
}

func (b *BaseContext) GetQueryInt(key string) int {
	value, exist := b.Context.GetQuery(key)
	if !exist {
		return 0
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return v
}
