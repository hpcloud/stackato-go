package client

import (
	"crypto/tls"
	"fmt"
	"github.com/ActiveState/httpapi/client"
	"net/http"
)

type RestClient struct {
	TargetURL string
	Token     string
	Group     string
	client    *client.Client
}

func NewRestClient(targetUrl, token, group string) *RestClient {
	return &RestClient{targetUrl, token, group, getInsecureHttpRestClient()}
}

type App struct {
	Name string
	URIs []string
	Instances int
	RunningInstances *int
	State string
	Version string
	Staging struct {
		Model string
		Stack string
	}
	Resources struct {
		Memory int
		Disk int
		FDs int
		Sudo bool
	}
	Meta struct {
		Version int
		Created int  // timestamp
	}
}

func (c *RestClient) ListApps() (apps []App, err error) {
	err = c.MakeRequest("GET", "/apps", nil, &apps)
	return
}

// CreateApp only creates the application. It is an equivalent of `s
// create-app --json`.
func (c *RestClient) CreateApp(name string) (int, error) {
	// Ensure that app name is unique for this user. We do this as
	// unfortunately the server doesn't enforce it.
	apps, err := c.ListApps()
	if err != nil {
		return -1, err
	}
	for _, app := range apps {
		if app.Name == name {
			return -1, fmt.Errorf("App by that name (%s) already exists", name)
		}
	}
	
	// The CC requires that a POST on /apps passes, at minimum, these
	// fields. The values for framework/runtime doesn't matter for our
	// purposes (they will get overwritten by a subsequent app push).
	createArgs := map[string]interface{}{
		"name": name,
		"staging": map[string]string{
			"framework": "buildpack",
			"runtime":   "python27",
		},
	}

	var resp struct { App_ID int }
	err = c.MakeRequest("POST", "/apps", createArgs, &resp)
	if err != nil {
		return -1, err
	}

	if resp.App_ID < 1 {
		return -1, fmt.Errorf("Invalid or missing AppID from CC: %d", resp.App_ID)
	}

	return resp.App_ID, nil
}

func (c *RestClient) MakeRequest(method string, path string, params interface{}, response interface{}) error {
	req, err := client.NewRequest(method, c.TargetURL+path, params)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.Token)
	if c.Group != "" {
		req.Header.Set("X-Stackato-Group", c.Group)
	}
	return c.client.DoRequest(req, response)
}

// emulate `curl -k ...`
func getInsecureHttpRestClient() *client.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	return &client.Client{Transport: tr}
}
