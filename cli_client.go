// Client interface using the official stackato client binary

package client

import (
	"fmt"
	"github.com/srid/log"
	"github.com/srid/run"
	"os/exec"
)

type CliClient struct {
	Endpoint string
	Token    string
	Group    string
}

func NewCliClient(endpoint string, token string, group string) (*CliClient, error) {
	if token == "" {
		return nil, fmt.Errorf("token string must not be empty")
	}
	c := &CliClient{endpoint, token, group}
	return c, nil
}

// PushAppNoCreate emulates `s push --no-create ...` and sends the
// output in outputCh channel.
func (c *CliClient) PushAppNoCreate(name string, dir string, outputCh chan string) error {
	// TODO: validate 'endpoint' and 'name' for security reasons.

	options := []string{
		name,
		"--no-tail", "--no-prompt",
		"--target", c.Endpoint,
		"--token", c.Token,
		"--path", dir}

	if c.Group != "" {
		options = append(options, "--group", c.Group)
	}

	pushOptions := append([]string{"push", "--no-create"}, options...)

	log.Infof("Deploying app %s", name)
	ret, err := run.Run(exec.Command("stackato", pushOptions...), outputCh)
	if err != nil {
		log.Error("cannot read line: ", err)
		return err
	}
	return ret
}
