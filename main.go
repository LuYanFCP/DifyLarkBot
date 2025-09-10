package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

	"dify_lark_bot/adapter"
	"dify_lark_bot/cmd"
	"dify_lark_bot/config"
	"dify_lark_bot/dify"
	larkservice "dify_lark_bot/lark"
)

const (
	appName    = "DifyLarkBot"
	appVersion = "1.0.0"
)

type Application struct {
	config       config.Config
	larkService  *larkservice.Service
	larkClient   *larkservice.LarkClient
	eventHandler *dispatcher.EventDispatcher
	startTime    time.Time
	shutdownChan chan struct{}
}

func main() {
	cli := cmd.ParseFlags()

	if cli.Help {
		cli.PrintHelp()
		os.Exit(0)
	}

	if cli.Version {
		cli.PrintVersion()
		os.Exit(0)
	}

	app := NewApplication(cli)

	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	app.Start()
}

func NewApplication(cli *cmd.CLI) *Application {
	envConfig := config.Load()

	if cli.ConfigFile != "" {
		fileConfig, err := config.LoadFromFile(cli.ConfigFile)
		if err != nil {
			log.Printf("Warning: Failed to load config file: %v", err)
		} else {
			envConfig = fileConfig.MergeWithEnv(envConfig)
		}
	}

	return &Application{
		config:       envConfig,
		startTime:    time.Now(),
		shutdownChan: make(chan struct{}),
	}
}

func (app *Application) Initialize() error {
	log.Printf("Initializing %s v%s...", appName, appVersion)

	if err := app.validateConfig(); err != nil {
		return fmt.Errorf("config validation: %w", err)
	}

	app.setupServices()
	app.setupEventHandler()

	return nil
}

func (app *Application) validateConfig() error {
	cfg := app.config
	fmt.Printf("%#v", cfg)

	requiredVars := []struct {
		name  string
		value string
	}{
		{"LARK_APP_ID", cfg.LarkAppID},
		{"LARK_APP_SECRET", cfg.LarkAppSecret},
		{"LARK_VERIFICATION_TOKEN", cfg.LarkVerificationToken},
		{"DIFY_API_KEY", cfg.DifyAPIKey},
	}

	for _, rv := range requiredVars {
		if rv.value == "" {
			return fmt.Errorf("%s is required", rv.name)
		}
	}

	return nil
}

func (app *Application) setupEventHandler() {
	// Set up event dispatcher for WebSocket connection
	app.eventHandler = dispatcher.NewEventDispatcher(
		app.config.LarkVerificationToken,
		app.config.LarkEncryptKey,
	)

	app.eventHandler.OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		return app.larkService.HandleMessageEvent(ctx, event)
	})
}

func (app *Application) setupServices() {
	// Initialize Lark client for WebSocket connection
	app.larkClient = larkservice.NewClient(app.config.LarkAppID, app.config.LarkAppSecret)

	// Initialize Dify client and service
	difyClient := dify.NewClient(app.config.DifyAPIKey, app.config.DifyBaseURL)
	difyAdapter := adapter.NewDifyAdapter(difyClient)
	app.larkService = larkservice.NewService(difyAdapter, app.config.LarkAppID, app.config.LarkAppSecret)
}


func (app *Application) Start() {
	log.Printf("%s v%s starting WebSocket connection...", appName, appVersion)
	log.Printf("Dify Base URL: %s", app.config.DifyBaseURL)
	log.Printf("Config file: %s", getConfigFile())

	// 启动WebSocket连接
	go func() {
		if err := app.startWebSocketConnection(); err != nil {
			log.Printf("WebSocket connection failed: %v", err)
		}
	}()

	// 等待关闭信号
	app.waitForShutdown()
}

func (app *Application) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, gracefully shutting down...")
	
	// 触发关闭流程
	app.shutdown()
}

func (app *Application) shutdown() {
	// 关闭WebSocket连接
	if app.larkClient != nil {
		log.Println("Closing WebSocket connection...")
		if err := app.larkClient.Stop(); err != nil {
			log.Printf("Error closing WebSocket connection: %v", err)
		}
	}
	
	// 等待所有异步任务完成
	if app.larkService != nil {
		log.Println("Waiting for async tasks to complete...")
		app.larkService.WaitForCompletion()
		log.Println("All async tasks completed")
	}
	
	// 其他清理工作
	log.Println("Application cleanup completed")
	os.Exit(0)
}

func (app *Application) startWebSocketConnection() error {
	log.Printf("Starting Lark WebSocket connection...")
	
	// Start WebSocket connection using the lark client
	err := app.larkClient.StartLongPolling(app.eventHandler)
	if err != nil {
		log.Printf("WebSocket connection error: %v", err)
		// In a production environment, you might want to implement retry logic
		// or graceful shutdown here
		return err
	}
	return nil
}


func getConfigFile() string {
	for i, arg := range os.Args {
		if (arg == "--config" || arg == "-c") && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
	}
	return ""
}
