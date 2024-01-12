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
	l := conn.NewLister(conn.NewClientFactory(logger))

	for _, res := range l.ListS3(context.TODO(), conn.ListS3Params{}) {
		result = append(result, resultData{
			Type: "S3",
			Name: res.Name,
		})
	}

	for _, res := range l.ListRDS(context.TODO(), conn.ListRDSParams{}) {
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
