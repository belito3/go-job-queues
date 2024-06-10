package handle

import (
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newTestServer(t *testing.T) {
	var sent, processed uint64
	wg := sync.WaitGroup{}

	process := func(w *Worker, p Job) error {
		atomic.AddUint64(&processed, 1)
		wg.Done()
		return nil
	}
	print(process)

	dispatcher := NewDispatcher(2, 10)
	dispatcher.Run()

	spam := func() {
		for i := 0; i < 5; i++ {
			atomic.AddUint64(&sent, 1)

			work := Job{Payload{"1"}}
			dispatcher.DispatchJob(&work)
			time.Sleep(time.Microsecond)
		}
	}

	wg.Add(10)
	go spam()
	go spam()

	wg.Wait()
	dispatcher.Stop()

	if sent != processed || sent != 10 {
		t.Errorf("not all jobs were processed %v jobs %v processed",
			sent, processed)
	}
}

func TestMain(m *testing.M) {
	var t *testing.T
	newTestServer(t)
	os.Exit(m.Run())
}
