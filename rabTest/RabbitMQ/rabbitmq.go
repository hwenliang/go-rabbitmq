package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//初始化参数
const(
	User string = "admin"
	Pass string = "admin"
	Host string = "127.0.0.1"
	Port string = "5672"
)

func Test(){
	fmt.Println("Test RabbitMQ")

}
//连接rabbitmq
func Connect() (*amqp.Connection,*amqp.Channel,error){
	//构造链接URL
	rabUrl := fmt.Sprintf("amqp://%s:%s@%s:%s%s", User, Pass, Host, Port, "/")
	fmt.Printf("url:%s\n",rabUrl)
	//创建链接
	mqConn, err := amqp.Dial(rabUrl)
	failOnError(err,"创建链接失败！")
	//defer mqConn.Close()
	if err != nil{
		fmt.Printf("RabbitMQ 打开链接失败：%s\n",err)
		return nil, nil, err
	}
	//打开一个通道
	mqChan, _err := mqConn.Channel()
	failOnError(_err,"打开管道失败！")
	//defer mqChan.Close()
	if _err != nil{
		fmt.Printf("RabbitMQ 打开管道失败：%s\n",_err)
		return nil, nil,  _err
	}
	//mqChan.Ack()
	return mqConn,mqChan,nil
}
//关闭RabbitMQ 连接
func Close(mqConn *amqp.Connection,mqChan *amqp.Channel) error {
	if mqConn != nil{
		err := mqConn.Close()
		return err
	}
	if mqChan != nil{
		err := mqChan.Close()
		return err
	}
	return nil
}
//推送消息
//qName 队列名称
//mqChan 信道指针
//body 要发送的消息
func Publish(mqChan *amqp.Channel,qName string, body string) error{
	//创建交换器 fanout
	err := mqChan.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	failOnError(err,"failed declare a Exchange")
	//创建队列
	q, err := mqChan.QueueDeclare(qName,false,false,false,false,nil)
	failOnError(err,"failed declare a queue")
	//生产者推送消息
	err = mqChan.Publish("logs",q.Name,false,false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(body),
		})
	failOnError(err,"failed Publish")
	fmt.Printf("send:%s\n",body)
	return err
}
//创建消费者并打印接收到的消息
//qName 队列名称
//mqChan 信道指针
func Consume(mqChan *amqp.Channel,qName string) error{
	//创建交换器 fanout
	err := mqChan.ExchangeDeclare("logs","fanout",true, false,false, false, nil)
	failOnError(err,"failed declare a Exchange")
	//创建队列
	q, err := mqChan.QueueDeclare(qName,false,false,false,false,nil)
	failOnError(err,"failed declare a queue")
	//绑定队列
	err = mqChan.QueueBind(q.Name,"","logs",false,nil)
	failOnError(err,"Failed to bind a queue")
	//创建消费者
	msgs, err := mqChan.Consume(q.Name,"",true,false,false,false,nil)
	failOnError(err,"create Consume failed")
	res := make(chan  bool)
	go func() {
		for m := range msgs{
			fmt.Printf("[x]recv:%s\n",m.Body)
		}
	}()
	<-res
	return err
}

//错误日志打印
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}