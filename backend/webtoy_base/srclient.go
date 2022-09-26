package webtoy_base

import (
	"errors"
	"strings"
	"time"

	"github.com/MuggleWei/srclient/clb"
	"github.com/MuggleWei/srclient/srd"
	log "github.com/sirupsen/logrus"
)

const (
	ServiceStatus_Pass           = 1
	ServiceStatus_ReadyToOffline = 2
)

type SRClient struct {
	ClientSD      *srd.ServiceDiscoveryClient
	ClientLB      *clb.ClientLoadBalancer
	ServiceStatus int
}

type SRClientArgs struct {
	SrAddr        string
	SrServiceID   string
	SrServiceName string
	SrServiceHost string
	SrServicePort uint
	SrServiceTag  string
	SrServiceTTL  time.Duration
}

func NewSRClient(args *SRClientArgs) (*SRClient, error) {
	client := &SRClient{
		ClientSD:      nil,
		ClientLB:      nil,
		ServiceStatus: ServiceStatus_Pass,
	}

	clientSD, err := srd.NewConsulClient(args.SrAddr)
	if err != nil {
		log.Errorf("failed init service discovery client: %v", err.Error())
		return nil, err
	}

	if args.SrServiceID != "" {
		tags := strings.Split(args.SrServiceTag, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
		registration := srd.ServiceRegistration{
			ID:    args.SrServiceID,
			Name:  args.SrServiceName,
			Addr:  args.SrServiceHost,
			Port:  int(args.SrServicePort),
			Tag:   tags,
			TTL:   args.SrServiceTTL,
			Check: client.Check,
		}

		for {
			log.Infof("try register to service registry")
			err = clientSD.Register(&registration)
			if err != nil {
				log.Errorf("failed init service registry client: %v", err.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
	}

	clientLB := clb.NewClientLoadBalancer(clientSD, args.SrServiceTTL*3)

	client.ClientSD = &clientSD
	client.ClientLB = clientLB

	return client, nil
}

func (this *SRClient) Check() (bool, error) {
	if this.ServiceStatus == ServiceStatus_Pass {
		return true, nil
	}
	return false, errors.New("ready to offline")
}
