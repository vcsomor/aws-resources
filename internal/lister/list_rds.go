package lister

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vcsomor/aws-resources/internal/executor"
	"github.com/vcsomor/aws-resources/internal/lister/rds_tasks"
)

func (l *taskBasedLister) listRDS(ctx context.Context) (results []any) {
	logger := l.logger.
		WithField(logKeyResourceType, rdsResourceType)

	return assembleRDBTasksResults(
		l.executor.ExecuteAll(
			l.makeRdsListTasks(ctx, logger)), logger)

}

func (l *taskBasedLister) makeRdsListTasks(ctx context.Context, logger *logrus.Entry) []executor.Task {
	if l.regions == nil {
		t, err := l.rdsTaskInRegion(ctx, logger.WithField(logKeyRegion, "default"), nil)
		if err != nil {
			logger.WithError(err).
				Error("unable to add RDS list task in default region")
			return nil
		}
		return []executor.Task{t}
	}

	var tasks []executor.Task
	for _, region := range l.regions {
		currRegion := region
		t, err := l.rdsTaskInRegion(ctx, logger.WithField(logKeyRegion, region), &currRegion)
		if err != nil {
			logger.WithError(err).
				Errorf("unable to add RDS list task in region")
			continue
		}
		tasks = append(tasks, t)
	}

	return tasks
}

func assembleRDBTasksResults(execResults []executor.SynchronousResult, logger *logrus.Entry) []any {
	var results []any
	for _, r := range execResults {
		if err := r.Error; err != nil {
			logger.WithError(err).
				Error("error while fetching the RDS instances")
			continue
		}

		listResult := r.Outcome.(rds_tasks.ListResult)
		if err := listResult.Error; err != nil {
			logger.WithError(err).
				Error("rds list task error")
			continue
		}

		for _, rds := range listResult.RDSInstances {
			results = append(results, Result[RDSData]{
				Arn:          rds.Arn,
				ID:           rds.ID,
				CreationTime: rds.CreationTime,
				Data: RDSData{
					Tags: rds.Tags,
				},
			})
		}
	}
	return results
}

func (l *taskBasedLister) rdsTaskInRegion(ctx context.Context, logger *logrus.Entry, region *string) (executor.Task, error) {
	client, err := l.clientFactory.RDSClient(ctx, region)
	if err != nil {
		logger.WithError(err).
			Error("unable to create the client")
		return nil, err
	}

	return rds_tasks.NewListTask(ctx, logger, client), nil
}
