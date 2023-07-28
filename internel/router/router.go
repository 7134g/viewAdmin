package router

import (
	_ "embed"
	"github.com/7134g/viewAdmin/internel/handle"
	"github.com/7134g/viewAdmin/internel/view"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	//go:embed dist/index.html
	html []byte
	//go:embed dist/assets/index-dac6ce4d.css
	css []byte
	//go:embed dist/assets/index-e2a54b59.js
	js []byte
)

func InitRouter(r *gin.Engine, c *view.Config) {
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", html)
	})
	r.GET("/assets/index-dac6ce4d.css", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/css; charset=utf-8", css)
	})
	r.GET("/assets/index-e2a54b59.js", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/javascript", js)
	})

	baseApi := r.Group("/")

	//baseApi.GET("/", handle.HomeHandler(c))

	modelGroup := baseApi.Group("model")
	{
		modelGroup.GET("/tables", handle.ViewTableHandler(c))
		// 查询该数据表的数据
		modelGroup.POST("/list", handle.ListHandler(c))
		// 插入一条数据
		modelGroup.POST("/insert", handle.InsertHandler(c))
		// 更新某条数据
		modelGroup.PUT("/update", handle.UpdateHandler(c))
		// 删除某条数据
		modelGroup.DELETE("/delete", handle.DeleteHandler(c))
	}
}
