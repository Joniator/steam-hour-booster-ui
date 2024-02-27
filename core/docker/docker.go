package docker

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	apiClient     client.APIClient
	ContainerName string
}

func New(containerName string) DockerClient {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("Failed to initialize docker: %v", err)
	}
	client.NegotiateAPIVersion(context.Background())

	return DockerClient{
		apiClient:     client,
		ContainerName: containerName,
	}
}

func (dc DockerClient) IsAvailable() bool {
	return dc.apiClient != nil
}

func (dc DockerClient) GetStatus() string {
	containers, err := dc.apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return err.Error()
	}
	for _, c := range containers {
		if slices.Contains(c.Names, fmt.Sprintf("/%s", dc.ContainerName)) {
			return c.Status
		}
	}
	return ""
}

func (dc DockerClient) GetLogs() []string {
	logs, err := dc.apiClient.ContainerLogs(context.Background(), dc.ContainerName, container.LogsOptions{Tail: "10", ShowStdout: true, Details: false})
	if err != nil {
		return []string{"Failed to get logs", err.Error()}
	}

	headerbytes := make([]byte, 8)
	buf := new(bytes.Buffer)
	buf.Read(headerbytes)
	buf.ReadFrom(logs)
	raw := buf.String()
	if err != nil {
		return []string{"Failed to get logs", err.Error()}
	}
	s := string(raw)
	s = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))").ReplaceAllString(s, "")
	s = regexp.MustCompile(`[^a-zA-Z0-9\-\(\)\n\[\] ]+`).ReplaceAllString(s, "")
	return strings.Split(s, "\n")
}

func (dc DockerClient) Start() error {
	err := dc.apiClient.ContainerStart(context.Background(), dc.ContainerName, container.StartOptions{})
	if err != nil {
		log.Printf("Failed to start container: %v", err)
	}
	return err
}

func (dc DockerClient) Restart() error {
	err := dc.apiClient.ContainerRestart(context.Background(), dc.ContainerName, container.StopOptions{})
	if err != nil {
		log.Printf("Failed to restart container: %v", err)
	}
	return err
}

func (dc DockerClient) Stop() error {
	err := dc.apiClient.ContainerStop(context.Background(), dc.ContainerName, container.StopOptions{})
	if err != nil {
		log.Printf("Failed to stop container: %v", err)
	}
	return err
}
