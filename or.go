// Package or provides a function to combine multiple channels into a single channel
// that closes when any of the input channels closes or receives a value.
package or

import "sync"

// Or returns a channel that closes when any of the provided input channels
// close or send a value. It is useful for combining multiple cancellation
// or signaling channels into a single channel.
//
// Example usage:
//
//	done := Or(ch1, ch2, ch3)
//	<-done // unblocks when any of ch1, ch2, or ch3 is closed or sends a value
//
// The function ignores any nil channels. If no channels are provided, it
// returns a closed channel immediately. If exactly one channel is provided,
// it returns that channel directly.
func Or(channels ...<-chan any) <-chan any {

	var nonNil []<-chan any
	for _, ch := range channels {
		if ch != nil {
			nonNil = append(nonNil, ch)
		}
	}

	switch len(nonNil) {
	case 0:
		c := make(chan any)
		close(c)
		return c
	case 1:
		return nonNil[0]
	}

	orDone := make(chan any)
	go func() {
		var once sync.Once
		notify := func() {
			once.Do(func() { close(orDone) })
		}
		for _, ch := range nonNil {
			go func(c <-chan any) {
				select {
				case <-c:
					notify()
				case <-orDone:
				}
			}(ch)
		}
	}()

	return orDone

}
