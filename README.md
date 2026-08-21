# ConfigCenter

纯 Go 标准库实现的配置中心后端服务。

## 运行

```bash
cd origin
go build ./cmd/server
./server
```

环境变量：
- `PORT`：监听端口，默认 8080
- `ADDR`：监听地址，默认 ":8080"
- `MAX_PAGE_SIZE`：最大分页大小，默认 100
- `LOG_LEVEL`：日志级别，默认 info

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST   | /api/applications            | 创建应用 |
| GET    | /api/applications            | 应用列表 |
| GET    | /api/applications/{id}       | 应用详情 |
| PUT    | /api/applications/{id}       | 更新应用 |
| DELETE | /api/applications/{id}       | 删除应用 |
| POST   | /api/environments            | 创建环境 |
| GET    | /api/environments            | 环境列表 |
| GET    | /api/environments/{id}       | 环境详情 |
| PUT    | /api/environments/{id}       | 更新环境 |
| DELETE | /api/environments/{id}       | 删除环境 |
| POST   | /api/config-items            | 创建配置项 |
| GET    | /api/config-items            | 配置项列表 |
| GET    | /api/config-items/{id}       | 配置项详情 |
| PUT    | /api/config-items/{id}       | 更新配置项 |
| DELETE | /api/config-items/{id}       | 删除配置项 |
| GET    | /api/config-versions         | 版本历史列表 |
| POST   | /api/releases                | 创建发布 |
| GET    | /api/releases                | 发布列表 |
| GET    | /api/releases/{id}           | 发布详情 |
| PUT    | /api/releases/{id}/status    | 流转发布状态 |
| DELETE | /api/releases/{id}           | 删除发布 |
| GET    | /api/audit-logs              | 审计日志列表 |
| GET    | /api/stats/config-item-counts| 按应用/环境统计配置项数量 |
| GET    | /api/stats/release-status    | 按状态统计发布次数 |
| GET    | /api/stats/audit-actions     | 按操作类型统计审计日志 |
