package service

import (
	"configcenter/internal/importflow/cache"
	"configcenter/internal/importflow/exporter"
	"configcenter/internal/importflow/parser"
)

type Importer struct {
	Parser   *parser.Parser
	Cache    *cache.Cache
	Exporter *exporter.Exporter
}

func (i *Importer) Submit(id, body string) {
	parsed := i.Parser.Parse(body)
	i.Cache.Save(id, parsed)
	i.Exporter.Queue(id, parsed)
}
