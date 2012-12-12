package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RestClient struct {
	TargetURL    string
	Token  string
	Group  string
	client *http.Client
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

// TODO: extract these into srid/httpapi

// A version of http.NewRequest including token and group params; and taking path and params.
func (c *RestClient) NewRequest(method string, path string, params map[string]interface{}) (*http.Request, error) {
	url := c.TargetURL + path

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Token)
	if c.Group != "" {
		req.Header.Set("X-Stackato-Group", c.Group)
	}
	return req, nil
}

func (c *RestClient) MakeRequest(method string, path string, params map[string]interface{}) (map[string]interface{}, error) {
	req, err := c.NewRequest(method, path, params)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return c.ParseResponse(resp)
}

func (c *RestClient) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	err = json.Unmarshal(data, &body)
	if err != nil {
		return nil, fmt.Errorf("Expecting JSON http response; %s", err)
	}

	// XXX: accept other codes
	if !(resp.StatusCode == 200 || resp.StatusCode == 302) {
		return nil, fmt.Errorf("HTTP request with failure code (%d); body -- %v",
			resp.StatusCode, body)
	}

	return body, nil
}

// emulate `curl -k ...`
func getInsecureHttpRestClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	return &http.Client{Transport: tr}
}
