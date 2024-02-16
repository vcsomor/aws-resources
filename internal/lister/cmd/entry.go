package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vcsomor/aws-resources/config"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
	"github.com/vcsomor/aws-resources/internal/lister"
	"github.com/vcsomor/aws-resources/internal/lister/args"
	"github.com/vcsomor/aws-resources/internal/lister/writer"
	"github.com/vcsomor/aws-resources/internal/lister/writer/jsonfile"
	"github.com/vcsomor/aws-resources/log"
	"strconv"
	"strings"
)

// ListResources is the command entry point
func ListResources(command *cobra.Command, _ []string) {
	logger := log.NewLogger(config.Config())

	th := command.Flag("threads").
		Value.
		String()

	threadCount, err := strconv.Atoi(th)
	if err != nil {
		logger.WithError(err).
			Error("invalid thread count")
		return
	}

	userRegions := args.ParseRegions(command.Flag("regions").
		Value.
		String())
	logger.Debugf("regions: %v", userRegions)

	userResources := args.ParseResources(command.Flag("resources").
		Value.
		String())
	logger.Debugf("resources: %v", userResources)

	threadpool, err := executor.NewThreadpool(threadCount)
	if err != nil {
		logger.WithError(err).
			Error("threadpool error")
		return
	}

	defer func() {
		if threadpool != nil {
			threadpool.Shutdown()
		}
	}()

	l := lister.NewLister().
		Dependencies(
			lister.WithClientFactory(conn.NewClientFactory(logger)),
			lister.WithLogger(logger),
			lister.WithExecutor(executor.NewSynchronousExecutor(threadpool)),
			lister.WithWriterFactory(func(r lister.Result) writer.Writer {
				w, _ := jsonfile.NewWriter(fmt.Sprintf("./output/%s", stripArn(r.Arn)))
				return w
			}),
		).
		Parameters(
			lister.WithRegions(userRegions),
			lister.WithResources(userResources),
		).
		Build()
	resources := l.List(context.TODO())

	fmt.Printf("Data fetched %d\n", len(resources))
}

func stripArn(arn string) any {
	return strings.ReplaceAll(arn, ":", "_")
}
