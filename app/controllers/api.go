package controllers

import (
	"github.com/robfig/revel"
	"awsviz/app/controllers/aws"
)

type Api struct {
	*revel.Controller
}

func (c Api) Everything() revel.Result {
	instances := aws.ListInstances()
	return c.RenderJson(instances)
}

func (c Api) Prod() revel.Result {
	instances := aws.ListInstances()
	return c.RenderJson(instances)
}
