package main

import (
	"fmt"
	rab "rabTest/RabbitMQ"
)

func main()  {
	fmt.Println("生产者")
	mqConn,mqChan,err := rab.Connect()
	if err != nil{
		fmt.Println("connect rabbitmq fail")
		return
	}
	fmt.Println("connect success")
	defer mqConn.Close()
	defer mqChan.Close()
	rab.Publish(mqChan,"test","hello word")
}
