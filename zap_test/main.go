package main

import (
	"time"

	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./myproject.log",
		"stderr",
	}
	return cfg.Build()
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	url := "http://test_url"
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
	// logger, err := NewLogger()
	// if err != nil {
	// 	panic(err)
	// }

	// su := logger.Sugar()
	// defer su.Sync()
	// url := "https://imooc.com"
	// su.Info("failed to fetch URL",
	// 	zap.String("url", url),
	// 	zap.Int("attempt", 3),
	// 	zap.Duration("backoff", time.Second),
	// )
}
