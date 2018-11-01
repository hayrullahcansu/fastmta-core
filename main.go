package main

import (
	"fmt"

	"./entity"
	"./initialize"
)

func main() {
	fmt.Println("hello")
	conf := initialize.Run()
	fmt.Println(conf.IPAddresses[0].IP)
	entity.Run()

}
