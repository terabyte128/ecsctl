package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Tail logs for a service",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

		var serviceNames []string

		for _, svc := range serviceDescs.Services {
			serviceNames = append(serviceNames, *svc.ServiceName)
		}

		return serviceNames, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("service is required")
		}

		serviceDescs, err := client.DescribeServices(context.TODO(), &ecs.DescribeServicesInput{
			Services: []string{args[0]},
			Cluster:  &cluster,
		})

		if err != nil {
			log.Fatalf("failed to describe services: %v", err)
		}

		if len(serviceDescs.Services) == 0 {
			log.Fatalf("No service named %s", args[0])
		}

		service := serviceDescs.Services[0]

		taskDefinition, err := client.DescribeTaskDefinition(
			context.TODO(),
			&ecs.DescribeTaskDefinitionInput{
				TaskDefinition: service.TaskDefinition,
			},
		)

		if err != nil {
			log.Fatalf("failed to get task definitions: %v", err)
		}

		logConfiguration := taskDefinition.TaskDefinition.ContainerDefinitions[0].LogConfiguration

		if logConfiguration.LogDriver != "awslogs" {
			log.Fatal("log driver must be awslogs")
		}

		cloudwatch := cloudwatchlogs.NewFromConfig(awsConfig)
		groupPrefix := logConfiguration.Options["awslogs-group"]

		groups, err := cloudwatch.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: &groupPrefix,
		})

		if len(groups.LogGroups) != 1 {
			log.Fatalf("log group prefix %s returned multiple groups", groupPrefix)
		}

		groupArn := groups.LogGroups[0].LogGroupArn

		tail, err := cloudwatch.StartLiveTail(context.TODO(), &cloudwatchlogs.StartLiveTailInput{
			LogGroupIdentifiers: []string{*groupArn},
		})

		if err != nil {
			log.Fatalf("failed to start tail: %v", err)
		}

		stream := tail.GetStream()
		defer stream.Close()

		eventsChan := stream.Events()

		for {
			event := <-eventsChan
			switch e := event.(type) {
			case *types.StartLiveTailResponseStreamMemberSessionStart:
				log.Printf("Log stream connected to %s\n", *groupArn)
			case *types.StartLiveTailResponseStreamMemberSessionUpdate:
				for _, logEvent := range e.Value.SessionResults {
					fmt.Println(*logEvent.Message)
				}
			default:
				// Handle on-stream exceptions
				if err := stream.Err(); err != nil {
					log.Fatalf("Error occured during streaming: %v", err)
				} else if event == nil {
					log.Println("Stream is Closed")
					return
				} else {
					log.Fatalf("Unknown event type: %T", e)
				}
			}
		}

		fmt.Println(tail)
	},
}

func init() {
	registerClustersArg(logsCmd)
}
