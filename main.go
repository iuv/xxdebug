package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iuv/registry-hub/adapter"
	"github.com/iuv/registry-hub/api/nacos"
	"github.com/iuv/registry-hub/config"
	"github.com/spf13/viper"
)
import "github.com/iuv/registry-hub/api/consul"

func main() {
	// read config
	err := config.Init()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	gin.SetMode(viper.GetString("runmode"))
	consulGroup := r.Group("v1")
	{
		// agent
		consulGroup.PUT("/agent/service/register", consul.Register)
		consulGroup.PUT("/agent/service/deregister/:serviceId", consul.DeRegister)
		consulGroup.PUT("/agent/check/pass/:checkId", consul.Pass)
		// catalog
		consulGroup.GET("/catalog/services", consul.GetAllServices)
		consulGroup.GET("/catalog/service/:name", consul.GetServiceByName)
		// health
		consulGroup.GET("/health/service/:name", consul.GetHealthServiceByName)
		// kv
		consulGroup.GET("/kv/:all/:name", consul.Kv)
		consulGroup.GET("/kv/:all", consul.KvAll)
	}
	nacosGroup := r.Group("nacos/v1")
	{
		nacosGroup.POST("/ns/instance", nacos.Register)
		nacosGroup.DELETE("/ns/instance", nacos.DeRegister)
		nacosGroup.GET("/ns/instance/list", nacos.ServiceList)
		nacosGroup.PUT("/ns/instance/beat", nacos.Beat)
	}
	adapterGroup := r.Group("adapter")
	{
		adapterGroup.Any("/*api", adapter.Adapter)
	}
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome Registry-hub")
	})
	r.Run(":80") // listen and serve on 0.0.0.0:8080

}
