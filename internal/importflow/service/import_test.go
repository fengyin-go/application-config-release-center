package service_test

import (
	"configcenter/internal/importflow/cache"
	"configcenter/internal/importflow/exporter"
	"configcenter/internal/importflow/parser"
	"configcenter/internal/importflow/service"
	"testing"
)

func TestQueuedImportOwnsItsPayload(t *testing.T) {
	c := cache.New()
	e := exporter.New()
	i := service.Importer{Parser: &parser.Parser{}, Cache: c, Exporter: e}
	i.Submit("first", "alpha")
	i.Submit("second", "bravo")
	if got := string(c.Get("first")); got != "alpha" {
		t.Errorf("first cached import became %q", got)
	}
	if got := string(e.Flush("first")); got != "alpha" {
		t.Errorf("first exported import became %q", got)
	}
}
