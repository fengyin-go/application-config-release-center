package exporter

type Exporter struct{ pending map[string][]byte }

func New() *Exporter                             { return &Exporter{pending: map[string][]byte{}} }
func (e *Exporter) Queue(id string, data []byte) { e.pending[id] = data }
func (e *Exporter) Flush(id string) []byte       { return e.pending[id] }
