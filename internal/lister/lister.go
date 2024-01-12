package lister

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vcsomor/aws-resources/config"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/log"
	"os"
)

type resultData struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func ListResources(*cobra.Command, []string) {
	writeResult(fetchData())
}

func fetchData() []resultData {
	logger := log.NewLogger(config.Config())

	var result []resultData
	l := conn.NewLister(logger, conn.NewClientFactory(logger))

	s3Items, err := l.ListS3(context.TODO(), conn.ListS3Params{})
	if err != nil {
		logger.WithError(err).
			Error("unable to list S3 resources")
		return nil
	}
	for _, res := range s3Items {
		result = append(result, resultData{
			Type: "S3",
			Name: res.Name,
		})
	}

	rdsItems, err := l.ListRDS(context.TODO(), conn.ListRDSParams{})
	if err != nil {
		logger.WithError(err).
			Error("unable to list RDS resources")
		return nil
	}
	for _, res := range rdsItems {
		result = append(result, resultData{
			Type: "RDS",
			Name: res.Name,
		})
	}
	return result
}

func writeResult(res []resultData) {
	js, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to display data %s", err)
		return
	}
	fmt.Printf("%s\n", js)
}
