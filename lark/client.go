package lark

import (
	"context"
	"log"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

type LarkClient struct {
	appID     string
	appSecret string
	client    *larkws.Client
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewClient(appID, appSecret string) *LarkClient {
	return &LarkClient{
		appID:     appID,
		appSecret: appSecret,
	}
}

func (c *LarkClient) StartLongPolling(eventHandler *dispatcher.EventDispatcher) error {
	log.Printf("Starting Lark WebSocket connection with app ID: %s", c.appID)

	// Create WebSocket client with event handler
	cli := larkws.NewClient(c.appID, c.appSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
	)

	c.client = cli

	// Create context with cancel for graceful shutdown
	c.ctx, c.cancel = context.WithCancel(context.Background())

	log.Printf("WebSocket connection started for Lark events")

	// Start the client
	return cli.Start(c.ctx)
}

func (c *LarkClient) Stop() error {
	log.Printf("Stopping Lark WebSocket connection...")
	
	if c.cancel != nil {
		c.cancel()
	}
	
	if c.client != nil {
		// 等待连接关闭
		return nil
	}
	
	return nil
}
