package consul

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Kv(c *gin.Context) {
	//serviceName := c.Param("name")
	wait := c.Query("wait")
	if wait != "" {
		time.Sleep(55 * time.Second)
	}
	c.Status(http.StatusNotFound)
}
func KvAll(c *gin.Context) {
	//serviceName := c.Param("name")
	wait := c.Query("wait")
	if wait != "" {
		time.Sleep(55 * time.Second)
	}
	c.Status(http.StatusNotFound)
}
