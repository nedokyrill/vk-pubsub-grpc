package subpub

import (
	"context"
	"fmt"
	"github.com/nedokyrill/vk-pubsub-grpc/pkg/logger"
	"sync"
)

type MessageHandler func(msg interface{})

type Subscription interface {
	Unsubscribe()
}

type SubPub interface {
	Subscribe(subject string, cb MessageHandler) (Subscription, error)
	Publish(subject string, msg interface{}) error
	Close(ctx context.Context) error
}

type SubPubAgent struct {
	mu    sync.Mutex
	subs  map[string][]*SubscriptionAgent
	quit  chan struct{}
	close bool
}

type SubscriptionAgent struct {
	ch      chan interface{}
	mh      MessageHandler
	Done    chan struct{}
	subject string
	spa     *SubPubAgent
}

func NewSubPubAgent() *SubPubAgent {
	return &SubPubAgent{
		subs: make(map[string][]*SubscriptionAgent),
		quit: make(chan struct{}),
	}
}

func NewSubscriptionAgent(mh MessageHandler, spa *SubPubAgent) *SubscriptionAgent {
	return &SubscriptionAgent{
		ch:   make(chan interface{}, 64),
		mh:   mh,
		Done: make(chan struct{}),
		spa:  spa,
	}
}

func (s *SubscriptionAgent) Unsubscribe() {
	s.spa.mu.Lock()
	defer s.spa.mu.Unlock()

	logger.Logger.Info("Unsubscribing from subject ", s.subject)

	subs := s.spa.subs[s.subject]
	for i, v := range subs {
		if v == s {
			s.spa.subs[s.subject] = append(subs[:i], subs[i+1:]...)

			// Проверка: если канал не закрыт, закрываем
			select {
			case <-v.Done:
			default:
				close(v.Done)
			}

			select {
			case <-v.ch:
			default:
				close(v.ch)
			}

			logger.Logger.Info("Successfully unsubscribed from subject ", s.subject)
			break
		}
	}
}

func (s *SubPubAgent) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.close {
		logger.Logger.Error("Subscribe failed: SubPubAgent already closed for subject ", subject)
		return nil, fmt.Errorf("publisher already closed")
	}

	if _, ok := s.subs[subject]; !ok {
		s.subs[subject] = []*SubscriptionAgent{}
	}

	sub := NewSubscriptionAgent(cb, s)
	sub.subject = subject
	s.subs[subject] = append(s.subs[subject], sub)

	logger.Logger.Info("Subscribed to subject ", subject)

	go func(subIn *SubscriptionAgent) {
		for {
			select {
			case msg := <-subIn.ch:
				subIn.mh(msg)
			case <-subIn.Done:
				logger.Logger.Info("Subscription done for subject ", subIn.subject)
				return
			case <-s.quit:
				logger.Logger.Info("Global shutdown: unsubscribing from subject ", subIn.subject)
				subIn.Unsubscribe()
				return
			}
		}
	}(sub)

	return sub, nil
}

func (s *SubPubAgent) Publish(subject string, msg interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.close {
		logger.Logger.Error("Publish failed: SubPubAgent already closed for subject ", subject)
		return fmt.Errorf("publisher already closed")
	}

	logger.Logger.Info("Publishing message to subject ", subject)

	for _, sub := range s.subs[subject] {
		select {
		case sub.ch <- msg:
		case <-sub.Done:
			logger.Logger.Info("Skipping dead subscription for subject ", subject)
		}
	}

	return nil
}

func (s *SubPubAgent) Close(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.close {
		logger.Logger.Error("SubPubAgent already closed")
		return fmt.Errorf("publisher already closed")
	}

	logger.Logger.Info("SubPubAgent closing")

	success := make(chan struct{})
	go func() {
		s.close = true
		close(s.quit)

		for _, subject := range s.subs {
			for _, sub := range subject {
				close(sub.ch)
				close(sub.Done)
				logger.Logger.Info("Closed subscription on subject ", sub.subject)
			}
		}

		close(success)
	}()

	select {
	case <-success:
		logger.Logger.Info("SubPubAgent closed successfully")
		return nil
	case <-ctx.Done():
		logger.Logger.Error("SubPubAgent close timed out: ", ctx.Err())
		return ctx.Err()
	}
}

func isClosed(ch chan interface{}) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}
