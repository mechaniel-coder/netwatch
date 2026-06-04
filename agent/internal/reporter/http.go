package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/mechaniel-coder/netwatch/agent/internal/collector"
)

// Reporter handles all communication with the NetWatch server.
type Reporter struct {
	serverURL string
	token     string
	client    *http.Client
}

type registrationPayload struct {
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ip_address"`
	OS        string `json:"os"`
	Version   string `json:"version"`
}

type registrationResponse struct {
	AgentID string `json:"agent_id"`
}

func New(serverURL, token string) *Reporter {
	return &Reporter{
		serverURL: serverURL,
		token:     token,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Register calls /api/v1/agents/register and returns the assigned agent ID.
func (r *Reporter) Register() (string, error) {
	payload := registrationPayload{
		Hostname:  hostname(),
		IPAddress: "0.0.0.0", // server can override from request IP
		OS:        runtime.GOOS,
		Version:   "1.0.0",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, r.serverURL+"/api/v1/agents/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("register request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("register: server returned %d", resp.StatusCode)
	}

	var out registrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode registration response: %w", err)
	}
	return out.AgentID, nil
}

// Send ships a metric snapshot to the server.
func (r *Reporter) Send(agentID string, snap *collector.Snapshot) error {
	body, _ := json.Marshal(snap)
	url := fmt.Sprintf("%s/api/v1/agents/%s/metrics", r.serverURL, agentID)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send: server returned %d", resp.StatusCode)
	}
	return nil
}

func hostname() string {
	// TODO: use os.Hostname() — placeholder for cross-platform handling
	return "unknown-host"
}
