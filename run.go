package viewAdmin

import (
	_ "embed"
	"fmt"
	"github.com/7134g/viewAdmin/internel/router"
	"github.com/7134g/viewAdmin/internel/view"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

//go:embed etc/base.yaml
var baseYaml string

func Cat() {
	fmt.Println(baseYaml)
}

func Run(c *view.Config) {
	gin.SetMode(gin.ReleaseMode)
	run(c)
}

// RunDebug 打开接口日志和数据库日志
func RunDebug(c *view.Config) {
	view.OpenLogDB()
	gin.SetMode(gin.DebugMode)

	run(c)
}

// RunDebugYaml 通过配置打开接口日志和数据库日志
func RunDebugYaml(etc ...string) {
	view.OpenLogDB()
	gin.SetMode(gin.DebugMode)

	var cfgPath string
	if etc == nil || len(etc) < 1 {
		cfgPath = "etc/admin.yaml"
	} else {
		cfgPath = etc[0]
	}

	if f, err := os.Stat(cfgPath); err != nil || f == nil {
		log.Fatal("cfgPath error: ", err)
	}
	c := view.InitConfig(cfgPath)
	run(c)
}

func run(c *view.Config) {
	serve := gin.Default()
	serve.Use(cors.Default())

	router.InitRouter(serve, c)
	log.Printf("Starting admin server at %s...\n", c.Listen)
	log.Println(serve.Run(c.Listen))
}
