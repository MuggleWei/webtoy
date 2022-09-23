package webtoy_base

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type SrClientComponent struct {
	Client *SRClient
}

var (
	instSrClientComponent *SrClientComponent
	onceSrClientComponent sync.Once
)

func GetSrClientComponent() *SrClientComponent {
	if instSrClientComponent == nil {
		onceSrClientComponent.Do(func() {
			instSrClientComponent = &SrClientComponent{
				Client: nil,
			}
		})
	}
	return instSrClientComponent
}

func (this *SrClientComponent) Init(args *SRClientArgs) error {
	client, err := NewSRClient(args)
	if err != nil {
		log.Errorf("failed init srclient: %v", err.Error())
		return err
	}

	this.Client = client

	return nil
}
