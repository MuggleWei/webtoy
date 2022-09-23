package webtoy_base

import (
	"sync"
)

type SessionComponent struct {
	Handler *SessionHandler
}

var (
	instSessionComponent *SessionComponent
	onceSessionComponent sync.Once
)

func GetSessionComponent() *SessionComponent {
	if instSessionComponent == nil {
		onceSessionComponent.Do(func() {
			instSessionComponent = &SessionComponent{
				Handler: NewSessionHandler(nil, 0),
			}
		})
	}
	return instSessionComponent
}
