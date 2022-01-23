package nacos

import (
	"strings"
	"time"
)

// register model
type Service struct {
	InstanceId  string `json:"instanceId"`
	Ip          string `json:"ip" form:"ip"`
	Port        int    `json:"port" form:"port"`
	ServiceName string `json:"serviceName" form:"serviceName"`
	ClusterName string `json:"clusterName" form:"clusterName"`
	App         string `json:"App" form:"App"`
	NamespaceId string `json:"namespaceId" form:"namespaceId"`
}

// response service
type RespService struct {
	Name                     string            `json:"name"`
	GroupName                string            `json:"groupName"`
	Clusters                 string            `json:"clusters"`
	CacheMillis              int               `json:"cacheMillis"`
	LastRefTime              int64             `json:"lastRefTime"`
	Checksum                 string            `json:"checksum"`
	AllIPs                   bool              `json:"allIPs"`
	ReachProtectionThreshold bool              `json:"reachProtectionThreshold"`
	Valid                    bool              `json:"valid"`
	Hosts                    []RespServiceHost `json:"hosts"`
}
type RespServiceHost struct {
	InstanceId                string  `json:"instanceId"`
	Ip                        string  `json:"ip"`
	Port                      int     `json:"port"`
	Weight                    float32 `json:"weight"`
	Healthy                   bool    `json:"healthy"`
	Enabled                   bool    `json:"enabled"`
	Ephemeral                 bool    `json:"ephemeral"`
	ClusterName               string  `json:"clusterName"`
	ServiceName               string  `json:"serviceName"`
	InstanceHeartBeatInterval int     `json:"instanceHeartBeatInterval"`
	InstanceHeartBeatTimeOut  int     `json:"instanceHeartBeatTimeOut"`
	IpDeleteTimeout           int     `json:"ipDeleteTimeout"`
}

func GetDefaultRespService(serviceName string) RespService {
	var resp RespService
	resp.Name = serviceName
	ss := strings.Split(serviceName, "@@")
	resp.GroupName = ss[0]
	resp.CacheMillis = 10000
	resp.LastRefTime = time.Now().UnixNano() / 1e6
	resp.Valid = true
	resp.Hosts = make([]RespServiceHost, 0)
	return resp
}
func GetRespServiceByService(ss []Service, serviceName string) RespService {
	var resp RespService
	resp.Name = serviceName
	sns := strings.Split(serviceName, "@@")
	resp.GroupName = sns[0]
	resp.CacheMillis = 10000
	resp.LastRefTime = time.Now().UnixNano() / 1e6
	resp.Valid = true
	resp.Hosts = make([]RespServiceHost, 0)
	for _, s := range ss {
		resp.Hosts = append(resp.Hosts, getRespServiceHost(s))
	}
	return resp
}

func getRespServiceHost(s Service) RespServiceHost {
	var host RespServiceHost
	host.InstanceId = s.InstanceId
	host.Ip = s.Ip
	host.Port = s.Port
	host.ServiceName = s.ServiceName
	host.ClusterName = s.ClusterName
	host.Weight = 1.0
	host.Healthy = true
	host.Enabled = true
	host.Ephemeral = true
	host.InstanceHeartBeatTimeOut = 15000
	host.InstanceHeartBeatInterval = 5000
	host.IpDeleteTimeout = 30000
	return host
}
