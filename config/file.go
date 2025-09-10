package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileConfig struct {
	Lark     LarkConfig     `json:"lark"`
	Dify     DifyConfig     `json:"dify"`
	Logging  LoggingConfig  `json:"logging"`
}

type LarkConfig struct {
	AppID             string `json:"app_id"`
	AppSecret         string `json:"app_secret"`
	VerificationToken string `json:"verification_token"`
	EncryptKey        string `json:"encrypt_key"`
}

type DifyConfig struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
	Timeout int    `json:"timeout"`
}

type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

func LoadFromFile(filePath string) (*FileConfig, error) {
	if filePath == "" {
		return nil, nil
	}
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	
	var config FileConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}
	
	return &config, nil
}

func (fc *FileConfig) MergeWithEnv(envConfig Config) Config {
	if fc == nil {
		return envConfig
	}
	
	result := envConfig
	
	if fc.Lark.AppID != "" {
		result.LarkAppID = fc.Lark.AppID
	}
	if fc.Lark.AppSecret != "" {
		result.LarkAppSecret = fc.Lark.AppSecret
	}
	if fc.Lark.VerificationToken != "" {
		result.LarkVerificationToken = fc.Lark.VerificationToken
	}
	if fc.Lark.EncryptKey != "" {
		result.LarkEncryptKey = fc.Lark.EncryptKey
	}
	
	if fc.Dify.APIKey != "" {
		result.DifyAPIKey = fc.Dify.APIKey
	}
	if fc.Dify.BaseURL != "" {
		result.DifyBaseURL = fc.Dify.BaseURL
	}
	
	return result
}