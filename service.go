package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"math/rand"
)

const LanAddrKey = "lan_ipv4"
const WanAddrKey = "wan_ipv4"

type Service struct {
	ID             string
	Schema         string
	Name           string
	Address        string
	MatchBody      string
	CheckInterval  string
	Port           int
	Tags           []string
	Meta           map[string]string
	HealthEndpoint string
	ServiceAddress map[string]api.ServiceAddress
}

func (s *Service) AddTags(tags ...string) {
	s.Tags = append(s.Tags, tags...)
}

func (c *client) GetServices(id, tag string) ([]*api.AgentService, error) {

	ss, _, err := c.client.Health().Service(id, tag, true, c.queryOptions(nil))
	if nil != err {
		return nil, err
	}
	if len(ss) == 0 {
		return nil, errors.New("service not found")
	}

	services := make([]*api.AgentService, 0)
	for e := range ss {
		services = append(services, ss[e].Service)
	}
	return services, nil
}

func (c *client) GetService(id, tags string) (*api.AgentService, error) {
	ss, err := c.GetServices(id, tags)
	if nil != err {
		return nil, err
	}
	return ss[rand.Intn(len(ss))], nil
}

func (c *client) GetServiceAddrPort(id string, useLan bool, tags string) (host string, port int, err error) {
	s, err := c.GetService(id, tags)
	if nil != err {
		return "", 0, err
	}
	var addr api.ServiceAddress
	var ok bool
	if useLan {
		addr, ok = s.TaggedAddresses[LanAddrKey]
	} else {
		addr, ok = s.TaggedAddresses[WanAddrKey]
	}

	if !ok {
		return "", 0, errors.New("service not found")
	}

	return addr.Address, addr.Port, nil
}
