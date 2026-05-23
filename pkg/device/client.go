package device

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// App represents an application on the BB10/PlayBook device.
type App struct {
	ID      string `json:"id" xml:"id"`
	Name    string `json:"name" xml:"name"`
	Version string `json:"version" xml:"version"`
	Status  string `json:"status" xml:"status"`
}

// Client represents a connection to a BB10/PlayBook device.
type Client struct {
	BaseURL    string
	Password   string
	HTTPClient *http.Client
	Token      string
}

// NewClient creates a new device client.
func NewClient(ip string, password string, insecure bool) *Client {
	baseURL := fmt.Sprintf("https://%s:1337", ip)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			// InsecureSkipVerify allows connecting to devices with self-signed certificates.
			// This is common for BB10/PlayBook devices in developer mode.
			InsecureSkipVerify: insecure,
		},
	}

	return &Client{
		BaseURL:  baseURL,
		Password: password,
		HTTPClient: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

// Login authenticates with the device and retrieves a session token.
func (c *Client) Login() error {
	formData := url.Values{}
	formData.Set("password", c.Password)

	resp, err := c.HTTPClient.PostForm(c.BaseURL+"/auth", formData)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth failed with status: %s", resp.Status)
	}

	// Assuming the token is returned in a header or body.
	// For this implementation, we'll check common patterns.
	// Some devices might use cookies instead of a token header.

	// Mocking token retrieval for now as it depends on exact API behavior.
	c.Token = resp.Header.Get("X-Auth-Token")
	if c.Token == "" {
		// Try to read from body if not in header
		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
			c.Token = result["token"]
		}
	}

	return nil
}

// ListApps retrieves a list of installed applications.
func (c *Client) ListApps() ([]App, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/apps", nil)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("X-Auth-Token", c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list apps: %s", resp.Status)
	}

	var apps []App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, fmt.Errorf("failed to decode apps list: %w", err)
	}

	return apps, nil
}

// InstallApp uploads and installs a BAR file.
func (c *Client) InstallApp(barPath string) error {
	file, err := os.Open(barPath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(barPath))
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, file); err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", c.BaseURL+"/install", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.Token != "" {
		req.Header.Set("X-Auth-Token", c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("install failed with status: %s", resp.Status)
	}

	return nil
}

// UninstallApp uninstalls an application by its ID.
func (c *Client) UninstallApp(appID string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/uninstall?id=%s", c.BaseURL, appID), nil)
	if err != nil {
		return err
	}

	if c.Token != "" {
		req.Header.Set("X-Auth-Token", c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("uninstall failed with status: %s", resp.Status)
	}

	return nil
}

// ManageApp manages app state (launch/terminate).
func (c *Client) ManageApp(appID string, action string) error {
	// action can be "launch" or "terminate"
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s?id=%s", c.BaseURL, action, appID), nil)
	if err != nil {
		return err
	}

	if c.Token != "" {
		req.Header.Set("X-Auth-Token", c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s failed with status: %s", action, resp.Status)
	}

	return nil
}
