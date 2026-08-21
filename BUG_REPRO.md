# Bug
导入解析器、缓存和导出队列共享复用缓冲区，后一次导入覆盖前一次内容。

# 触发
先提交 `alpha` 批次，再提交同长度的 `bravo` 批次，随后读取第一批缓存和导出。

# 错误信息
`first cached import became "bravo"`；`first exported import became "bravo"`。
