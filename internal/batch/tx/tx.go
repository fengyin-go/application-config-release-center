package tx

type Transaction struct{ committed, rolledBack bool }

func (t *Transaction) Commit() error    { t.committed = true; return nil }
func (t *Transaction) Rollback() error  { t.rolledBack = true; return nil }
func (t *Transaction) Committed() bool  { return t.committed }
func (t *Transaction) RolledBack() bool { return t.rolledBack }
