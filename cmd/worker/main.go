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
	// 处理系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down worker...")
		queue.CloseRabbitMQ()
		os.Exit(0)
	}()

	// 启动任务消费者
	err := queue.ConsumeTask(func(body []byte) error {
		var message service.TaskMessage
		if err := json.Unmarshal(body, &message); err != nil {
			return err
		}

		// TODO: 根据任务类型执行不同的处理逻辑
		// 例如执行SQL语句或发送HTTP请求

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to start task consumer: %v", err)
	}
}
