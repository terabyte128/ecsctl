package cmd

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var getTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List all tasks in a cluster",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := client.ListTasks(context.TODO(), &ecs.ListTasksInput{
			Cluster: &cluster,
		})

		if err != nil {
			log.Fatalf("failed to list services: %v", err)
		}

		taskDescs, err := client.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
			Tasks:   tasks.TaskArns,
			Cluster: &cluster,
		})

		if err != nil {
			log.Fatalf("failed to describe services: %v", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Containers", "Age"})

		for _, task := range taskDescs.Tasks {
			age := time.Since(*task.CreatedAt)
			splat := strings.Split(*task.TaskArn, "/")
			id := splat[len(splat)-1]

			var containers []string

			for _, container := range task.Containers {
				containers = append(containers, *container.Name)
			}

			table.Append([]string{
				id,
				strings.Join(containers, "\n"),
				age.Round(time.Second).String(),
			})
		}

		table.Render()
	},
}

func init() {
	registerClustersArg(getTasksCmd)
}
