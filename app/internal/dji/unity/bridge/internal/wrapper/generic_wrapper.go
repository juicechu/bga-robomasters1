package wrapper

type generic struct {
	callbacks map[uint64]EventCallback
	parent    Wrapper
}

func newGeneric(parent Wrapper) *generic {
	return &generic{
		make(map[uint64]EventCallback),
		parent,
	}
}

func (g *generic) UnitySetEventCallback(eventCode uint64,
	eventCallback EventCallback) {
	eventType := eventCode >> 32

	add := false
	if eventCallback == nil {
		delete(g.callbacks, eventType)
	} else {
		g.callbacks[eventType] = eventCallback
		add = true
	}

	g.parent.unitySetEventCallback(eventCode, add)
}

func (g *generic) Callback(eventCode uint64) EventCallback {
	eventType := eventCode >> 32

	cb, ok := g.callbacks[eventType]
	if !ok {
		return nil
	}

	return cb
}
