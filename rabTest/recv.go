package main

import (
	"fmt"
	rab "rabTest/RabbitMQ"
)

func main()  {
	fmt.Println("消费者")
	mqConn,mqChan,err := rab.Connect()
	if err != nil{
		fmt.Println("connect rabbitmq fail")
		return
	}
	fmt.Println("connect success")
	defer mqConn.Close()
	defer mqChan.Close()
	rab.Consume(mqChan,"test") //消费者接收推送的消息
}
