package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var region string

var rootCmd = &cobra.Command{
	Use: "ecsctl",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		makeClient(region)
	},
}

var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "AWS region")
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(execCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
