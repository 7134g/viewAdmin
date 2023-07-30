package logic

import (
	"github.com/7134g/viewAdmin/config"
	"github.com/7134g/viewAdmin/internel/serve"
)

type Home struct {
	cfg *config.Config
}

func NewHomeLogic(c *config.Config) Home {
	return Home{cfg: c}
}

func (h *Home) Home(_ *serve.BaseContext) (interface{}, error) {
	return "ok", nil
}
