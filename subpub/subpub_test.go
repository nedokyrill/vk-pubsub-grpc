package subpub_test

import (
	"context"
	"github.com/nedokyrill/vk-pubsub-grpc/pkg/logger"
	"go.uber.org/zap"
	"sync"
	"testing"
	"time"

	"github.com/nedokyrill/vk-pubsub-grpc/subpub"
	"github.com/stretchr/testify/assert"
)

func TestPublishSubscribe(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()

	agent := subpub.NewSubPubAgent()
	defer agent.Close(context.Background())

	var (
		wg      sync.WaitGroup
		recvMsg interface{}
	)

	wg.Add(1)
	_, err := agent.Subscribe("test-key", func(msg interface{}) {
		recvMsg = msg
		wg.Done()
	})
	assert.NoError(t, err)

	err = agent.Publish("test-key", "hello")
	assert.NoError(t, err)

	waitWithTimeout(t, &wg, time.Second)

	assert.Equal(t, "hello", recvMsg)
}

func TestUnsubscribe(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()

	agent := subpub.NewSubPubAgent()
	defer agent.Close(context.Background())

	var received bool
	sub, err := agent.Subscribe("key", func(msg interface{}) {
		received = true
	})
	assert.NoError(t, err)

	sub.Unsubscribe()
	time.Sleep(100 * time.Millisecond)

	err = agent.Publish("key", "data")
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	assert.False(t, received)
}

func TestPublishAfterClose(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()

	agent := subpub.NewSubPubAgent()

	err := agent.Close(context.Background())
	assert.NoError(t, err)

	err = agent.Publish("key", "data")
	assert.Error(t, err)
}

func TestSubscribeAfterClose(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()

	agent := subpub.NewSubPubAgent()

	err := agent.Close(context.Background())
	assert.NoError(t, err)

	_, err = agent.Subscribe("key", func(msg interface{}) {})
	assert.Error(t, err)
}

func TestCloseWithTimeout(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()

	agent := subpub.NewSubPubAgent()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	err := agent.Close(ctx)
	assert.NoError(t, err)
}

// Вспомогательная функция для ожидания с таймаутом
func waitWithTimeout(t *testing.T, wg *sync.WaitGroup, timeout time.Duration) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatal("timeout waiting for waitgroup")
	}
}
