package factory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Payload interface{}

func GetPayload(file string) (*interface{}, *bytes.Buffer) {
	var payload interface{}
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &payload); err != nil {
		fmt.Println("There was an error:", err)
	}
	return &payload, bytes.NewBuffer(byteValue)
}
