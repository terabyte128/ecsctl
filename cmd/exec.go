package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

type sessionData struct {
	SessionId  string `json:"sessionId"`
	StreamUrl  string `json:"streamUrl"`
	TokenValue string `json:"tokenValue"`
}

type SsmRequestParams struct {
	Target string `json:"Target"`
}

func buildSsmRequestParams(rsp *ecs.ExecuteCommandOutput) string {
	splitArn := strings.Split(*rsp.ClusterArn, "/")
	splitTask := strings.Split(*rsp.TaskArn, "/")

	clusterName := splitArn[len(splitArn)-1]
	taskID := splitTask[len(splitTask)-1]
	containerName := *rsp.ContainerName

	taskRsp, err := client.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
		Cluster: &clusterName,
		Tasks:   []string{taskID},
	})
	if err != nil {
		log.Fatalf("failed to describe tasks: %v", err)
	}
	if len(taskRsp.Tasks) == 0 {
		log.Fatalf("no tasks with ID %s", taskID)
	}

	var containerRuntimeID *string

	for _, container := range taskRsp.Tasks[0].Containers {
		if *container.Name == containerName {
			containerRuntimeID = container.RuntimeId
		}
	}
	if containerRuntimeID == nil {
		log.Fatalf("no containers with name %s", containerName)
	}

	target := fmt.Sprintf("ecs:%s_%s_%s", clusterName, taskID, *containerRuntimeID)
	params := SsmRequestParams{
		Target: target,
	}

	marshalled, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("failed to marshal params %v", err)
	}

	return string(marshalled)
}

var execCmd = &cobra.Command{
	Use:   "exec command",
	Short: "Execute an interactive command on an ECS container",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("You must specify a command to run")
		}

		input := ecs.ExecuteCommandInput{
			Task:        &taskID,
			Cluster:     &cluster,
			Command:     &args[0],
			Interactive: true,
		}

		if container != "" {
			input.Container = &container
		}

		rsp, err := client.ExecuteCommand(context.TODO(), &input)
		if err != nil {
			log.Fatalf("failed to execute command: %v", err)
		}

		data := sessionData{
			SessionId:  *rsp.Session.SessionId,
			StreamUrl:  *rsp.Session.StreamUrl,
			TokenValue: *rsp.Session.TokenValue,
		}

		marshalledData, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("failed to marshal session data: %v", err)
		}

		ssmRequestParams := buildSsmRequestParams(rsp)
		sessionManagerArgs := []string{
			string(marshalledData),
			client.Options().Region,
			"StartSession",
			"", // profile (ignored)
			ssmRequestParams,
			fmt.Sprintf("https://ecs.%s.amazonaws.com", client.Options().Region),
		}

		smCmd := exec.Command("session-manager-plugin", sessionManagerArgs...)
		smCmd.Stdin = os.Stdin
		smCmd.Stdout = os.Stdout
		smCmd.Stderr = os.Stderr

		smCmd.Run()
	},
}

var taskID string
var container string

func init() {
	registerClustersArg(execCmd)
	execCmd.Flags().StringVarP(&taskID, "task-id", "t", "", "task ID")
	execCmd.MarkFlagRequired("task-id")
	execCmd.RegisterFlagCompletionFunc("task-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		tasks, err := client.ListTasks(context.TODO(), &ecs.ListTasksInput{
			Cluster: &cluster,
		})
		if err != nil {
			fmt.Printf("failed to list tasks %v", err)
		}

		var taskIDs []string

		for _, task := range tasks.TaskArns {
			splitArn := strings.Split(task, "/")
			id := splitArn[len(splitArn)-1]
			taskIDs = append(taskIDs, id)
		}

		return taskIDs, cobra.ShellCompDirectiveNoFileComp
	})

	execCmd.Flags().StringVar(&container, "container", "", "container name")
}
