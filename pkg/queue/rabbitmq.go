package queue

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
	mu      sync.Mutex
)

// InitRabbitMQ 初始化RabbitMQ连接
func InitRabbitMQ() error {
	mu.Lock()
	defer mu.Unlock()

	if conn != nil {
		return nil
	}

	// 从配置中获取RabbitMQ连接信息

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		viper.GetString("rabbitmq.username"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
		viper.GetString("rabbitmq.vhost"),
	)

	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %v", err)
	}

	// 声明任务队列
	_, err = channel.QueueDeclare(
		"task_queue", // 队列名称
		true,         // 持久化
		false,        // 自动删除
		false,        // 独占
		false,        // 不等待
		nil,          // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	return nil
}

// PublishMessage 发布消息到队列
func PublishMessage(queueName string, body []byte) error {
	mu.Lock()
	defer mu.Unlock()

	if channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	err := channel.Publish(
		"",        // 交换机
		queueName, // 路由键
		false,     // 强制
		false,     // 立即
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

// ConsumeMessages 消费队列消息
func ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	mu.Lock()
	defer mu.Unlock()

	if channel == nil {
		return nil, fmt.Errorf("RabbitMQ channel not initialized")
	}

	// 设置QoS
	err := channel.Qos(
		1,     // 预取计数
		0,     // 预取大小
		false, // 全局
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set QoS: %v", err)
	}

	msgs, err := channel.Consume(
		queueName, // 队列
		"",        // 消费者
		false,     // 自动确认
		false,     // 独占
		false,     // 不等待
		false,     // 参数
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %v", err)
	}

	return msgs, nil
}

// CloseConnection 关闭RabbitMQ连接
func CloseConnection() error {
	mu.Lock()
	defer mu.Unlock()

	var err error

	if channel != nil {
		err = channel.Close()
		if err != nil {
			return fmt.Errorf("failed to close channel: %v", err)
		}
		channel = nil
	}

	if conn != nil {
		err = conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close connection: %v", err)
		}
		conn = nil
	}

	return nil
}
