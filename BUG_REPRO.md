# Bug
已发布的配置视图与后台注册表共享同一个 map，后续更新会改写旧结果，并发访问会产生数据竞争。

# 触发
发布 `mode=stable`，保留返回值；后台更新为 `mode=canary`，同时读取原返回值。

# 错误信息
`published snapshot changed to "canary"`，并伴随 `WARNING: DATA RACE`。
