package main

import "github.com/davidsbond/lux/examples/service/service"

func main() {
	service, err := service.New()

	if err != nil {
		panic(err)
	}

	service.Start()
}
