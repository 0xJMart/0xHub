package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a client for the backend API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Project represents a project in the backend API
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Icon        string `json:"icon,omitempty"`
	Category    string `json:"category,omitempty"`
	Status      string `json:"status,omitempty"`
}

// NewClient creates a new backend client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateProject creates a project in the backend
func (c *Client) CreateProject(project *Project) error {
	url := fmt.Sprintf("%s/api/projects", c.baseURL)
	return c.doRequest(http.MethodPost, url, project, nil)
}

// UpdateProject updates a project in the backend
func (c *Client) UpdateProject(id string, project *Project) error {
	url := fmt.Sprintf("%s/api/projects/%s", c.baseURL, id)
	return c.doRequest(http.MethodPut, url, project, nil)
}

// DeleteProject deletes a project from the backend
func (c *Client) DeleteProject(id string) error {
	url := fmt.Sprintf("%s/api/projects/%s", c.baseURL, id)
	return c.doRequest(http.MethodDelete, url, nil, nil)
}

// GetProject retrieves a project from the backend
func (c *Client) GetProject(id string) (*Project, error) {
	url := fmt.Sprintf("%s/api/projects/%s", c.baseURL, id)
	var project Project
	err := c.doRequest(http.MethodGet, url, nil, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// HealthCheck checks if the backend is healthy
func (c *Client) HealthCheck() error {
	url := fmt.Sprintf("%s/api/health", c.baseURL)
	return c.doRequest(http.MethodGet, url, nil, nil)
}

func (c *Client) doRequest(method, url string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("backend API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

