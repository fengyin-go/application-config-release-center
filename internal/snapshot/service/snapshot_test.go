package service_test

import (
	"configcenter/internal/snapshot/cache"
	"configcenter/internal/snapshot/service"
	"configcenter/internal/snapshot/store"
	"configcenter/internal/snapshot/worker"
	"runtime"
	"sync"
	"testing"
)

func TestPublishedSnapshotIsolation(t *testing.T) {
	r := store.New(map[string]string{"mode": "stable"})
	p := service.New(r, &cache.Cache{})
	u := worker.New(r)
	published := p.Publish()
	u.Apply("mode", "canary")
	if published["mode"] != "stable" {
		t.Errorf("published snapshot changed to %q", published["mode"])
	}

	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 2000; i++ {
			service.Checksum(published)
			runtime.Gosched()
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 2000; i++ {
			u.Apply("counter", string(rune(i)))
			runtime.Gosched()
		}
	}()
	close(start)
	wg.Wait()
}
