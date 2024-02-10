package lister

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
	"github.com/vcsomor/aws-resources/internal/lister/rds_tasks"
	"github.com/vcsomor/aws-resources/internal/lister/s3_tasks"
	"os"
	"slices"
)

type Lister interface {
	List(ctx context.Context) []any
}

type defaultLister struct {
	clientFactory conn.ClientFactory
	executor      executor.SynchronousExecutor
	logger        *logrus.Logger

	regions   []string
	resources []string
}

var _ Lister = (*defaultLister)(nil)

func (l *defaultLister) List(ctx context.Context) []any {
	var res []any

	if slices.Contains(l.resources, "s3") {
		res = append(res, l.listS3(ctx, l.regions)...)
	}
	if slices.Contains(l.resources, "rds") {
		res = append(res, l.listRDS(ctx, l.regions)...)
	}

	return res
}

func (l *defaultLister) listS3(ctx context.Context, regions []string) (results []any) {
	logger := l.logger.WithField(logKeyResourceType, s3ResourceType)

	client, err := l.clientFactory.S3Client(ctx, nil)
	if err != nil {
		logger.WithError(err).
			Error("unable to create the client")
		return
	}

	listResult, err := l.executor.Execute(s3_tasks.NewListTask(ctx, logger, client))
	if err != nil {
		logger.WithError(err).
			Error("unable to start listing task")
		return
	}

	s3Buckets := listResult.(s3_tasks.ListTaskResult)
	if err = s3Buckets.Error; err != nil {
		logger.WithError(err).
			Error("bucket fetch error")
		return
	}

	var getRegionTasks []executor.Task
	for _, b := range s3Buckets.Buckets {
		getRegionTasks = append(getRegionTasks, s3_tasks.NewS3GetRegionTask(ctx, logger, client, b.Name))
	}

	getRegionResults := l.executor.ExecuteAll(getRegionTasks)
	var getRegionOutcomes []s3_tasks.GetRegionResult

	var getTagsTasks []executor.Task
	for _, r := range getRegionResults {
		if err = r.Error; err != nil {
			logger.WithError(err).
				Error("error while fetching the region")
			continue
		}

		getRegionResult := r.Outcome.(s3_tasks.GetRegionResult)
		if err = getRegionResult.Error; err != nil {
			logger.WithError(err).
				Error("region fetch error")
			continue
		}

		if regionFiler(getRegionResult.Region, regions) {
			getTagsClient, errClient := l.clientFactory.S3Client(ctx, &getRegionResult.Region)
			if errClient != nil {
				logger.WithError(errClient).
					Error("client build error")
				continue
			}
			getRegionOutcomes = append(getRegionOutcomes, getRegionResult)
			getTagsTasks = append(getTagsTasks, s3_tasks.NewS3GetTagsTask(ctx, logger, getTagsClient, getRegionResult.BucketName))
		}
	}

	getTagsResult := l.executor.ExecuteAll(getTagsTasks)
	for _, res := range getTagsResult {
		if res.Error != nil {
			logger.WithError(err).
				Error("tags fetch error")
			continue

		}
		gtRes := res.Outcome.(s3_tasks.GetTagsResult)
		bucketName := gtRes.BucketName
		results = append(results, anS3Result(
			s3_tasks.FindListBucketData(bucketName, s3Buckets.Buckets),
			s3_tasks.FindRegionResult(bucketName, getRegionOutcomes),
			gtRes,
		))
	}

	return
}

func (l *defaultLister) listRDS(ctx context.Context, regions []string) (results []any) {
	logger := l.logger.
		WithField(logKeyResourceType, rdsResourceType)

	var tasks []executor.Task
	if regions == nil {
		t, err := l.rdsTaskInRegion(ctx, logger.WithField(logKeyRegion, "default"), nil)
		if err != nil {
			logger.WithError(err).
				Error("unable to add RDS list task in default region")
		} else {
			tasks = append(tasks, t)
		}
	} else {
		for _, region := range regions {
			currRegion := region
			t, err := l.rdsTaskInRegion(ctx, logger.WithField(logKeyRegion, region), &currRegion)
			if err != nil {
				logger.WithError(err).
					Errorf("unable to add RDS list task in region")
				continue
			}
			tasks = append(tasks, t)
		}
	}

	execResults := l.executor.ExecuteAll(tasks)

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
	return
}

func (l *defaultLister) rdsTaskInRegion(ctx context.Context, logger *logrus.Entry, region *string) (executor.Task, error) {
	client, err := l.clientFactory.RDSClient(ctx, region)
	if err != nil {
		logger.WithError(err).
			Error("unable to create the client")
		return nil, err
	}

	return rds_tasks.NewListTask(ctx, logger, client), nil
}

func writeResult(res []any) {
	js, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to display data %s", err)
		return
	}
	fmt.Printf("%s\n", js)
}

func anS3Result(baseData s3_tasks.ListTaskBucketData, region s3_tasks.GetRegionResult, tagsResult s3_tasks.GetTagsResult) Result[S3Data] {
	return Result[S3Data]{
		Arn:          fmt.Sprintf("arn:aws:s3:::%s", baseData.Name),
		ID:           baseData.Name,
		CreationTime: baseData.Created,
		Data: S3Data{
			LocationConstraint: region.Region,
			Tags:               tagsResult.Tags,
		},
	}
}

func regionFiler(region string, regions []string) bool {
	if len(regions) == 0 {
		return true
	}
	for _, r := range regions {
		if region == r {
			return true
		}
	}
	return false
}
