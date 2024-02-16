package lister

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vcsomor/aws-resources/internal/executor"
	"github.com/vcsomor/aws-resources/internal/lister/s3_tasks"
	"slices"
)

func (l *taskBasedLister) listS3(ctx context.Context) []Result {
	logger := l.logger.WithField(logKeyResourceType, s3ResourceType)

	buckets, err := l.fetchAllS3Buckets(ctx, logger)
	if err != nil {
		logger.WithError(err).
			Error("unable fetch all buckets")
		return nil
	}

	regionMappings, err := l.fetchAllS3BucketRegions(ctx, logger, buckets.Buckets)
	if err != nil {
		logger.WithError(err).
			Error("unable fetch regions for buckets")
		return nil
	}

	tagMappings := l.fetchTagsForBuckets(ctx, logger, regionMappings)

	return assembleResults(buckets.Buckets, regionMappings, tagMappings)
}

func assembleResults(
	buckets []s3_tasks.ListTaskBucketData,
	regionMappings map[string]string,
	tags map[string]map[string]*string,
) []Result {
	var result []Result

	for _, bucket := range buckets {
		bucketName := bucket.Name
		if r, exist := regionMappings[bucketName]; exist {
			result = append(result, anS3Result(bucket, r, tags[bucketName]))
		}
	}

	return result
}

func (l *taskBasedLister) fetchAllS3Buckets(ctx context.Context, logger *logrus.Entry) (s3_tasks.ListTaskResult, error) {
	client, err := l.clientFactory.S3Client(ctx, nil)
	if err != nil {
		return s3_tasks.ListTaskResult{},
			errors.Join(errors.New("unable to create the client"), err)
	}

	taskResult, err := l.executor.Execute(s3_tasks.NewListTask(ctx, logger, client))
	if err != nil {
		return s3_tasks.ListTaskResult{},
			errors.Join(errors.New("unable to start listing task"), err)
	}

	buckets := taskResult.(s3_tasks.ListTaskResult)
	if err = buckets.Error; err != nil {
		return s3_tasks.ListTaskResult{},
			errors.Join(errors.New("task error"), err)
	}
	return buckets, nil
}

func (l *taskBasedLister) fetchAllS3BucketRegions(
	ctx context.Context,
	logger *logrus.Entry,
	buckets []s3_tasks.ListTaskBucketData,
) (map[string]string, error) {
	client, err := l.clientFactory.S3Client(ctx, nil)
	if err != nil {
		return nil,
			errors.Join(errors.New("unable to create the client"), err)
	}

	var tasks []executor.Task
	for _, b := range buckets {
		tasks = append(tasks, s3_tasks.NewS3GetRegionTask(ctx, logger, client, b.Name))
	}

	result := map[string]string{}
	for _, execResult := range l.executor.ExecuteAll(tasks) {
		if err = execResult.Error; err != nil {
			logger.WithError(err).
				Error("error while fetching the region")
			continue
		}

		taskResult := execResult.Outcome.(s3_tasks.GetRegionResult)
		if err = taskResult.Error; err != nil {
			logger.WithError(err).
				Error("region fetch error")
			continue
		}

		if region := taskResult.Region; regionFiler(region, l.regions) {
			result[taskResult.BucketName] = region
		}
	}
	return result, nil
}

func (l *taskBasedLister) fetchTagsForBuckets(
	ctx context.Context,
	logger *logrus.Entry,
	mappings map[string]string,
) map[string]map[string]*string {
	tags := map[string]map[string]*string{}

	var tasks []executor.Task
	for name, region := range mappings {
		r := region // avoid taking the address of the auto var
		client, errClient := l.clientFactory.S3Client(ctx, &r)
		if errClient != nil {
			logger.WithError(errClient).
				Error("client build error")
			continue
		}

		tasks = append(tasks, s3_tasks.NewS3GetTagsTask(ctx, logger, client, name))
	}

	for _, execResult := range l.executor.ExecuteAll(tasks) {
		if err := execResult.Error; err != nil {
			logger.WithError(err).
				Error("error while fetching the bucket tags")
			continue
		}

		taskResult := execResult.Outcome.(s3_tasks.GetTagsResult)
		if err := taskResult.Error; err != nil {
			logger.WithError(err).
				Error("tags fetch error")
			continue
		}

		tags[taskResult.BucketName] = taskResult.Tags
	}
	return tags
}

func anS3Result(baseData s3_tasks.ListTaskBucketData, region string, tags map[string]*string) Result {
	return Result{
		Arn:          fmt.Sprintf("arn:aws:s3:::%s", baseData.Name),
		ID:           baseData.Name,
		CreationTime: baseData.Created,
		Data: S3Data{
			LocationConstraint: region,
			Tags:               tags,
		},
	}
}

func regionFiler(region string, regions []string) bool {
	if len(region) == 0 {
		return false
	}
	return slices.Contains(regions, region)
}
