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

func (c *RestClient) GetLogs(appGUID string, num int) ([]AppLogLine, error) {
	path := fmt.Sprintf("/v2/apps/%s/stackato_logs?num=%d&monolith=1", appGUID, num)
	var response struct {
		Lines []AppLogLine `json:"lines"`
	}
	err := c.MakeRequest("GET", path, nil, &response)
	return response.Lines, err
}

func (c *RestClient) GetLogsRaw(appGUID string, num int) ([]string, error) {
	path := fmt.Sprintf(
		"/v2/apps/%s/stackato_logs?num=%d&as_is=1&monolith=1", appGUID, num)
	var response struct {
		Lines []string `json:"lines"`
	}
	err := c.MakeRequest("GET", path, nil, &response)
	return response.Lines, err
}

func (c *RestClient) ListApps() (apps []App, err error) {
	if c.Space == "" {
		panic("empty Space")
	}
	path := fmt.Sprintf("/v2/spaces/%s/summary", c.Space)
	var response struct {
		GUID string
		Name string
		Apps []App
	}
	response.Apps = apps
	err = c.MakeRequest("GET", path, nil, &response)
	return
}

// CreateApp only creates the application. It is an equivalent of `s
// create-app --json`.
func (c *RestClient) CreateApp(name string) (string, error) {
	// Ensure that app name is unique for this user. We do this as
	// unfortunately the server doesn't enforce it.
	apps, err := c.ListApps()
	if err != nil {
		return "", err
	}
	for _, app := range apps {
		if app.Name == name {
			return "", fmt.Errorf("App by that name (%s) already exists", name)
		}
	}

	// The CC requires that a POST on /apps sends, at minimum, these
	// fields. The values for framework/runtime doesn't matter for our
	// purpose (they will get overwritten by a subsequent app push).
	createArgs := map[string]interface{}{
		"name":       name,
		"space_guid": c.Space,
	}

	var resp struct {
		Metadata struct {
			GUID string
		}
	}
	err = c.MakeRequest("POST", "/v2/apps", createArgs, &resp)
	if err != nil {
		return "", err
	}

	if resp.Metadata.GUID == "" {
		return "", fmt.Errorf("Missing App GUID from CC")
	}

	return resp.Metadata.GUID, nil
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
