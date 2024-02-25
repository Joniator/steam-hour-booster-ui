package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var apiClient, err = client.NewClientWithOpts(client.FromEnv)

type DockerClient struct {
	ContainerName string
}

func New(containerName string) DockerClient {
	return DockerClient{
		ContainerName: containerName,
	}
}

func (dc DockerClient) IsAvailable() bool {
	return err != nil
}

func (dc DockerClient) Restart() error {
	if err != nil {
		return err
	}
	err := apiClient.ContainerRestart(context.Background(), dc.ContainerName, container.StopOptions{})
	return err
}
