package or

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func TestOr(t *testing.T) {
	tests := []struct {
		name     string
		channels []<-chan interface{}
		maxTime  time.Duration
	}{
		{
			name: "single channel",
			channels: []<-chan interface{}{
				sig(50 * time.Millisecond),
			},
			maxTime: 100 * time.Millisecond,
		},
		{
			name: "many channels",
			channels: []<-chan interface{}{
				sig(200 * time.Millisecond),
				sig(30 * time.Millisecond),
				sig(300 * time.Millisecond),
			},
			maxTime: 80 * time.Millisecond,
		},
		{
			name:     "nil for zero channels",
			channels: []<-chan interface{}{},
			maxTime:  0,
		},
		{
			name: "immediately closed channel closes result",
			channels: func() []<-chan interface{} {
				c := make(chan interface{})
				close(c)
				return []<-chan interface{}{c, sig(500 * time.Millisecond)}
			}(),
			maxTime: 20 * time.Millisecond,
		},
		{
			name: "many slow channels",
			channels: []<-chan interface{}{
				sig(150 * time.Millisecond),
				sig(500 * time.Millisecond),
				sig(300 * time.Millisecond),
			},
			maxTime: 200 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ch := Or(tt.channels...)

			if len(tt.channels) == 0 {
				require.Nil(t, ch)
				return
			}

			start := time.Now()

			<-ch

			elapsed := time.Since(start)

			require.LessOrEqualf(t, elapsed, tt.maxTime,
				"should close within %v but took %v", tt.maxTime, elapsed)
		})
	}
}
