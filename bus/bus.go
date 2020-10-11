package bus

import (
	"messageQueue/sub"
	"os"
	"sync"
)

type bus struct {
	// For Add/Del exchanges
	rw        sync.RWMutex
	exchanges map[string]Exchange
}

func newbus() *bus{
	return &bus{
		exchanges: make(map[string]Exchange),
	}
}

func (b *bus) publish(msg []byte, exchange Exchange) error {
	return exchange.Publish(msg)
}

func (b *bus) consume() error {
	for _, ex := range b.exchanges {
		go ex.Consume()
	}

	return nil
}

func (b *bus) addExchange(name string) error {
	if _, ok := b.exchanges[name]; ok {
		return ErrExchangeAlreadyExist
	}

	b.exchanges[name] = newexchange(name)
	return nil
}

var b = newbus()

func Start() {
	var err error

	if err = b.consume(); err != nil {
		os.Exit(127)
	}
}

func AddExchange(name string) error {
	if err := b.addExchange(name); err != nil {
		return err
	}

	go b.exchanges[name].Consume()
	return nil
}

func AddSubsToExchange(exchange string, subs ...sub.Subscriber) error {
	if ex, ok := b.exchanges[exchange]; ok {
		for _, suber := range subs {
			if err := ex.AddSub(suber); err != nil {
				return err
			}
		}
	}

	return ErrExchangeNotExist
}

func Publish(msg []byte, exchange string) error {
	if ex, ok := b.exchanges[exchange]; ok {
		return ex.Publish(msg)
	}

	return ErrExchangeNotExist
}