package main

import (
	"github.com/gin-gonic/gin"
	"github.com/margostino/gobox/common"
	"github.com/margostino/gobox/factory"
	"net/http"
)

var responseMocks = make(map[string]string)
var healthcheckMocks = make(map[string]string)

func main() {
	var servers = common.GetServerConfig("./configuration/configuration.yml")
	wg := common.WaitGroup(len(servers))
	for _, server := range servers {
		responseMocks[server.GetAddress()] = server.ResponseFile
		healthcheckMocks[server.GetAddress()] = server.HealthcheckFile
		go start(server)
	}
	wg.Wait()
}

func start(server *common.Server) {
	router := gin.Default()
	router.POST(server.Path, response)
	router.POST(server.HealthcheckPath, healthcheck)
	router.Run(server.Host + ":" + server.Port)
}

func response(c *gin.Context) {
	//time.Sleep(4000 * time.Millisecond)
	//success(c)
	failure(c)
}

func success(c *gin.Context)  {
	response, _ := factory.GetPayload(responseMocks[c.Request.Host])
	c.IndentedJSON(http.StatusOK, response)
}

func failure(c *gin.Context)  {
	c.AbortWithStatus(500)
}

func healthcheck(c *gin.Context) {
	response, _ := factory.GetPayload(healthcheckMocks[c.Request.Host])
	c.IndentedJSON(http.StatusOK, response)
}
