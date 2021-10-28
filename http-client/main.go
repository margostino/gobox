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
	callsNumber := 20
	wg := common.WaitGroup(callsNumber)
	for _, client := range clients {
		data := io.OpenFile(client.RequestFile)
		go call(wg, callsNumber, data, client.Url)
	}
	wg.Wait()
}

func call(wg *sync.WaitGroup, callsNumber int, data *os.File, url string) {
	for i := 0; i < callsNumber; i++ {
		go execute(wg, i, data, url)
	}
}

func execute(wg *sync.WaitGroup, requestId int, data *os.File, url string) {
	payload := io.ReadAll(data)
	client := http.GetClient()
	request := http.GetRequest(url, payload)
	start := time.Now()
	//RegisterTime("Request", requestId)
	response := http.Call(client, request)
	end := time.Now()
	fmt.Printf("Request #%d Elapsed time %s\n", requestId, end.Sub(start).String())
	http.Print(response)
	wg.Done()
}
