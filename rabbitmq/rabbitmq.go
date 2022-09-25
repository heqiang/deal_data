package rabbitmq

import (
	"deal_data/config"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// MQ_URL 格式 amqp://账号：密码@rabbitmq服务器地址：端口号/vhost (默认是5672端口)
// 端口可在 /etc/rabbitmq/rabbitmq-env.conf 配置文件设置，也可以启动后通过netstat -tlnp查看
const (
	MQURL = "amqp://guest:guest@127.0.0.1:5673"
)

type RabbitMQ struct {
	Conn       *amqp.Connection
	Channel    *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string
	MqUrl      string
}

// NewRabbitMQ 创建结构体实例
func NewRabbitMQ(config config.MqConfig) *RabbitMQ {
	rabbitMQ := RabbitMQ{
		QueueName: config.QueueName,
		Exchange:  config.ExchangeName,
		MqUrl:     MQURL,
	}
	var err error
	//创建rabbitmq连接
	rabbitMQ.Conn, err = amqp.Dial(rabbitMQ.MqUrl)
	checkErr(err, "创建连接失败")

	//创建Channel
	rabbitMQ.Channel, err = rabbitMQ.Conn.Channel()
	checkErr(err, "创建channel失败")

	_, err = rabbitMQ.Channel.QueueDeclare( // 返回的队列对象内部记录了队列的一些信息，这里没什么用
		rabbitMQ.QueueName,            // 队列名
		config.QueueConfig.AutoDelete, // 是否持久化
		false,                         // 是否自动删除(前提是至少有一个消费者连接到这个队列，之后所有与这个队列连接的消费者都断开时，才会自动删除。注意：生产者客户端创建这个队列，或者没有消费者客户端与这个队列连接时，都不会自动删除这个队列)
		false,                         // 是否为排他队列（排他的队列仅对“首次”声明的conn可见[一个conn中的其他channel也能访问该队列]，conn结束后队列删除）
		false,                         // 是否阻塞
		nil,                           //额外属性（我还不会用）
	)
	if err != nil {
		fmt.Println("声明队列失败", err)
		return nil
	}

	err = rabbitMQ.Channel.ExchangeDeclare(
		rabbitMQ.Exchange,             //交换器名
		config.ExchangeConfig.Kind,    //exchange type：一般用fanout、direct、topic
		config.ExchangeConfig.DurAble, // 是否持久化
		false,                         //是否自动删除（自动删除的前提是至少有一个队列或者交换器与这和交换器绑定，之后所有与这个交换器绑定的队列或者交换器都与此解绑）
		false,                         //设置是否内置的。true表示是内置的交换器，客户端程序无法直接发送消息到这个交换器中，只能通过交换器路由到交换器这种方式
		false,                         // 是否阻塞
		nil,                           // 额外属性
	)
	if err != nil {
		fmt.Println("声明交换器失败", err)
		return nil
	}

	// 3.建立Binding(可随心所欲建立多个绑定关系)
	err = rabbitMQ.Channel.QueueBind(
		rabbitMQ.QueueName,  // 绑定的队列名称
		rabbitMQ.RoutingKey, // bindkey 用于消息路由分发的key
		rabbitMQ.Exchange,   // 绑定的exchange名
		false,               // 是否阻塞
		nil,                 // 额外属性
	)

	if err != nil {
		fmt.Println("绑定队列和交换器失败", err)
		return nil
	}

	return &rabbitMQ

}
func (mq *RabbitMQ) SentMessage(info string) {
	err := mq.Channel.Publish(
		mq.Exchange,
		mq.RoutingKey,
		false, // 是否返回消息(匹配队列)，如果为true, 会根据binding规则匹配queue，如未匹配queue，则把发送的消息返回给发送者
		false, // 是否返回消息(匹配消费者)，如果为true, 消息发送到queue后发现没有绑定消费者，则把发送的消息返回给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(info),
		},
	)
	if err != nil {
		return
	}
}

// ReleaseRes 释放资源,建议NewRabbitMQ获取实例后 配合defer使用
func (mq *RabbitMQ) ReleaseRes() {
	mq.Conn.Close()
	mq.Channel.Close()
}

func checkErr(err error, meg string) {
	if err != nil {
		log.Fatalf("%s:%s\n", meg, err)
	}
}
