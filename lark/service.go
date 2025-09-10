package lark

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"dify_lark_bot/dify"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type Service struct {
	difyClient DifyClient
	appID      string
	appSecret  string
	wg         sync.WaitGroup
}

type DifyClient interface {
	Chat(message, userID string) (any, error)
}

type ChatResponse struct {
	MessageID      string
	ConversationID string
	Answer         string
}

func NewService(difyClient DifyClient, appID, appSecret string) *Service {
	return &Service{
		difyClient: difyClient,
		appID:      appID,
		appSecret:  appSecret,
	}
}

func (s *Service) HandleMessageEvent(ctx context.Context, event any) error {
	msgEvent, ok := event.(*larkim.P2MessageReceiveV1)
	if !ok {
		return nil
	}

	return s.handleLarkEvent(ctx, msgEvent)
}

func (s *Service) handleLarkEvent(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	if event.Event.Message == nil || event.Event.Message.Content == nil {
		return nil
	}

	var content map[string]any
	if err := json.Unmarshal([]byte(*event.Event.Message.Content), &content); err != nil {
		return fmt.Errorf("parse content: %w", err)
	}

	text, ok := content["text"].(string)
	if !ok || text == "" {
		return nil
	}

	if !s.isBotMentioned(text) {
		return nil
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.processMessageAsync(ctx, event, text)
	}()

	return nil
}

func (s *Service) isBotMentioned(text string) bool {
	return strings.Contains(text, "@") || strings.Contains(text, "机器人")
}

func (s *Service) buildInteractiveMessage(content, senderName string) map[string]any {
	// 构建interactive富文本消息，包含@提及
	return map[string]any{
		"elements": []map[string]any{
			{
				"tag": "div",
				"text": map[string]any{
					"content": fmt.Sprintf("<at id=\"%s\"></at>%s", senderName, content),
					"tag":     "lark_md",
				},
			},
		},
	}
}

func (s *Service) processMessageAsync(ctx context.Context, event *larkim.P2MessageReceiveV1, text string) {
	userMessage := s.extractUserMessage(text)
	userID := s.getUserID(event)

	senderName := s.getSenderName(event)

	fmt.Printf("Processing message from user %s: %s\n", senderName, userMessage)

	response, err := s.difyClient.Chat(userMessage, userID)
	if err != nil {
		fmt.Printf("Dify chat error: %v\n", err)
		return
	}

	if resp, ok := response.(*dify.ChatResponse); ok {
		interactiveMsg := s.buildInteractiveMessage(resp.Answer, senderName)
		if err := s.replyInteractiveMessage(ctx, event, interactiveMsg); err != nil {
			fmt.Printf("Failed to send reply: %v\n", err)
		} else {
			fmt.Printf("Successfully replied to %s: %s\n", senderName, resp.Answer)
		}
	} else {
		fmt.Printf("Invalid response type from Dify\n")
	}
}

func (s *Service) getSenderName(event *larkim.P2MessageReceiveV1) string {
	if event.Event.Sender != nil {
		if event.Event.Sender.SenderId != nil {
			if openID := event.Event.Sender.SenderId.OpenId; openID != nil {
				return *openID
			}
		}
	}
	return "用户"
}

func (s *Service) extractUserMessage(text string) string {
	return strings.TrimSpace(text)
}

func (s *Service) getUserID(event *larkim.P2MessageReceiveV1) string {
	if event.Event.Sender != nil && event.Event.Sender.SenderId != nil {
		if openID := event.Event.Sender.SenderId.OpenId; openID != nil {
			return *openID
		}
	}
	return "anonymous"
}

func (s *Service) replyMessage(ctx context.Context, event *larkim.P2MessageReceiveV1, replyText string) error {
	// Create Lark client
	client := lark.NewClient(s.appID, s.appSecret)

	// Get message information
	messageID := event.Event.Message.MessageId
	receiveID := event.Event.Message.ChatId

	if messageID == nil || receiveID == nil {
		return fmt.Errorf("missing message ID or chat ID")
	}

	// Create reply content with proper formatting
	replyContent := map[string]any{
		"text": replyText,
	}

	contentBytes, err := json.Marshal(replyContent)
	if err != nil {
		return fmt.Errorf("failed to marshal reply content: %w", err)
	}

	// Create reply request using the message reply API
	req := larkim.NewReplyMessageReqBuilder().
		MessageId(*messageID).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			Content(string(contentBytes)).
			MsgType("text").
			Build()).
		Build()

	// Send the reply
	resp, err := client.Im.Message.Reply(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send reply message: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("failed to send reply message: %s", resp.Msg)
	}

	fmt.Printf("Successfully sent reply: %s\n", replyText)
	return nil
}

func (s *Service) replyInteractiveMessage(ctx context.Context, event *larkim.P2MessageReceiveV1, interactiveContent map[string]any) error {
	// Create Lark client
	client := lark.NewClient(s.appID, s.appSecret)

	// Get message information
	messageID := event.Event.Message.MessageId
	receiveID := event.Event.Message.ChatId

	if messageID == nil || receiveID == nil {
		return fmt.Errorf("missing message ID or chat ID")
	}

	// Convert interactive content to JSON
	contentBytes, err := json.Marshal(interactiveContent)
	if err != nil {
		return fmt.Errorf("failed to marshal interactive content: %w", err)
	}

	// Create interactive message request
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(*receiveID).
			Content(string(contentBytes)).
			Build()).
		Build()

	// Send the interactive message
	resp, err := client.Im.Message.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send interactive message: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("failed to send interactive message: %s", resp.Msg)
	}

	fmt.Printf("Successfully sent interactive message\n")
	return nil
}

func (s *Service) WaitForCompletion() {
	done := make(chan struct{})

	go func() {
		s.wg.Wait()
		close(done)
	}()

	// 设置30秒超时
	select {
	case <-done:
		log.Println("All async messages processed")
	case <-time.After(30 * time.Second):
		log.Println("Timeout waiting for async tasks, forcing shutdown")
	}
}
