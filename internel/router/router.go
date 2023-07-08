package router

import (
	"github.com/7134g/viewAdmin/internel/handle"
	"github.com/7134g/viewAdmin/internel/view"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, c *view.Config) {
	baseApi := r.Group("/")

	baseApi.GET("/", handle.HomeHandler(c))

	modelGroup := baseApi.Group("model")
	{
		modelGroup.GET("/tables", handle.ViewTableHandler(c))
		// 查询该数据表的数据
		modelGroup.POST("/", handle.ListHandler(c))
		// 插入一条数据
		modelGroup.POST("/insert", handle.InsertHandler(c))
		// 更新某条数据
		modelGroup.PUT("/update", handle.UpdateHandler(c))
		// 删除某条数据
		modelGroup.DELETE("/delete", handle.DeleteHandler(c))
	}
}
