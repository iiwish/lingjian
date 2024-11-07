package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iiwish/lingjian/internal/model"
	"github.com/iiwish/lingjian/internal/service"
	"github.com/iiwish/lingjian/pkg/queue"
	"github.com/spf13/viper"
)

func init() {
	// 加载配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 初始化数据库连接
	if err := model.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化RabbitMQ连接
	if err := queue.InitRabbitMQ(); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
}

func main() {
	taskService := &service.TaskService{}

	// 创建任务队列通道
	taskQueue := make(chan model.TaskMessage)
	defer close(taskQueue)

	// 启动消息消费者
	go func() {
		msgs, err := queue.ConsumeMessages("task_queue")
		if err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}

		for msg := range msgs {
			var taskMsg model.TaskMessage
			if err := json.Unmarshal(msg.Body, &taskMsg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false)
				continue
			}

			// 发送任务到处理通道
			taskQueue <- taskMsg
			msg.Ack(false)
		}
	}()

	// 启动任务处理器
	go func() {
		for task := range taskQueue {
			log.Printf("Processing task %d", task.TaskID)
			if err := taskService.ExecuteTask(task.TaskID); err != nil {
				log.Printf("Failed to execute task %d: %v", task.TaskID, err)
			}
		}
	}()

	// 优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down worker...")

	// 关闭RabbitMQ连接
	if err := queue.CloseConnection(); err != nil {
		log.Printf("Error closing RabbitMQ connection: %v", err)
	}

	log.Println("Worker shutdown complete")
}
