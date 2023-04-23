package main

import (
	"context"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func RunRedisContainers(t *testing.T) testcontainers.Container {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}

	return redisC
}

func TearDownContainers(ctx context.Context, t *testing.T, containers testcontainers.Container) {
	if err := containers.Terminate(ctx); err != nil {
		t.Fatalf("failed to terminate container: %s", err.Error())
	}
}

func TestSimpleRedis(t *testing.T) {
	ctx := context.Background()
	redisC := RunRedisContainers(t)
	defer TearDownContainers(ctx, t, redisC)

	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		t.Error(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	key := "hello"
	value := "world"

	client.Set(ctx, key, value, time.Minute)
	got, _ := client.Get(ctx, key).Result()

	assert.Equal(t, got, value)
}
