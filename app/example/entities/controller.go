package entities

import (
	"git.bug-br.org.br/bga/robomasters1/app/example2/components"
	"github.com/EngoEngine/ecs"
)

type Controller struct {
	*ecs.BasicEntity
	*components.Controller
}
