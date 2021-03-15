package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client client using gui socket
type Client struct {
	client http.Client
}

const timeout = 10 * time.Minute

// Restart restarts the engine.
func (c *Client) Restart() error {
	return c.post("engine/restart", &RestartParameters{}, nil)
}

type RestartParameters struct {
	OpenContainerView bool `json:"openContainerView"`
}

// ResetToFactoryDefaults resets to factory defaults (keeping the container mode)
func (c *Client) ResetToFactoryDefaults() error {
	return c.post("desktop/factory-reset", nil, nil)
}

// GetLastEvent gets the Kubernetes state to gui.
func (c *Client) GetLastEvent() (*Event, error) {
	var events []Event
	if err := c.get("vm/events", &events); err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, errors.New("no events found")
	}
	event := events[0]
	return &event, nil
}

func (c Client) get(route string, result interface{}) error {
	return doGet(c.client, "http://unix/"+route, result)
}

func (c Client) post(route string, payload, result interface{}) error {
	return doPost(c.client, "http://unix/"+route, payload, result)
}

func doPost(client http.Client, url string, payload, result interface{}) error {
	var body []byte
	if payload != nil {
		var err error
		body, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}
	//log.Printf("POST %s %v", url, string(body))
	res, err := client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := checkSuccess(res); err != nil {
		return err
	}
	respBody, err := ioutil.ReadAll(res.Body)
	if result == nil || err != nil {
		return err
	}
	return json.Unmarshal(respBody, result)
}

func doGet(client http.Client, url string, result interface{}) error {
	//log.Printf("GET %s", url)
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := checkSuccess(res); err != nil {
		return err
	}
	data, err := ioutil.ReadAll(res.Body)
	if result == nil || err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

func isSuccess(response *http.Response) bool {
	return response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices
}

func checkSuccess(response *http.Response) error {
	if isSuccess(response) {
		return nil
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Errorf("status code not OK but %d: %s", response.StatusCode, string(data))
	}
	return errors.Errorf("status code not OK but %d: %s", response.StatusCode, string(data))
}

// State vm state
type State string

const (
	// NotValid ...
	NotValid State = "not valid"
	// Stopped ...
	Stopped State = "stopped"
	// Stopping ...
	Stopping State = "stopping"
	// Starting ...
	Starting State = "starting"
	// FailedToStart ...
	FailedToStart State = "failed to start"
	// Running ...
	Running State = "running"
)

// Event vm event
type Event struct {
	State     State  `json:"state"`
	Timestamp int    `json:"date"`
	Mode      string `json:"mode,omitempty"`
}
