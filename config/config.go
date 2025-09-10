package config

import (
	"fmt"
	"os"
)

type Config struct {
	LarkAppID             string
	LarkAppSecret         string
	LarkVerificationToken string
	LarkEncryptKey        string
	DifyAPIKey            string
	DifyBaseURL           string
	Port                  string
}

func Load() Config {
	ret := Config{
		LarkAppID:             getEnv("LARK_APP_ID", ""),
		LarkAppSecret:         getEnv("LARK_APP_SECRET", ""),
		LarkVerificationToken: getEnv("LARK_VERIFICATION_TOKEN", ""),
		LarkEncryptKey:        getEnv("LARK_ENCRYPT_KEY", ""),
		DifyAPIKey:            getEnv("DIFY_API_KEY", ""),
		DifyBaseURL:           getEnv("DIFY_BASE_URL", "https://api.dify.ai"),
		Port:                  getEnv("PORT", "8080"),
	}
	fmt.Printf("%#v", ret)
	return ret
}

func getEnv(key, defaultValue string) string {

	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
