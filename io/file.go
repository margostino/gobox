package io

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func OpenFile(file string) *os.File {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	return jsonFile
}

func ReadAll(payload *os.File) *bytes.Buffer {
	byteValue, _ := ioutil.ReadAll(payload)
	//if err := json.Unmarshal(byteValue, &someStruct); err != nil {
	//	fmt.Println("There was an error:", err)
	//}
	return bytes.NewBuffer(byteValue)
}
