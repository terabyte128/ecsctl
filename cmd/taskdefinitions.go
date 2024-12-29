package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var EXCLUDED_TASK_DEFINITION_KEYS = []string{
	"TaskDefinitionArn",
	"Family",
	"Revision",
	"Status",
	"RequiresAttributes",
	"Compatibilities",
	"RegisteredAt",
	"RegisteredBy",
	"Tags",
}

var editTaskDefinition = &cobra.Command{
	Use:   "taskdefinitions",
	Short: "Edit a task definition",
	Run: func(cmd *cobra.Command, args []string) {
		defn := "scorecard-development"
		definition, err := client.DescribeTaskDefinition(context.TODO(), &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: &defn, // TODO
		})
		exitWithError("describe task definition", err)

		out, err := json.MarshalIndent(definition.TaskDefinition, "", "  ")

		fmt.Println(string(out))

		// vi := exec.Command("nvim")
		// vi.Stdin = os.Stdin
		// vi.Stdout = os.Stdout
		// vi.Stderr = os.Stderr
		//
		// vi.Run()
	},
}

var getTaskDefinitions = &cobra.Command{
	Use:   "taskdefinitions",
	Short: "Get a list of task definitions",
	Run: func(cmd *cobra.Command, args []string) {
		families, err := client.ListTaskDefinitionFamilies(
			context.TODO(),
			&ecs.ListTaskDefinitionFamiliesInput{},
		)
		exitWithError("list task defintion families", err)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Family", "Revision"})

		for _, family := range families.Families {
			taskDefinition, err := client.DescribeTaskDefinition(context.TODO(), &ecs.DescribeTaskDefinitionInput{
				TaskDefinition: &family,
			})
			revision := "n/a"

			if err == nil {
				revision = fmt.Sprintf("%d", taskDefinition.TaskDefinition.Revision)
			} else if !*includeInactive {
				continue
			}

			table.Append([]string{
				family,
				revision,
			})
		}

		table.Render()
	},
}

var includeInactive *bool

func init() {
	includeInactive = getTaskDefinitions.Flags().Bool("include-inactive", false, "whether to include inactive definitions")
}
