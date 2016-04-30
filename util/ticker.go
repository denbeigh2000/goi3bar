package util

import "time"

func NewTicker(d time.Duration, insta bool) *Ticker {
	t := &Ticker{d, nil, insta, nil, nil}

	t.Start()

	return t
}

type Ticker struct {
	// Duration represents frequently the ticker should tick
	Duration time.Duration
	// Ticks are sent down this channel
	C <-chan time.Time

	// If set to true, a tick is sent down the channel immediately
	InstaTick bool

	t   *time.Ticker
	out chan time.Time
}

func (t *Ticker) Start() {
	if t.t != nil {
		panic("Ticker already started")
	}

	t.start()
}

func (t *Ticker) start() {
	t.t = time.NewTicker(t.Duration)

	out := make(chan time.Time)
	t.C = out
	t.out = out

	go func() {
		if t.InstaTick {
			out <- time.Now()
		}

		for now := range t.t.C {
			out <- now
		}
	}()
}

func (t *Ticker) Kill() {
	t.Stop()
	close(t.out)
}

func (t *Ticker) Stop() {
	if t.t != nil {
		t.stop()
	}
}

func (t *Ticker) stop() {
	t.t.Stop()
	t.t = nil
}
