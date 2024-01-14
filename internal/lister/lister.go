package lister

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vcsomor/aws-resources/config"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/log"
	"os"
	"time"
)

const (
	logKeyResourceType  = "resource-type"
	logKeyResourceCount = "resource-count"

	s3ResourceType  = "S3"
	rdsResourceType = "RDS"
)

type Result struct {
	Arn          string     `json:"arn"`
	ID           string     `json:"id"`
	CreationTime *time.Time `json:"creationTime"`
}

type Lister interface {
	List(ctx context.Context) []Result
}

type defaultLister struct {
	clientFactory conn.ClientFactory
	logger        *logrus.Logger
}

var _ Lister = (*defaultLister)(nil)

// CmdListResources is the command entry point
func CmdListResources(*cobra.Command, []string) {
	logger := log.NewLogger(config.Config())

	l := NewDefaultLister(logger, conn.NewClientFactory(logger))
	resources := l.List(context.TODO())
	logger.WithField(logKeyResourceCount, len(resources)).
		Debug("resources listed")

	// TODO vcsomor do the write
	writeResult(resources)
}

func NewDefaultLister(logger *logrus.Logger, clientFactory conn.ClientFactory) Lister {
	return &defaultLister{
		clientFactory: clientFactory,
		logger:        logger,
	}
}

func (l *defaultLister) List(ctx context.Context) (result []Result) {
	result = append(result, l.listS3(ctx)...)
	result = append(result, l.listRDS(ctx)...)
	return
}

func (l *defaultLister) listS3(ctx context.Context) (result []Result) {
	logger := l.logger.WithField(logKeyResourceType, s3ResourceType)

	client, err := l.clientFactory.S3Client(ctx)
	if err != nil {
		logger.WithError(err).
			Error("unable to create the client")
		return
	}

	resources, err := conn.NewDefaultS3Operations(l.logger, client).
		ListS3(ctx, conn.ListS3Params{})
	if err != nil {
		logger.WithError(err).
			Error("unable to list resources")
		return
	}

	for _, res := range resources {
		result = append(result, Result{
			Arn:          res.Arn,
			ID:           res.Name,
			CreationTime: res.CreationTime,
		})
	}

	logger.WithField(logKeyResourceCount, len(resources)).
		Debugf("resouces listed")
	return
}

func (l *defaultLister) listRDS(ctx context.Context) (result []Result) {
	logger := l.logger.WithField(logKeyResourceType, rdsResourceType)

	client, err := l.clientFactory.RDSClient(ctx)
	if err != nil {
		logger.WithError(err).
			Error("unable to create the client")
		return
	}

	resources, err := conn.NewDefaultRDSOperations(l.logger, client).
		ListRDS(ctx, conn.ListRDSParams{})
	if err != nil {
		logger.WithError(err).
			Error("unable to list resources")
		return
	}

	for _, res := range resources {
		result = append(result, Result{
			Arn:          res.Arn,
			ID:           res.ID,
			CreationTime: res.CreateTime,
		})
	}

	logger.WithField(logKeyResourceCount, len(resources)).
		Debug("resources listed")
	return
}

func writeResult(res []Result) {
	js, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to display data %s", err)
		return
	}
	fmt.Printf("%s\n", js)
}
