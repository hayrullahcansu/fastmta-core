package main

import (
	"fmt"

	"./boss"
	"./initialize"
)

func main() {
	fmt.Println("hello")
	conf := initialize.Run()
	boss := boss.New(conf)
	boss.Run()
	//entity.Run()
	var e int

	fmt.Scanf("%#X", &e)

}
