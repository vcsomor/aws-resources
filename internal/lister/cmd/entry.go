package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vcsomor/aws-resources/config"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
	"github.com/vcsomor/aws-resources/internal/lister"
	"github.com/vcsomor/aws-resources/internal/lister/args"
	"github.com/vcsomor/aws-resources/log"
	"os"
	"strconv"
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

	l := lister.NewLister(logger, conn.NewClientFactory(logger), executor.NewSynchronousExecutor(threadpool)).
		WithRegions(userRegions).
		WithResources(userResources).
		Build()
	resources := l.List(context.TODO())

	// TODO vcsomor do the write
	writeResult(resources)
}

func writeResult(res []any) {
	js, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to display data %s", err)
		return
	}
	fmt.Printf("%s\n", js)
}
