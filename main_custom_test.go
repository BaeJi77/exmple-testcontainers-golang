package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/testcontainers/testcontainers-go"
)

func RunCustomContainers(t *testing.T) testcontainers.Container {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "testcontainers-go-test:1.0.0",
		ExposedPorts: []string{"10000/tcp"},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}

	return container
}

func TestSimpleCustomContainers(t *testing.T) {
	ctx := context.Background()
	customContainers := RunCustomContainers(t)
	defer TearDownContainers(ctx, t, customContainers)

	endpoint, err := customContainers.Endpoint(ctx, "")
	if err != nil {
		t.Error(err)
	}

	res, err := http.DefaultClient.Get(fmt.Sprintf("http://%s", endpoint))
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
