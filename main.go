package main

import (
	_ "github.com/lib/pq"

	"github.com/dl-watson/pg-go/controller"
)

func main() {
	controller.SetupServer()
}
