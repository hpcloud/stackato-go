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
	Space     string
	client    *client.Client
}

func NewRestClient(targetUrl, token, space string) *RestClient {
	if space == "" {
		panic("empty Space")
	}
	if token == "" {
		panic("empty Token")
	}
	return &RestClient{targetUrl, token, space, getInsecureHttpRestClient()}
}

type App struct {
	GUID              string
	Name              string
	URLs              []string
	Instances         int
	RunningInstances  *int
	Version           string
	Buildpack         string
	DetectedBuildpack string
	Memory            int
	DiskQuota         int
}

func (c *RestClient) ListApps() (apps []App, err error) {
	path := fmt.Sprintf("/v2/spaces/%s/summary", c.Space)
	err = c.MakeRequest("GET", path, nil, &apps)
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

	// The CC requires that a POST on /apps sends, at minimum, these
	// fields. The values for framework/runtime doesn't matter for our
	// purpose (they will get overwritten by a subsequent app push).
	createArgs := map[string]interface{}{
		"name": name,
		"staging": map[string]string{
			"framework": "buildpack",
			"runtime":   "python27",
		},
	}

	var resp struct{ App_ID int }
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
	err = c.client.DoRequest(req, response)
	if err != nil {
		return fmt.Errorf("CC API %v %v failed: %v", method, path, err)
	}
	return nil
}

// emulate `curl -k ...`
func getInsecureHttpRestClient() *client.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	return &client.Client{Transport: tr}
}
