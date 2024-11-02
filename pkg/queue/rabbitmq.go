package queue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var (
	connection *amqp.Connection
	channel    *amqp.Channel
)

// InitRabbitMQ 初始化RabbitMQ连接
func InitRabbitMQ() error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		viper.GetString("rabbitmq.username"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
		viper.GetString("rabbitmq.vhost"),
	)

	var err error
	connection, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err = connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %v", err)
	}

	// 声明任务队列
	_, err = channel.QueueDeclare(
		"scheduled_tasks", // 队列名称
		true,              // 持久化
		false,             // 自动删除
		false,             // 排他性
		false,             // 不等待
		nil,               // 额外参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	return nil
}

// CloseRabbitMQ 关闭RabbitMQ连接
func CloseRabbitMQ() {
	if channel != nil {
		channel.Close()
	}
	if connection != nil {
		connection.Close()
	}
}

// PublishTask 发布任务到队列
func PublishTask(task interface{}) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return channel.Publish(
		"",                // 交换机
		"scheduled_tasks", // 队列名称
		false,             // 强制
		false,             // 立即
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// ConsumeTask 消费任务队列
func ConsumeTask(handler func([]byte) error) error {
	msgs, err := channel.Consume(
		"scheduled_tasks", // 队列名称
		"",                // 消费者
		false,             // 自动确认
		false,             // 排他性
		false,             // 不等待
		false,             // 不阻塞
		nil,               // 额外参数
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error processing task: %v", err)
				d.Nack(false, true) // 消息重新入队
			} else {
				d.Ack(false) // 确认消息
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}
