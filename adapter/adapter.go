package adapter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// get adapter url
func GetAdapterUrl(c *gin.Context) string {
	var adapterUrl string
	if viper.GetString("local.url") != "" {
		adapterUrl = viper.GetString("local.url") + "/adapter/"
	} else {
		host := c.Request.Host
		adapterUrl = host + "/adapter/";
	}
	return adapterUrl
}

// proxy request
func Adapter(c *gin.Context) {
	api := c.Param("api")
	println(c.Request.Method + "#" + api)
	// /ip:port/api
	target := api[1:]
	targets := strings.Split(target, "/")
	target = targets[0]
	path := strings.Join(targets[1:], "/")
	u := &url.URL{}
	u.Scheme = "http"
	u.Host = target
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		log.Printf("consul adapter error :%v", err)
		ret := fmt.Sprintf("consul adapter error :%v", err)
		writer.Write([]byte(ret))
	}
	request := c.Request
	request.URL.Path = path
	proxy.ServeHTTP(c.Writer, request)
}
