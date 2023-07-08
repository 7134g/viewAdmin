package logic

import (
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/7134g/viewAdmin/internel/view"
)

type Home struct {
	cfg *view.Config
}

func NewHomeLogic(c *view.Config) Home {
	return Home{cfg: c}
}

func (h *Home) Home(_ *serve.BaseContext) (interface{}, error) {
	return "ok", nil
}
