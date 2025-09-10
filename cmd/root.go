package cmd

import (
	"flag"
	"fmt"
)

type CLI struct {
	ConfigFile string
	Version    bool
	Help       bool
}

func ParseFlags() *CLI {
	cli := &CLI{}
	
	flag.StringVar(&cli.ConfigFile, "config", "", "Configuration file path")
	flag.StringVar(&cli.ConfigFile, "c", "", "Configuration file path (short)")
	flag.BoolVar(&cli.Version, "version", false, "Show version information")
	flag.BoolVar(&cli.Version, "v", false, "Show version information (short)")
	flag.BoolVar(&cli.Help, "help", false, "Show help information")
	flag.BoolVar(&cli.Help, "h", false, "Show help information (short)")
	
	flag.Parse()
	
	return cli
}

func (cli *CLI) PrintVersion() {
	fmt.Printf("DifyLarkBot v%s\n", appVersion)
	fmt.Printf("Go Version: %s\n", "1.21+")
	fmt.Printf("Build Time: %s\n", buildTime)
	fmt.Printf("Git Commit: %s\n", gitCommit)
}

func (cli *CLI) PrintHelp() {
	fmt.Printf(`DifyLarkBot - Lark Bot with Dify AI Integration

Usage:
  dify_lark_bot [options]

Options:
  -c, --config string    Configuration file path
  -v, --version         Show version information
  -h, --help           Show help information

Environment Variables:
  LARK_APP_ID              Lark application ID
  LARK_APP_SECRET          Lark application secret
  LARK_VERIFICATION_TOKEN  Lark verification token
  LARK_ENCRYPT_KEY         Lark encrypt key (optional)
  DIFY_API_KEY             Dify API key
  DIFY_BASE_URL            Dify base URL (default: https://api.dify.ai)

Examples:
  # Run with environment variables
  export LARK_APP_ID=your_app_id
  export LARK_APP_SECRET=your_app_secret
  export DIFY_API_KEY=your_dify_key
  ./dify_lark_bot

  # Run with config file
  ./dify_lark_bot --config config.json

  # Show version
  ./dify_lark_bot --version
`)
}

var (
	appVersion = "1.0.0"
	buildTime  = "unknown"
	gitCommit  = "unknown"
)