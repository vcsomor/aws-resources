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
	S3Client(ctx context.Context) (*s3.Client, error)
	RDSClient(ctx context.Context) (*rds.Client, error)
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

func (f *defaultAwsClientFactory) S3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := f.loadConfig(ctx)
	log := f.logger.WithField("client", "S3")

	if err != nil {
		log.WithError(err).
			Error("client init failed")
		return &s3.Client{}, err
	}
	log.Debugf("client init successful")
	return s3.NewFromConfig(cfg), nil
}

func (f *defaultAwsClientFactory) RDSClient(ctx context.Context) (*rds.Client, error) {
	cfg, err := f.loadConfig(ctx)
	log := f.logger.WithField("client", "RDS")

	if err != nil {
		log.WithError(err).
			Error("client init failed")
		return &rds.Client{}, err
	}
	log.Debugf("client init successful")
	return rds.NewFromConfig(cfg), nil
}

func (f *defaultAwsClientFactory) loadConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}
