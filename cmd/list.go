package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vcsomor/aws-resources/internal/lister"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List AWS resources.",
		Long:  "Listing the AWS Resources",
		Run:   lister.ListResources,
	})
}
