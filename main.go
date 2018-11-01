package main

import (
	"fmt"

	"./initialize"
)

func main() {
	fmt.Println("hello")
	conf := initialize.Run()
	fmt.Println(conf.IPAddresses[0].IP)

}
