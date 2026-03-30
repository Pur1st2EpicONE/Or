package or

import (
	"testing"
	"time"
)

func sig(after time.Duration) <-chan any {
	c := make(chan any)
	go func() {
		time.Sleep(after)
		close(c)
	}()
	time.Sleep(10 * time.Millisecond)
	return c
}

func assertDuration(t *testing.T, got time.Duration, want time.Duration, tolerance time.Duration) {
	if got < want-tolerance || got > want+tolerance {
		t.Fatalf("expected around %v, got %v", want, got)
	}
}

func TestOr(t *testing.T) {

	start := time.Now()
	<-Or(
		sig(2*time.Second),
		sig(100*time.Millisecond),
		sig(1*time.Second),
	)

	elapsed := time.Since(start)
	assertDuration(t, elapsed, 100*time.Millisecond, 50*time.Millisecond)

}

func TestSingleChannel(t *testing.T) {

	c := sig(300 * time.Millisecond)

	testOr := func(name string, orFunc func(...<-chan any) <-chan any) {
		start := time.Now()
		<-orFunc(c)
		elapsed := time.Since(start)
		t.Logf("%s elapsed: %v", name, elapsed)
		assertDuration(t, elapsed, 300*time.Millisecond, 50*time.Millisecond)
	}

	testOr("or1", Or)

}

func TestNoChannels(t *testing.T) {

	if ch := Or(); ch != nil {
		t.Fatalf("expected nil for or1 with no channels")
	}

}
