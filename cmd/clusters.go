package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func completeClusters(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	clusters, err := client.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		log.Fatalf("failed to list clusters: %v", err)
	}

	clusterNames := lo.Map(clusters.ClusterArns, func(val string, _ int) string {
		splat := strings.Split(val, "/")
		return splat[len(splat)-1]
	})

	return clusterNames, cobra.ShellCompDirectiveNoFileComp
}

var cluster string

func registerClustersArg(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&cluster, "cluster", "c", "", "ECS cluster")
	cmd.RegisterFlagCompletionFunc("cluster", completeClusters)
	cmd.MarkFlagRequired("cluster")
}

var getClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "List all clusters in the current region",
	Run: func(cmd *cobra.Command, args []string) {
		clusters, err := client.ListClusters(context.TODO(), &ecs.ListClustersInput{})
		if err != nil {
			log.Fatalf("failed to list clusters: %v", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Status", "Tasks", "Services"})

		clusterDescs, err := client.DescribeClusters(
			context.TODO(),
			&ecs.DescribeClustersInput{
				Clusters: clusters.ClusterArns,
			},
		)

		for _, cluster := range clusterDescs.Clusters {
			tasks := fmt.Sprintf("%d/%d", cluster.RunningTasksCount, cluster.PendingTasksCount+cluster.RunningTasksCount)

			table.Append([]string{
				*cluster.ClusterName,
				*cluster.Status,
				tasks,
				fmt.Sprintf("%d", cluster.ActiveServicesCount),
			})
		}

		table.Render()
	},
}
