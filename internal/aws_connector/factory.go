package aws_connector

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type ClientFactory interface {
	S3Client(ctx context.Context, region *string) (S3Client, error)
	RDSClient(ctx context.Context, region *string) (RDSClient, error)
}

type defaultAwsClientFactory struct {
	logger *logrus.Logger
}

var _ ClientFactory = (*defaultAwsClientFactory)(nil)

func NewClientFactory(logger *logrus.Logger) ClientFactory {
	return &defaultAwsClientFactory{
		logger: logger,
	}
}

func (f *defaultAwsClientFactory) S3Client(ctx context.Context, region *string) (S3Client, error) {
	cfg, err := f.loadConfig(ctx, region)
	log := f.logger.WithField("client", "S3")

	if err != nil {
		log.WithError(err).
			Error("client init failed")
		return nil, err
	}
	log.Debugf("client init successful")
	return newS3Client(s3.NewFromConfig(cfg)), nil
}

func (f *defaultAwsClientFactory) RDSClient(ctx context.Context, region *string) (RDSClient, error) {
	cfg, err := f.loadConfig(ctx, region)
	log := f.logger.WithField("client", "RDS")

	if err != nil {
		log.WithError(err).
			Error("client init failed")
		return nil, err
	}
	log.Debugf("client init successful")
	return newRDSClient(rds.NewFromConfig(cfg)), nil
}

func (f *defaultAwsClientFactory) loadConfig(ctx context.Context, region *string) (aws.Config, error) {
	if region != nil {
		return config.LoadDefaultConfig(ctx, config.WithRegion(*region))
	}
	return config.LoadDefaultConfig(ctx)
}
