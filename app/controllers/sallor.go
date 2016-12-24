package controllers

import (
	"github.com/revel/revel"
	"wuhuaguo.com/whgv01/app/routes"
)

type Sellers struct {
	Application
}

func (c Sellers) Index() revel.Result {
	return c.Redirect(routes.Application.Index())
}
