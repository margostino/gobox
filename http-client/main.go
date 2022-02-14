package main

import (
	"fmt"
	"github.com/margostino/gobox/common"
	"github.com/margostino/gobox/http"
	"github.com/margostino/gobox/io"
	"log"
	"math/rand"
	"sync"
	"time"
)

var wg *sync.WaitGroup

func main() {
	var clients = common.GetClientConfig("./configuration/configuration.yml")
	delta := getDelta(clients)
	wg = common.WaitGroup(delta)
	for _, client := range clients {
		data, err := io.OpenFile(client.RequestFile)
		if err != nil {
			log.Println(fmt.Sprintf("Cannot open file %s", client.RequestFile), err)
			wg.Add(-client.CallsNumber)
		} else {
			go call(client.CallsNumber, client.MaxStepTime, data, client.Url)
		}
	}
	wg.Wait()
}

func getDelta(clients []*common.Client) int {
	var delta = 0
	for _, client := range clients {
		delta += client.CallsNumber
	}
	return delta
}

func call(callsNumber int, maxStepTime int, data []byte, url string) {
	for i := 0; i < callsNumber; i++ {
		go execute(i, data, url)
		wait(maxStepTime)
	}
}

func wait(maxStepTime int) {
	waitTime := time.Duration(rand.Intn(1) + maxStepTime)
	time.Sleep(waitTime * time.Millisecond)
}

func execute(requestId int, data []byte, url string) {
	payload := io.ReadAll(data)
	client := http.GetClient()
	request := http.GetRequest(url, payload)
	start := time.Now()
	//RegisterTime("Request", requestId)
	response := http.Call(client, request)
	if response != nil {
		end := time.Now()
		fmt.Printf("URL %s Request #%d Elapsed time %s with status: %s\n", url, requestId, end.Sub(start).String(), response.Status)
		http.Print(response)
	}
	wg.Done()
}
