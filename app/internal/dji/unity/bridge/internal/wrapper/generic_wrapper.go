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
	add := false
	if eventCallback == nil {
		delete(g.callbacks, eventCode)
	} else {
		g.callbacks[eventCode] = eventCallback
		add = true
	}

	g.parent.unitySetEventCallback(eventCode, add)
}

func (g *generic) Callback(eventCode uint64) EventCallback {
	cb, ok := g.callbacks[eventCode]
	if !ok {
		return nil
	}

	return cb
}
