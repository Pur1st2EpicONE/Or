// Package or provides a function to combine multiple channels into a single channel
// that closes when any of the input channels closes or receives a value.
package or

import "sync"

// Or returns a channel that closes when any of the input channels closes or receives a value.
// If called with 0 channels, returns nil. If called with 1 channel, returns it unchanged.
func Or(channels ...<-chan any) <-chan any {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan any)
	go func() {
		var once sync.Once
		for _, channel := range channels {
			go func(channel <-chan any) {
				select {
				case <-channel:
					once.Do(func() { close(orDone) })
				case <-orDone:
				}
			}(channel)
		}
	}()
	return orDone
}
