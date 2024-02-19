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
	"github.com/vcsomor/aws-resources/internal/lister/writer/jsonfile"
	"github.com/vcsomor/aws-resources/internal/lister/writer/stdout"
	"github.com/vcsomor/aws-resources/log"
	"slices"
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

	argRegions := args.ParseRegions(command.Flag("regions").
		Value.
		String())
	logger.Debugf("regions: %v", argRegions)

	argResources := args.ParseResources(command.Flag("resources").
		Value.
		String())
	logger.Debugf("resources: %v", argResources)

	argOutputs := args.ParseOutputs(command.Flag("output").
		Value.
		String())
	logger.Debugf("output: %v", argOutputs)

	argTarget := command.Flag("target").
		Value.
		String()
	logger.Debugf("target: %v", argTarget)

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

	deps := []lister.DependencyFn{
		lister.WithClientFactory(conn.NewClientFactory(logger)),
		lister.WithLogger(logger),
		lister.WithExecutor(executor.NewSynchronousExecutor(threadpool)),
	}

	res := lister.NewLister().
		Dependencies(deps...).
		Parameters(
			lister.WithRegions(argRegions),
			lister.WithResources(argResources),
		).
		Build().
		List(context.TODO())

	if slices.Contains(argOutputs, args.OutputFile) {
		writeOutputFiles(argTarget, res)
	}

	if slices.Contains(argOutputs, args.OutputStdout) {
		writeStandardOut(res)
	}
}

func writeOutputFiles(toFolder string, res []lister.Result) {
	for _, result := range res {
		w, err := jsonfile.NewWriter(
			toFolder,
			jsonfile.WithOutputFile(fmt.Sprintf("%s.json", stripArn(result.Arn))),
			jsonfile.WithIndentation("\t"))
		if err != nil {
			continue
		}
		_ = w.Write(result)
	}
}

func writeStandardOut(res []lister.Result) {
	w, err := stdout.NewWriter(stdout.WithIndentation("\t"))
	if err != nil {
		return
	}
	_ = w.Write(res)
}

func stripArn(arn string) any {
	return strings.ReplaceAll(arn, ":", "_")
}
