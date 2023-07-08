package handle

import (
	"github.com/7134g/viewAdmin/internel/logic"
	"github.com/7134g/viewAdmin/internel/logic/table"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/7134g/viewAdmin/internel/view"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomeHandler(cfg *view.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		baseContext := &serve.BaseContext{ctx}
		home := logic.NewHomeLogic(cfg)
		response, err := home.Home(baseContext)
		if err != nil {
			ctx.JSON(http.StatusOK, fail(err, response))
		} else {
			ctx.JSON(http.StatusOK, success(response))
		}
	}

}

func ViewTableHandler(cfg *view.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		baseContext := &serve.BaseContext{ctx}
		vb := table.NewViewTableLogic(cfg)
		response, err := vb.ViewTable(baseContext)
		if err != nil {
			ctx.JSON(http.StatusOK, fail(err, response))
		} else {
			ctx.JSON(http.StatusOK, success(response))
		}
	}
}

func ListHandler(cfg *view.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		baseContext := &serve.BaseContext{ctx}
		vb := table.NewListLogic(cfg)
		response, err := vb.List(baseContext)
		if err != nil {
			ctx.JSON(http.StatusOK, fail(err, response))
		} else {
			ctx.JSON(http.StatusOK, success(response))
		}
	}
}

func InsertHandler(cfg *view.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		baseContext := &serve.BaseContext{ctx}
		vb := table.NewInsertLogic(cfg)
		response, err := vb.Insert(baseContext)
		if err != nil {
			ctx.JSON(http.StatusOK, fail(err, response))
		} else {
			ctx.JSON(http.StatusOK, success(response))
		}
	}
}

func UpdateHandler(cfg *view.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		baseContext := &serve.BaseContext{ctx}
		vb := table.NewUpdateLogic(cfg)
		response, err := vb.Update(baseContext)
		if err != nil {
			ctx.JSON(http.StatusOK, fail(err, response))
		} else {
			ctx.JSON(http.StatusOK, success(response))
		}
	}
}

func DeleteHandler(cfg *view.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		baseContext := &serve.BaseContext{ctx}
		vb := table.NewDeleteLogic(cfg)
		response, err := vb.Delete(baseContext)
		if err != nil {
			ctx.JSON(http.StatusOK, fail(err, response))
		} else {
			ctx.JSON(http.StatusOK, success(response))
		}
	}
}
