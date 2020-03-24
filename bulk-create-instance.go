package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/zeebe-io/zeebe/clients/go/pkg/pb"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		convertedValue, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return convertedValue
	}
	return fallback
}

func deployWorkflow(client zbc.Client, filename string) (*pb.DeployWorkflowResponse, error) {
	return client.NewDeployWorkflowCommand().AddResourceFile(filename).Send(context.Background())
}

func createInstance(client zbc.Client, processID string, requestID int) (*pb.CreateWorkflowInstanceResponse, error) {
	variables := make(map[string]interface{})
	variables["requestId"] = requestID
	cmd, err := client.NewCreateInstanceCommand().BPMNProcessId(processID).LatestVersion().VariablesFromMap(variables)
	if err != nil {
		return nil, err
	}

	return cmd.Send(context.Background())
}

func main() {
	address := getEnv("ZEEBE_ADDRESS", "localhost:26500")
	processID := getEnv("ZEEBE_PROCESS_ID", "noop")
	workflowFilename := getEnv("ZEEEBE_WORKFLOW", "noop.bpmn")

	instances := getEnvInt("ZEEBE_INSTANCES", 10)

	client, err := zbc.NewClient(&zbc.ClientConfig{GatewayAddress: address})

	if err != nil {
		panic(err)
	}

	fmt.Println(deployWorkflow(client, workflowFilename))
	fmt.Println(createInstance(client, processID, 0))

	var wg sync.WaitGroup
	var success uint64
	var errors uint64

	wg.Add(instances)
	for i := 1; i <= instances; i++ {
		requestID := i
		go func() {
			defer wg.Done()
			_, err := createInstance(client, processID, requestID)
			if err == nil {
				atomic.AddUint64(&success, 1)
			} else {
				atomic.AddUint64(&errors, 1)
				if grpc.Code(err) != codes.ResourceExhausted {
					fmt.Println("Request ID", requestID, "Error", err)
				}
			}
		}()
	}

	wg.Wait()
	fmt.Println("Success:", success, "Error:", errors, "(", (float64(success)*100.0)/float64(success+errors), "%)")
}
