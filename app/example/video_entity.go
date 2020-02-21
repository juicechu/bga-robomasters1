package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type VideoEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}
