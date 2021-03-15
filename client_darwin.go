package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"log"
)

const guiAPISocketName = "gui-api.sock"

// NewClient constructs a client for to-GUI socket.
func NewClient() *Client {
	return newClient(data(guiAPISocketName))
}

// NewClient constructs a client for the to-GUI socket.
func newClient(socket string) *Client {
	return &Client{
		client: http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", socket)
				},
			},
		},
	}
}

func data(elem ...string) string {
	container := filepath.Join(home(), "Library", "Containers", "com.docker.docker")
	data := filepath.Join(container, "Data")
	path := filepath.Join(append([]string{data}, elem...)...)
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("cannot compute path to %q relative to cwd: %v", path, err)
		return path
	}
	return strings.TrimPrefix(path, wd+"/")
}

func home() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
