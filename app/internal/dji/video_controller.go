package dji

import (
	"fmt"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge"
)

type VideoController struct {
	eventHandlerIndexes []int
}

func NewVideoController() (*VideoController, error) {
	b := bridge.Instance()

	vc := &VideoController{}

	eventHandlerIndexes := make([]int, 3)

	var err error
	eventHandlerIndexes[0], err = b.AddEventHandler(
		unity.EventTypeGetNativeTexture, vc)
	if err != nil {
		return nil, err
	}
	eventHandlerIndexes[1], err = b.AddEventHandler(
		unity.EventTypeVideoTransferSpeed, vc)
	if err != nil {
		return nil, err
	}
	eventHandlerIndexes[2], err = b.AddEventHandler(
		unity.EventTypeVideoDataRecv, vc)
	if err != nil {
		return nil, err
	}

	vc.eventHandlerIndexes = eventHandlerIndexes

	return vc, nil
}

func (v *VideoController) StartVideo() {
	ub := bridge.Instance()

	ub.SendEvent(unity.NewEvent(unity.EventTypeStartVideo))
}

func (v *VideoController) StopVideo() {
	ub := bridge.Instance()

	ub.SendEvent(unity.NewEvent(unity.EventTypeStopVideo))
}

func (v *VideoController) HandleEvent(event *unity.Event, info []byte,
	tag uint64) {
	fmt.Printf("%s, %#+v, %d\n", unity.EventTypeName(event.Type()), info,
		tag)
}
