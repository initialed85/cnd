package main

import (
	"fmt"
	"github.com/initialed85/cnd/pkg/app"
	"log"
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Printf("%v\n\n", app.Usage)
		log.Fatal(err)
	}
}
