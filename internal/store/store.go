// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"configcenter/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Application
	CreateApplication(a *model.Application) error
	GetApplication(id string) (*model.Application, error)
	GetApplicationByCode(code string) (*model.Application, error)
	ListApplications() []*model.Application
	UpdateApplication(a *model.Application) error
	DeleteApplication(id string) error

	// Environment
	CreateEnvironment(e *model.Environment) error
	GetEnvironment(id string) (*model.Environment, error)
	ListEnvironments() []*model.Environment
	ListEnvironmentsByAppID(appID string) []*model.Environment
	UpdateEnvironment(e *model.Environment) error
	DeleteEnvironment(id string) error

	// ConfigItem
	CreateConfigItem(c *model.ConfigItem) error
	GetConfigItem(id string) (*model.ConfigItem, error)
	GetConfigItemByAppEnvKey(appID, envID, key string) (*model.ConfigItem, error)
	ListConfigItems() []*model.ConfigItem
	ListConfigItemsByAppID(appID string) []*model.ConfigItem
	ListConfigItemsByEnvID(envID string) []*model.ConfigItem
	UpdateConfigItem(c *model.ConfigItem) error
	DeleteConfigItem(id string) error

	// ConfigVersion
	CreateConfigVersion(c *model.ConfigVersion) error
	GetConfigVersion(id string) (*model.ConfigVersion, error)
	ListConfigVersions() []*model.ConfigVersion
	ListConfigVersionsByConfigItemID(configItemID string) []*model.ConfigVersion

	// Release
	CreateRelease(r *model.Release) error
	GetRelease(id string) (*model.Release, error)
	ListReleases() []*model.Release
	ListReleasesByAppID(appID string) []*model.Release
	ListReleasesByEnvID(envID string) []*model.Release
	UpdateRelease(r *model.Release) error
	DeleteRelease(id string) error

	// AuditLog
	CreateAuditLog(a *model.AuditLog) error
	GetAuditLog(id string) (*model.AuditLog, error)
	ListAuditLogs() []*model.AuditLog
	ListAuditLogsByAction(action string) []*model.AuditLog
	ListAuditLogsByResource(resource string) []*model.AuditLog
}
