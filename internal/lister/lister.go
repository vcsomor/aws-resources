package lister

import (
	"fmt"
	"github.com/spf13/cobra"
)

func ListResources(_ *cobra.Command, args []string) {
	fmt.Printf("Listing with args %v\n", args)
}
