package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"github.com/srid/httpapi/client"
)

type RestClient struct {
	TargetURL    string
	Token  string
	Group  string
	client *client.Client
}

func NewRestClient(targetUrl, token, group string) *RestClient {
	return &RestClient{targetUrl, token, group, getInsecureHttpRestClient()}
}

// CreateApp only creates the application. It is an equivalent of `s
// create-app --json`.
func (c *RestClient) CreateApp(name string) (int, error) {
	// POST on /apps seems to required these at minimum.
	createArgs := map[string]interface{}{
		"name": name,
		"staging": map[string]string{
			"framework": "buildpack", // server requires this field to be set.
			"runtime":   "python27",
		},
	}
	response, err := c.MakeRequest("POST", "/apps", createArgs)
	if err != nil {
		return -1, err
	}
	if app_id, ok := response["app_id"].(float64); ok {
		return int(app_id), nil
	} else {
		return -1, fmt.Errorf("Invalid json response from the app-create API")
	}
	return -1, err
}

func (c *RestClient) MakeRequest(method string, path string, params client.Hash) (client.Hash, error) {
	req, err := client.NewRequest(method, c.TargetURL + path, params)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Token)
	if c.Group != "" {
		req.Header.Set("X-Stackato-Group", c.Group)
	}
	return c.client.DoRequest(req)
}

// emulate `curl -k ...`
func getInsecureHttpRestClient() *client.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	return &client.Client{Transport: tr}
}
