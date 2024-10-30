package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
)

const (
	// app argument default values
	defaultFluentHost           = "127.0.0.1"
	defaultFluentPort           = 24224
	defaultSteadyRate           = 10000
	defaultBurst                = 0
	defaultTime                 = 1
	defaultIsAsync              = true
	defaultIsSubSecondPrecision = true
)

var (
	logger    *fluent.Fluent
	tagFormat = "logger-firelens-%s"

	// args
	fluentHost           string
	fluentPort           int
	steadyRate           int
	burst                int
	timeInMinutes        int
	isAsync              bool
	isSubSecondPrecision bool
	printUsage           bool
)

type TaskMetadata struct {
	TaskARN string `json:"TaskARN"`
}

func getTaskID() string {
	metadataURI := os.Getenv("ECS_CONTAINER_METADATA_URI_V4")
	if metadataURI == "" {
		fmt.Println("ECS_CONTAINER_METADATA_URI_V4 environment variable is not set")
		os.Exit(1)
	}

	taskEndpoint := metadataURI + "/task"
	resp, err := http.Get(taskEndpoint)
	if err != nil {
		fmt.Printf("Error making request to TMDS: %v\n", err)
		os.Exit(1)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		os.Exit(1)
	}

	var taskMetadata TaskMetadata
	err = json.Unmarshal(body, &taskMetadata)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		os.Exit(1)
	}

	return extractTaskID(taskMetadata.TaskARN)
}

func extractTaskID(taskARN string) string {
	parts := strings.Split(taskARN, "/")
	return parts[len(parts)-1]
}

func parseArgs() {
	flag.StringVar(&fluentHost, "host", defaultFluentHost, "Fluent host address")
	flag.IntVar(&fluentPort, "port", defaultFluentPort, "Fluent port number")
	flag.IntVar(&steadyRate, "steadyRate", defaultSteadyRate, "Steady state rate")
	flag.IntVar(&burst, "burst", defaultBurst, "Burst value")
	flag.IntVar(&timeInMinutes, "time", defaultTime, "Time in minutes")
	flag.BoolVar(&isAsync, "async", defaultIsAsync, "Whether Fluent Async mode is enabled")
	flag.BoolVar(&isSubSecondPrecision, "subSecond", defaultIsSubSecondPrecision,
		"Whether Fluent Sub-second precision is enabled")
	flag.BoolVar(&printUsage, "help", false, "Print app usage")

	flag.Parse()
}

func initLogger() error {
	var err error
	logger, err = fluent.New(fluent.Config{
		FluentHost:         fluentHost,
		FluentPort:         fluentPort,
		Async:              isAsync,
		SubSecondPrecision: isSubSecondPrecision,
	})
	return err
}

// sendEvent sends a log event to Fluent bit directly
func sendEvent(tag string, i int) {
	data := map[string]interface{}{
		"message": fmt.Sprintf("Message: %d", i),
	}

	err := logger.Post(tag, data)
	if err != nil {
		fmt.Printf("Error sending event %d: %v\n", i, err)
	}
}

func main() {
	parseArgs()

	if printUsage {
		flag.Usage()
		os.Exit(0)
	}

	err := initLogger()
	if err != nil {
		fmt.Printf("Error initializing Fluent logger: %v\n", err)
		os.Exit(1)
	}
	defer func(logger *fluent.Fluent) {
		err := logger.Close()
		if err != nil {
			fmt.Printf("Error closing Fluent logger: %v\n", err)
		}
	}(logger)

	taskId := getTaskID()
	if taskId == "" {
		fmt.Println("Error getting task ID from TMDS")
		os.Exit(1)
	}
	tag := fmt.Sprintf(tagFormat, taskId)

	timeInSeconds := timeInMinutes * 60
	numEvents := timeInSeconds * steadyRate
	startTime := time.Now()
	for i := 1; i < numEvents; i++ {
		sendEvent(tag, i)
		if i%steadyRate == 0 {
			currentTime := time.Now()
			toWait := time.Second - currentTime.Sub(startTime)
			if toWait > 0 {
				time.Sleep(toWait)
			}
			startTime = time.Now()
		}
	}

	soFar := timeInSeconds * steadyRate
	eventsWithBurst := soFar + burst + 1

	for i := soFar; i < eventsWithBurst; i++ {
		sendEvent(tag, i)
	}
}
