package main

import (
	"github.com/gin-gonic/gin"
	"github.com/margostino/gobox/common"
	"github.com/margostino/gobox/factory"
	"net/http"
)

var configuration = common.GetServerConfig("./configuration/configuration.yml")

func main() {
	router := gin.Default()
	router.POST(configuration.Path, response)
	router.POST(configuration.HealthcheckPath, healthcheck)
	router.Run("localhost:" + configuration.Port)
}

func response(c *gin.Context) {
	response, _ := factory.GetPayload(configuration.ResponseFile)
	//time.Sleep(4000 * time.Millisecond)
	c.IndentedJSON(http.StatusOK, response)
}

func healthcheck(c *gin.Context) {
	response, _ := factory.GetPayload(configuration.HealthcheckFile)
	//time.Sleep(4000 * time.Millisecond)
	c.IndentedJSON(http.StatusOK, response)
}
