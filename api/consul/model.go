package consul

import (
	"github.com/go-basic/uuid"
)

/**
{
  "id": "web1",
  "name": "web",
  "port": 80,
  "check": {
    "name": "ping check",
    "args": ["ping", "-c1", "learn.hashicorp.com"],
    "interval": "30s",
    "status": "passing"
  }
}
*/
// service register json
type ServiceVo struct {
	Id      string  `json:"id" from:"id"`
	Name    string  `json:"name"`
	Port    int     `json:"port"`
	Address string  `json:"address"`
	Check   CheckVo `json:"check"`
}
type CheckVo struct {
	Name     string   `json:"name"`
	Args     []string `json:"args"`
	Interval string   `json:"interval"`
	Status   string   `json:"status"`
}

// catalog service response json
type ServiceResp struct {
	ID                     string
	Node                   string
	Address                string
	Datacenter             string
	ServiceAddress         string
	ServiceID              string
	ServiceName            string
	ServicePort            int
	ServiceTaggedAddresses Addr
}
type Addr struct {
	Lan IpPort `json:"lan"`
	Wan IpPort `json:"wan"`
}
type IpPort struct {
	Address string `json:"address"`
	Port    int    `json:"port""`
}

// health service response
type HealthServiceResp struct {
	Node    HealthServiceRespNode
	Service HealthServiceRespService
	Checks  []HealthServiceRespChecks
}

// init GetHealthServiceResponse by ServiceVo
func GetHealthServiceResp(s ServiceVo) HealthServiceResp {
	var node HealthServiceRespNode
	node.ID = uuid.New()
	node.Node = "node1"
	node.Address = "127.0.0.1"
	node.Datacenter = "dc1"
	var ta TaggedAddresses
	ta.Lan = node.Address
	ta.Wan = node.Address
	ta.LanIpv4 = node.Address
	ta.WanIpv4 = node.Address
	node.TaggedAddresses = ta
	var service HealthServiceRespService
	service.ID = s.Id
	service.Service = s.Name
	service.Port = s.Port
	adapterAddress := s.Address
	service.Address = adapterAddress
	var ipPort = IpPort{adapterAddress, s.Port}
	var healthAddr = HealthAddr{ipPort, ipPort}
	service.TaggedAddresses = healthAddr
	var check = HealthServiceRespChecks{"node1", "serfHealth", "Serf Health Status", "passing"}
	var checks []HealthServiceRespChecks = make([]HealthServiceRespChecks, 0)
	checks = append(checks, check)
	return HealthServiceResp{node, service, checks}
}

type HealthServiceRespNode struct {
	ID              string
	Node            string
	Address         string
	Datacenter      string
	TaggedAddresses TaggedAddresses
}
type HealthServiceRespService struct {
	ID              string
	Service         string
	Address         string
	TaggedAddresses HealthAddr
	Port            int
}
type HealthServiceRespChecks struct {
	Node    string
	CheckID string
	Name    string
	Status  string
}
type HealthAddr struct {
	LanIpv4 IpPort `json:"lan_ipv4"`
	WanIpv4 IpPort `json:"wan_ipv4"`
}

type TaggedAddresses struct {
	Lan     string `json:"lan"`
	LanIpv4 string `json:"lan_ipv4"`
	Wan     string `json:"wan"`
	WanIpv4 string `json:"wan_ipv4"`
}
