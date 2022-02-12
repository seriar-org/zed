package gzc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type API struct {
	HTTPClient *http.Client
	BaseURL    string
}

type Client struct {
	AuthToken string
	API       *API
}

func CreateAPI(c *http.Client, baseURL string) *API {
	return &API{
		HTTPClient: c,
		BaseURL:    baseURL,
	}
}

func (a *API) WithTimeout(timeout int) *API {
	td, _ := time.ParseDuration(fmt.Sprintf("%ds", timeout))
	a.HTTPClient.Timeout = td
	return a
}

func CreateClient(api *API, authToken string) *Client {
	return &Client{
		AuthToken: authToken,
		API:       api,
	}
}

func (c *Client) Request(req *http.Request) (resp *http.Response, err error) {

	req.Header.Set("X-Authentication-Token", c.AuthToken)

	resp, err = c.API.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("Request %v failed. Response was %v", req, resp)
	}
	return resp, nil
}

func (c *Client) createURL(resource string) string {
	return fmt.Sprintf("%s%s", c.API.BaseURL, resource)
}

type Estimate struct {
	Value int `json:"value"`
}
type Pipeline struct {
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	PipelineID  string `json:"pipeline_id"`
}
type Issue struct {
	IssueNumber int        `json:"issue_number"`
	IsEpic      bool       `json:"is_epic"`
	RepoID      int        `json:"repo_id"`
	Estimete    Estimate   `json:"estimate"`
	Pipelines   []Pipeline `json:"pipelines"`
}
type Epic struct {
	TotalEpicEstimates Estimate   `json:"total_epic_estimates"`
	Estimate           Estimate   `json:"estimate"`
	Pipelines          []Pipeline `json:"pipelines"`
	Issues             []Issue    `json:"issues"`
}
type Dependency struct {
	Blocking Issue `json:"blocking"`
	Blocked  Issue `json:"blocked"`
}
type Dependencies struct {
	Dependencies []Dependency `json:"dependencies"`
}

func (c *Client) RequestEpic(repoID, epicID int) (*Epic, error) {
	req, err := http.NewRequest("GET", c.createURL(fmt.Sprintf("/p1/repositories/%d/epics/%d", repoID, epicID)), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var epic Epic
	err = json.Unmarshal(body, &epic)
	if err != nil {
		return nil, err
	}

	return &epic, err
}

func (c *Client) RequestDependencies(repoID int) (*Dependencies, error) {
	req, err := http.NewRequest("GET", c.createURL(fmt.Sprintf("/p1/repositories/%d/dependencies", repoID)), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dependencies Dependencies
	err = json.Unmarshal(body, &dependencies)
	if err != nil {
		return nil, err
	}

	return &dependencies, err
}
