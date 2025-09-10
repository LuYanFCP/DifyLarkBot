package adapter

import (
	"dify_lark_bot/dify"
	"dify_lark_bot/lark"
)

type DifyAdapter struct {
	client *dify.Client
}

func NewDifyAdapter(client *dify.Client) lark.DifyClient {
	return &DifyAdapter{client: client}
}

func (da *DifyAdapter) Chat(message, userID string) (interface{}, error) {
	return da.client.Chat(message, userID)
}