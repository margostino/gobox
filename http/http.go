package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func GetClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func GetRequest(url string, payload *bytes.Buffer) *http.Request {
	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	//id := uuid.New()
	//request.Header.Add("Id", id.String())
	request.Header.Add("Content-Type", "application/json")
	request.SetBasicAuth("username1", "password")
	return request
}

func Call(client *http.Client, request *http.Request) *http.Response {
	response, error := client.Do(request)
	if error != nil {
		fmt.Println(error.Error())
		//log.Fatal(error)
	}
	return response
}

func Print(response *http.Response) {
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
	}
}
