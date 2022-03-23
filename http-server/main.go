package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/margostino/gobox/common"
	"github.com/margostino/gobox/factory"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var responseMocks = make(map[string]string)
var healthcheckMocks = make(map[string]string)
var hotStatus = make(map[string]int)

type HotStatusRequest struct {
	Status string `json:"status"`
}

type Request struct {
	Namespace string            `json:"namespace"`
	Variables []string          `json:"variables"`
	Arguments map[string]string `json:"arguments"`
}

func main() {
	var servers = common.GetServerConfig("./configuration/configuration.yml")
	wg := common.WaitGroup(len(servers))
	for _, server := range servers {
		setDefaultHotStatus(server)
		responseMocks[server.GetAddress()] = server.ResponseFile
		healthcheckMocks[server.GetAddress()] = server.HealthcheckFile
		go start(server)
	}
	wg.Wait()
}

func setDefaultHotStatus(server *common.Server) {
	hotStatus[server.Host+":"+server.Port] = common.SuccessEnabled
}

func start(server *common.Server) {
	router := gin.Default()
	router.POST(server.Path, response)
	router.PUT(server.HotStatusPath, updateHotStatus)
	router.GET(server.HealthcheckPath, healthcheck)
	router.GET("/ping", pong)
	router.Run(server.Host + ":" + server.Port)
}

func response(c *gin.Context) {
	var request Request
	jsonRequest, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(jsonRequest, &request)

	log.Println(fmt.Sprintf("Request: %s", string(jsonRequest)))

	if isSuccessEnabled(c) {
		filename := strings.Replace(responseMocks[c.Request.Host], "{0}", request.Namespace, -1)
		response, _ := factory.GetPayload(filename)
		successWithResponse(c, response)
	} else if isFailureEnabled(c) {
		failure(c)
	} else if isRandomnessEnabled(c) {
		randomness(c)
	}
	// TODO: tbd no option
}

// TODO: config thresholds
func randomness(c *gin.Context) {
	randomDelay := rand.Intn(100)
	randomSuccess := rand.Intn(100)
	if randomDelay < 5 {
		time.Sleep(1200 * time.Millisecond)
	}
	if randomSuccess > 70 {
		success(c)
	} else {
		failure(c)
	}
}

func success(c *gin.Context) {
	response, _ := factory.GetPayload(responseMocks[c.Request.Host])
	c.IndentedJSON(http.StatusOK, response)
}

func successWithResponse(c *gin.Context, response *interface{}) {
	c.IndentedJSON(http.StatusOK, response)
}

func failure(c *gin.Context) {
	c.AbortWithStatus(500)
}

func healthcheck(c *gin.Context) {
	response, _ := factory.GetPayload(healthcheckMocks[c.Request.Host])
	c.IndentedJSON(http.StatusOK, response)
}

func pong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func updateHotStatus(c *gin.Context) {

	if status, ok := hotStatus[c.Request.Host]; ok {
		//request := HotStatusRequest{}
		//if c.ShouldBind(&request) == nil {
		//	status = getHotStatusFrom(&request)
		//	hotStatus[c.Request.Host] = status
		//	c.IndentedJSON(http.StatusNoContent, status)
		//}
		param := c.Query("id")
		status = getHotStatusFrom(param)

		if status >= 0 {
			hotStatus[c.Request.Host] = status
			c.IndentedJSON(http.StatusNoContent, status)
		} else {
			c.IndentedJSON(http.StatusBadRequest, "status param is wrong")
		}
	} else {
		c.IndentedJSON(http.StatusNotFound, "config not found")
	}

}

// TODO: define another simpler contract to avoid strings
func getHotStatusFrom(param string) int {
	if param == "success" {
		return common.SuccessEnabled
	}
	if param == "failure" {
		return common.FailureEnabled
	}
	if param == "randomness" {
		return common.RandomnessEnabled
	}
	return -1
}

func isRandomnessEnabled(c *gin.Context) bool {
	value, ok := hotStatus[c.Request.Host]
	return ok && value == common.RandomnessEnabled
}

func isSuccessEnabled(c *gin.Context) bool {
	value, ok := hotStatus[c.Request.Host]
	return ok && value == common.SuccessEnabled
}

func isFailureEnabled(c *gin.Context) bool {
	value, ok := hotStatus[c.Request.Host]
	return ok && value == common.FailureEnabled
}
