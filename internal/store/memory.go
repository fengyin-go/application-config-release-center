package store

import (
	"sync"

	"configcenter/internal/model"
)

type MemoryStore struct {
	mu              sync.RWMutex
	applications    map[string]*model.Application
	environments    map[string]*model.Environment
	configItems     map[string]*model.ConfigItem
	configVersions  map[string]*model.ConfigVersion
	releases        map[string]*model.Release
	auditLogs       map[string]*model.AuditLog
	configItemVerSeq map[string]int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		applications:     make(map[string]*model.Application),
		environments:     make(map[string]*model.Environment),
		configItems:      make(map[string]*model.ConfigItem),
		configVersions:   make(map[string]*model.ConfigVersion),
		releases:         make(map[string]*model.Release),
		auditLogs:        make(map[string]*model.AuditLog),
		configItemVerSeq: make(map[string]int),
	}
}

var _ Store = (*MemoryStore)(nil)
