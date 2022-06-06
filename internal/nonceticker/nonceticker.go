package nonceticker

import (
	"sync"
	"time"
)

var ticker = newNonceTicker()

type NonceTicker struct {
	ch   chan time.Time
	once sync.Once
}

func newNonceTicker() *NonceTicker {
	ch := make(chan time.Time, 100)

	return &NonceTicker{
		ch: ch,
	}
}

func (nt *NonceTicker) start() {
	nt.once.Do(func() {
		nt.run()
	})
}

func (nt *NonceTicker) GetTick() uint64 {
	t := <-nt.ch
	ms := uint64(t.UTC().UnixNano())

	return ms
}

func (nt *NonceTicker) run() {
	go func(nt *NonceTicker) {
		tick := time.NewTicker(time.Second)

		for t := range tick.C {
			nt.pushTime(t)
		}
	}(nt)
}

func (nt *NonceTicker) pushTime(t time.Time) {
	go func(nt *NonceTicker, t time.Time) {
		nt.ch <- t
	}(nt, t)
}

func GetNonceTicker() *NonceTicker {
	ticker.start()

	return ticker
}
