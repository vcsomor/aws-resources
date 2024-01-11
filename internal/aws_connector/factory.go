package aws_connector

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/vcsomor/aws-resources/log"
)

func NewClient(ctx context.Context) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Errorf("error")
		return
	}

	fmt.Printf("Config: %v", cfg)
}
