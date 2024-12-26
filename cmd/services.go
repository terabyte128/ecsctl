package cmd

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var getServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List all services in a cluster",
	Run: func(cmd *cobra.Command, args []string) {
		services, err := client.ListServices(context.TODO(), &ecs.ListServicesInput{
			Cluster: &cluster,
		})

		if err != nil {
			log.Fatalf("failed to list services: %v", err)
		}

		serviceDescs, err := client.DescribeServices(context.TODO(), &ecs.DescribeServicesInput{
			Services: services.ServiceArns,
			Cluster:  &cluster,
		})

		if err != nil {
			log.Fatalf("failed to describe services: %v", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Desired", "Running", "Pending", "Type", "Last Deployed"})

		for _, svc := range serviceDescs.Services {
			lastDeploy := "none"

			if len(svc.Deployments) > 0 {
				lastDeployTime := svc.Deployments[len(svc.Deployments)-1].CreatedAt
				lastDeploy = lastDeployTime.In(time.Local).Format(time.RFC822)
			}

			table.Append([]string{
				*svc.ServiceName,
				strconv.Itoa(int(svc.DesiredCount)),
				strconv.Itoa(int(svc.RunningCount)),
				strconv.Itoa(int(svc.PendingCount)),
				string(svc.LaunchType),
				lastDeploy,
			})
		}

		table.Render()
	},
}

func init() {
	registerClustersArg(getServicesCmd)
}
