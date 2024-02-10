package lister

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/vcsomor/aws-resources/config"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
	"github.com/vcsomor/aws-resources/internal/lister/args"
	"github.com/vcsomor/aws-resources/log"
	"strconv"
)

// CmdListResources is the command entry point
func CmdListResources(command *cobra.Command, _ []string) {
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

	l := newLister(logger, conn.NewClientFactory(logger), executor.NewSynchronousExecutor(threadpool)).
		withRegions(userRegions).
		withResources(userResources).
		build()
	resources := l.List(context.TODO())

	logger.WithField(logKeyResourceCount, len(resources)).
		Debug("resources listed")

	// TODO vcsomor do the write
	writeResult(resources)
}
