package lister

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/threads"
)

type rdsTask struct {
	ctx    context.Context
	logger *logrus.Entry
	client *rds.Client
}

type rdsTaskResult struct {
	results []Result
	error   error
}

var _ threads.Task = (*rdsTask)(nil)

func newRDSTask(
	ctx context.Context,
	logger *logrus.Entry,
	client *rds.Client,
) threads.Task {
	return &rdsTask{
		ctx:    ctx,
		logger: logger,
		client: client,
	}
}

func (t *rdsTask) Execute() any {
	resources, err := conn.NewDefaultRDSOperations(t.logger, t.client).
		ListRDS(t.ctx, conn.ListRDSParams{})
	if err != nil {
		t.logger.WithError(err).
			Error("unable to list resources")
		return rdsTaskResult{error: err}
	}

	var results []Result
	for _, res := range resources {
		results = append(results, Result{
			Arn:          res.Arn,
			ID:           res.ID,
			CreationTime: res.CreateTime,
		})
	}

	t.logger.WithField(logKeyResourceCount, len(resources)).
		Debug("resources listed")

	return rdsTaskResult{results: results}
}
