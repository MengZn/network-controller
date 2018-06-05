package docker

import (
	docker "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"golang.org/x/net/context"
)

type Client struct {
	*docker.Client
	Context context.Context
}

func New() (*Client, error) {
	cli, err := docker.NewEnvClient()
	return &Client{
		cli,
		context.Background(),
	}, err
}

func (c *Client) ListContainer() ([]types.Container, error) {
	return c.Client.ContainerList(c.Context, types.ContainerListOptions{})
}

// docker inspect <Container ID>
func (c *Client) InspectContainer(containerID string) (types.ContainerJSON, error) {
	return c.Client.ContainerInspect(c.Context, containerID)
}

// docker ps -a <Container ID>
func ListContainer(cli *docker.Client) ([]types.Container, error) {
	return cli.ContainerList(context.Background(), types.ContainerListOptions{})
}

// docker inspect <Container ID>
func InspectContainer(cli *docker.Client, containerID string) (types.ContainerJSON, error) {
	return cli.ContainerInspect(context.Background(), containerID)
}

// docker inspect <Container ID> | grep -E 'SandboxKey|Id'
func GetSandboxKey(containerInfo types.ContainerJSON) string {
	return containerInfo.NetworkSettings.SandboxKey
}
