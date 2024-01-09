package lister

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"os"
)

type resultData struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func ListResources(_ *cobra.Command, _ []string) {
	writeResult(fetchData())
}

func fetchData() []resultData {
	var result []resultData
	l := conn.NewLister()
	for _, res := range l.ListS3(conn.ListS3Params{}) {
		result = append(result, resultData{
			Type: "S3",
			Name: res.Name,
		})
	}

	for _, res := range l.ListRDS(conn.ListRDSParams{}) {
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
