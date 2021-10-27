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

var configuration = common.GetClientConfig("./configuration/configuration.yml")

func main() {
	data := io.OpenFile(configuration.RequestFile)
	callsNumber := 20
	wg := waitGroup(callsNumber)
	loop(&wg, callsNumber, data)
	wg.Wait()
}

func loop(wg *sync.WaitGroup, callsNumber int, data *os.File) {
	for i := 0; i < callsNumber; i++ {
		go call(wg, i, data)
	}
}

func call(wg *sync.WaitGroup, requestId int, data *os.File) {
	payload := io.ReadAll(data)
	client := http.GetClient()
	request := http.GetRequest(configuration.Url, payload)
	start := time.Now()
	//RegisterTime("Request", requestId)
	response := http.Call(client, request)
	end := time.Now()
	fmt.Printf("Request #%d Elapsed time %s\n", requestId, end.Sub(start).String())
	http.Print(response)
	wg.Done()
}

func waitGroup(delta int) sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(delta)
	return wg
}
