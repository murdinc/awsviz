package controllers

import (
	"github.com/robfig/revel"
	"awsgraph/app/controllers/aws"
)

type Api struct {
	*revel.Controller
}

func (c Api) Prod() revel.Result {
	instances := aws.ListInstances()
	return c.RenderJson(instances)
}
