package util

import "time"

// NewTicker returns a new Ticker which ticks repeatedly at the given duration.
// If insta is set to true, it will send a tick down its' channel on instantiation
func NewTicker(d time.Duration, insta bool) *Ticker {
	t := &Ticker{d, nil, insta, nil, nil}

	t.Start()

	return t
}

// Ticker behaves like time.Ticker, but also supports stopping and restarting
// without destroying the initial ticker, as well as sending the first tick
// immediately
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

// Start starts the Ticker, or resumes it from a stopped state.
// It will panic if called while the Ticker is in a started or killed state.
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

// Kill kills the ticker, preventing it from being started again.
// The Ticker will panic if Start() is called after Kill()
// Kill may be called when the ticker is in any state (other than killed)
func (t *Ticker) Kill() {
	t.Stop()
	close(t.out)
}

// Stop stops the ticker, with the intention of restarting it again.
// The outbound channel is not closed. If you are intending not to use the
// Ticker again, call Kill() to properly free resources. Stop may be called
// from any state.
func (t *Ticker) Stop() {
	if t.t != nil {
		t.stop()
	}
}

func (t *Ticker) stop() {
	t.t.Stop()
	t.t = nil
}
