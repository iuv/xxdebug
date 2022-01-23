package consul

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// check/pass/check_id
func Pass(c *gin.Context) {
	checkId := c.Param("checkId")
	checkIds := strings.Split(checkId, ":")
	_, ok := GetServiceById(checkIds[1])
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusInternalServerError)
	}
}
