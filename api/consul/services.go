package consul

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"github.com/iuv/registry-hub/adapter"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
)

var serviceMap = make(map[string][]ServiceVo)
var serviceIdMap = make(map[string]ServiceVo)

// add cache service
func addService(s ServiceVo) {
	oldS := serviceIdMap[s.Id]
	// update or add to ids map
	serviceIdMap[s.Id] = s
	// name changed remove old service in serviceMap
	if oldS.Name != s.Name {
		oldSs := serviceMap[oldS.Name]
		for i, old := range oldSs {
			if old.Id == s.Id {
				serviceMap[oldS.Name] = append(oldSs[:i], oldSs[i+1:]...)
			}
		}
	}
	// add or update serviceMap
	ss := serviceMap[s.Name]
	if len(ss) == 0 {
		serviceMap[s.Name] = []ServiceVo{s}
	} else {
		for i, t := range ss {
			// id equals update and return
			if t.Id == s.Id {
				ss[i] = s
				return
			}
		}
		// id is exists add to arrays
		serviceMap[s.Name] = append(ss, s)
	}
}

// del cache service by serviceId
func delService(serviceId string) {
	s := serviceIdMap[serviceId]
	delete(serviceIdMap, serviceId)
	ss := serviceMap[s.Name]
	for i, t := range ss {
		if t.Id == serviceId {
			serviceMap[s.Name] = append(ss[:i], ss[i+1:]...)
			return
		}
	}
}

// get service by name
func GetServices(name string) []ServiceVo {
	return serviceMap[name]
}

// get service by id
func GetServiceById(serviceId string) (ServiceVo, bool) {
	vo, ok := serviceIdMap[serviceId]
	return vo, ok
}

// service register
func Register(c *gin.Context) {
	var serviceVo ServiceVo
	if err := c.ShouldBindJSON(&serviceVo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	addService(serviceVo)
	c.Status(http.StatusOK)
}

// service deregister
func DeRegister(c *gin.Context) {
	serviceId := c.Param("serviceId")
	delService(serviceId)
}

// catalog services by serviceName
func GetServiceByName(c *gin.Context) {
	name := c.Param("name")
	services := GetServices(name)
	serviceResps := make([]ServiceResp, 0)
	for _, s := range services {
		var sr ServiceResp
		sr.ID = uuid.New()
		sr.ServiceName = s.Name
		sr.ServiceID = s.Id
		sr.ServiceAddress = s.Address
		sr.ServicePort = s.Port
		sr.Datacenter = "dc1"
		sr.Node = "node1"
		var ipport IpPort
		ipport.Address = s.Address
		ipport.Port = s.Port
		var addr Addr
		addr.Lan = ipport
		addr.Wan = ipport
		sr.ServiceTaggedAddresses = addr
		serviceResps = append(serviceResps, sr)
	}
	c.JSON(http.StatusOK, serviceResps)
}

// /health/service/:name
func GetHealthServiceByName(c *gin.Context) {
	name := c.Param("name")
	services := GetServices(name)
	var resps = make([]HealthServiceResp, 0)
	// get consul service by name
	consulResp, err := http.Get(getConsulUrl("/health/service/" + name))
	if err == nil {
		body, _ := ioutil.ReadAll(consulResp.Body)
		json.Unmarshal(body, &resps)
	}
	// fix address proxy
	if len(resps) > 0 {
		for i, s := range resps {
			resps[i].Service.Address = adapter.GetAdapterUrl(c) + s.Service.Address
		}
	}
	// add local service
	for _, s := range services {
		resp := GetHealthServiceResp(s)
		resps = append(resps, resp)
	}
	c.JSON(http.StatusOK, resps)
}
func getConsulUrl(path string) string {
	return "http://" + viper.GetString("consul.url") + "/v1" + path
}

// catalog all services
func GetAllServices(c *gin.Context) {
	wait := c.Query("wait")
	if wait != "" {
		time.Sleep(2 * time.Second)
	}
	retMap := make(map[string][]string)
	// get consul services
	resp, err := http.Get(getConsulUrl("/catalog/services"))
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &retMap)
	}
	// add local services
	rett := []string{"secure=false"}
	for k, _ := range serviceMap {
		retMap[k] = rett
	}
	c.JSON(http.StatusOK, retMap)
}
