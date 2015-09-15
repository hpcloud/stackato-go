// Client interface using the official stackato client binary

package client

import (
	"fmt"
	"os/exec"

	"github.com/hpcloud/log"
	"github.com/hpcloud/run"
)

type CliClient struct {
	TargetURL string
	Token     string
	Space     string
}

func NewCliClient(targetUrl, token, space string) (*CliClient, error) {
	if token == "" {
		return nil, fmt.Errorf("Token must not be empty")
	}
	if space == "" {
		return nil, fmt.Errorf("Space must not be empty")
	}
	c := &CliClient{targetUrl, token, space}
	return c, nil
}

// PushAppNoCreate emulates `s push ...` and sends the
// output in outputCh channel.
func (c *CliClient) PushAppNoCreate(name string, dir string, autoStart bool, outputCh chan string) (bool, error) {
	options := []string{
		"push",
		"--name", name,
		"--no-tail",
		"--no-prompt",
		"--target", c.TargetURL,
		"--token", c.Token,
		"--space-guid", c.Space,
		"--path", dir}

	if !autoStart {
		options = append(options, "--no-start")
	}

	cmd := exec.Command("stackato", options...)
	cmd.Dir = dir
	// Avoid messing up ~stackato/.stackato/client
	// cmd.Env = []string{"HOME=/tmp"}

	// log.Infof("Running command: %+v", cmd)
	ret, err := run.Run(cmd, outputCh)
	if err != nil {
		log.Error("cannot read line: ", err)
		return false, err
	}
	if r, ok := ret.(*exec.ExitError); ok {
		log.Errorf("Client exited abruptly: %v", r)
		return false, nil
	} else {
		return true, ret
	}
}
