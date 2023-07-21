package main

import (
	"context"
	"io/ioutil"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/xyzbit/gid/core/conf"
	"gopkg.in/yaml.v3"
)

func main() {
	ctx := context.Background()
	config := conf.Config{}
	time.Sleep(5 * time.Second)

	// 读取配置文件
	yamlFile, err := ioutil.ReadFile("./configs/configs.yml")
	if err != nil {
		log.Fatalf("Failed to read the YAML file: %v", err)
	}

	// 解析配置文件到结构体
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML into struct: %v", err)
	}

	s, err := initGrpcServer(ctx, config.Server, config.DBConfig, config.ConsulConfig)
	if err != nil {
		log.Fatalf("Failed to init grpc server: %v", err)
	}

	// 创建一个用于优雅关闭的等待组
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		// 监听系统信号
		ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		// 等待接收信号
		<-ctx.Done()
		// 接受第一次信号后关闭监听，下次可强制关闭
		stop()
		log.Println("Received termination signal. Shutting down...")

		// 关闭gRPC服务
		finshCh := make(chan struct{})
		go func() {
			s.Stop()
			finshCh <- struct{}{}
		}()

		// 等待超时或任务完成
		select {
		case <-finshCh:
			log.Println("Server shutdown success.")
		case <-time.After(10 * time.Second):
			log.Println("Server shutdown forcefully by timeout.")
		}

		// 通知等待组任务已完成
		wg.Done()
	}()

	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start grpc server: %v", err)
	}

	wg.Wait()
}
