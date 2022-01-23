package nacos

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/iuv/registry-hub/adapter"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var serviceMap = make(map[string][]Service)

func addService(s Service) {
	// add or update
	services := serviceMap[s.ServiceName]
	if len(services) > 0 {
		for i, service := range services {
			if s.InstanceId == service.InstanceId {
				services[i] = s
				return
			}
		}
		serviceMap[s.ServiceName] = append(services, s)
	} else {
		serviceMap[s.ServiceName] = []Service{s}
	}
}

func delService(s Service) {
	// add or update
	services := serviceMap[s.ServiceName]
	if len(services) > 0 {
		for i, service := range services {
			if s.InstanceId == service.InstanceId {
				serviceMap[s.ServiceName] = append(services[:i], services[i+1:]...)
				return
			}
		}
	}
}

// get service list
func ServiceList(c *gin.Context) {
	serviceName := c.Query("serviceName")
	log.Printf(serviceName)
	if serviceName == "" {
		c.String(http.StatusOK, "caused: Param 'serviceName' is required.;")
		return
	}
	services := serviceMap[serviceName]
	var respService RespService

	// rquest nacos get instance list
	nacosResp, err := http.Get(getNacosUrl("/ns/instance/list?serviceName=" + serviceName))
	if err == nil {
		body, _ := ioutil.ReadAll(nacosResp.Body)
		json.Unmarshal(body, &respService)
		// fix adapter proxy
		hosts := respService.Hosts
		for i, host := range hosts {
			hosts[i].Ip = adapter.GetAdapterUrl(c)+host.Ip
		}
	}
	// nacos service not found ,get by local serviceMap
	if len(respService.Hosts) <= 0{
		if len(services) > 0 {
			respService = GetRespServiceByService(services, serviceName)
		} else {
			respService = GetDefaultRespService(serviceName)
		}
	}
	c.JSON(http.StatusOK, respService)
}

func getNacosUrl(path string) string {
	return "http://" + viper.GetString("nacos.url") + "/nacos/v1" + path
}
func Register(c *gin.Context) {
	var s Service
	if err := c.ShouldBindQuery(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.InstanceId = s.Ip + "#" + strconv.Itoa(s.Port) + "#" + s.ClusterName + "#" + s.ServiceName
	addService(s)
	c.String(http.StatusOK, "ok")
}

// service health
func Beat(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
func DeRegister(c *gin.Context) {
	var s Service
	if err := c.ShouldBindQuery(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.InstanceId = s.Ip + "#" + strconv.Itoa(s.Port) + "#" + s.ClusterName + "#" + s.ServiceName
	delService(s)
	c.String(http.StatusOK, "ok")
}
