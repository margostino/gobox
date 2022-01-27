package main

import (
	"fmt"
	"github.com/margostino/gobox/common"
	"github.com/margostino/gobox/http"
	"github.com/margostino/gobox/io"
	"os"
	"sync"
	"time"
)

func main() {
	var clients = common.GetClientConfig("./configuration/configuration.yml")
	callsNumber := 1 // TODO: get from config
	wg := common.WaitGroup(callsNumber)
	for _, client := range clients {
		data := io.OpenFile(client.RequestFile)
		go call(wg, callsNumber, data, client.Url)
	}
	wg.Wait()
}

func call(wg *sync.WaitGroup, callsNumber int, data *os.File, url string) {
	//var waitTime time.Duration
	for i := 0; i < callsNumber; i++ {
		go execute(wg, i, data, url)
		//waitTime = time.Duration(rand.Intn(1) + 1000)
		//time.Sleep(1000 * time.Millisecond)
	}
}

func execute(wg *sync.WaitGroup, requestId int, data *os.File, url string) {
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
