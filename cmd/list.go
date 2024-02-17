package cmd

import (
	"github.com/spf13/cobra"
	listcmd "github.com/vcsomor/aws-resources/internal/lister/cmd"
)

func listCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "List AWS resources.",
		Long:  "Listing the AWS Resources",
		Run:   listcmd.ListResources,
	}

	cmd.PersistentFlags().
		String(
			"threads",
			"2",
			`Specify thread count for querying the resources.`)

	cmd.PersistentFlags().
		String(
			"regions",
			"all",
			`Specify regions to list for e.g.: --regions us-east-1,us-east-2. Use "all" for every supported region.`)

	cmd.PersistentFlags().
		String(
			"resources",
			"all",
			`Specify resources to list for e.g.: --resources s3,rds. Use "all" for every supported resource type.`)

	cmd.PersistentFlags().
		String(
			"output",
			"stdout",
			`Specify the output  e.g.: --output stdout,file. Possible values are: "stdout", "file"`)

	cmd.PersistentFlags().
		String(
			"target",
			"resources",
			`Specify the target directory if file output has been specified.`)

	return &cmd
}
