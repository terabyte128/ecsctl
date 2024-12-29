package cmd

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a list of resources",
}

func init() {
	getCmd.AddCommand(getServicesCmd)
	getCmd.AddCommand(getClustersCmd)
	getCmd.AddCommand(getTasksCmd)
	getCmd.AddCommand(getTaskDefinitions)
}
