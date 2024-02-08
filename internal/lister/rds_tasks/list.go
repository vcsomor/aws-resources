package rds_tasks

import (
	"context"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
	"time"
)

type listTask struct {
	ctx    context.Context
	logger *logrus.Entry
	client conn.RDSClient
}

type ListResultRDSData struct {
	Arn          string
	ID           string
	CreationTime *time.Time
	Tags         map[string]*string
}

type ListResult struct {
	RDSInstances []ListResultRDSData
	Error        error
}

var _ executor.Task = (*listTask)(nil)

func NewListTask(
	ctx context.Context,
	logger *logrus.Entry,
	client conn.RDSClient,
) executor.Task {
	return &listTask{
		ctx:    ctx,
		logger: logger,
		client: client,
	}
}

func (t *listTask) Execute() any {
	resources, err := t.client.List(t.ctx, conn.ListRDSParams{})
	if err != nil {
		t.logger.WithError(err).
			Error("unable to list resources")
		return ListResult{Error: err}
	}

	var rdsInstances []ListResultRDSData
	for _, res := range resources {
		rdsInstances = append(rdsInstances, ListResultRDSData{
			Arn:          res.Arn,
			ID:           res.ID,
			CreationTime: res.CreateTime,
			Tags:         res.Tags,
		})
	}

	t.logger.Debug("resources listed")

	return ListResult{RDSInstances: rdsInstances}
}
