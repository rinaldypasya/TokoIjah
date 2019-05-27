package main

import (
	"github.com/rinaldypasya/TokoIjah/api"
	"github.com/rinaldypasya/TokoIjah/model"
)

var (
	port = ":8080"
)

func main() {
	db := model.InitDB()
	server := api.InitRouter(db)
	server.Run(port)
}
