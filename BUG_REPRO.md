# Bug
请求 Scope 归池时未清理，延迟审计仍持有该对象，复用后发生跨租户污染。

# 触发
结束 tenant-a 请求并延迟审计，随后让 tenant-b 复用相同 Scope，再刷新审计。

# 错误信息
`tenant-b inherited labels [secret-a clean-b]`；tenant-a 审计变为 tenant-b。
