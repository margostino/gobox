package main

import (
	"fmt"
	"github.com/marianogappa/sqlparser"
	"log"
)

func main() {
	query, err := sqlparser.Parse("SELECT a, b, c FROM 'd' WHERE e = '1' AND f > '2'")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+#v", query)
}
