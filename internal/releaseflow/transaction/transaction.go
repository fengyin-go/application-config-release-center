package transaction

type Manager struct{ Commits int }
type Tx struct {
	manager *Manager
	closed  bool
}

func (m *Manager) Begin() *Tx { return &Tx{manager: m} }
func (t *Tx) Commit() {
	if !t.closed {
		t.closed = true
		t.manager.Commits++
	}
}
func (t *Tx) Rollback() { t.closed = true }
