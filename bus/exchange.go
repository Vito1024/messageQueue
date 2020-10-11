package bus

import (
	"errors"
	"fmt"
	"messageQueue/sub"
	"sync"
)

var (
	// Exchange
	ErrExchangeAlreadyExist = errors.New("exchange already exist")
	ErrExchangeNotExist     = errors.New("exchange does not exist")

	// Subscriber
	ErrSubAlreadyExist = errors.New("subscriber already exist")
	ErrSubNotExist     = errors.New("subscriber does not exist")
)

type Exchange interface {
	AddSub(sub.Subscriber) error
	DelSub(sub.Subscriber) error
	Publish([]byte) error
	Consume() error
}

type exchange struct {
	// exchange name
	name string
	// For Add/Del subscribers
	rw sync.RWMutex
	// Subscribers subscribing current exchange
	subs []sub.Subscriber
	// Message storage
	msgs chan []byte
}

func newexchange(name string) *exchange {
	return &exchange{
		name: name,
		// By default, we set the length of msgs queue to 1000
		msgs: make(chan []byte, 1000),
	}
}

func (ex *exchange) AddSub(suber sub.Subscriber) error {
	ex.rw.Lock()
	defer ex.rw.Unlock()

	for _, _sub := range ex.subs {
		if _sub == suber {
			return ErrSubAlreadyExist
		}
	}

	ex.subs = append(ex.subs, suber)
	return nil
}

func (ex *exchange) DelSub(suber sub.Subscriber) error {
	ex.rw.Lock()
	defer ex.rw.Unlock()

	for idx, value := range ex.subs {
		if value == suber {
			ex.subs = append(ex.subs[:idx], ex.subs[idx+1:]...)

			// Because subscribers are unique
			return nil
		}
	}

	return ErrSubNotExist
}

func (ex *exchange) AllSubs() []sub.Subscriber {
	return ex.subs
}

func (ex *exchange) Publish(msg []byte) error {
	ex.msgs <- msg
	return nil
}

func (ex *exchange) Consume() error {
	for i := 0; i < 10; i++ {
		go func() {
			for { fmt.Println("[I'm consuming]:", string(<-ex.msgs)) }
		}()
	}

	return nil
}
