package main

import (
	"github.com/gin-gonic/gin"
	"github.com/margostino/gobox/common"
	"github.com/margostino/gobox/factory"
	"net/http"
	"strings"
)

var responseMocks = make(map[string]string)
var healthcheckMocks = make(map[string]string)
var hotStatus = make(map[string]int)
var isTimeout bool

type Request struct {
	Namespace string            `json:"namespace"`
	Variables []string          `json:"variables"`
	Arguments map[string]string `json:"arguments"`
}

func main() {
	var server = common.GetDynamicTimerConfig("./configuration/configuration.yml")
	responseMocks[server.GetAddress()] = server.ResponseFile
	start(server)
}

func start(server *common.DynamicTimer) {
	router := gin.Default()
	router.GET(server.Path, response)
	router.Run(server.Host + ":" + server.Port)
}

func response(c *gin.Context) {
	filename := strings.Replace(responseMocks[c.Request.Host], "{0}", "request.Namespace", -1)
	response, _ := factory.GetPayload(filename)
	c.IndentedJSON(http.StatusOK, response)
}
