# Bug
请求 context 在入口、store、worker 和 client 四层断开，取消无法隔离或向下游传播。

# 触发
先让请求超时并保存已取消 context，再发 fresh 请求；另用已取消请求触发三次重试。

# 错误信息
入口 deadline 丢失，fresh 请求继承 `context canceled`，取消后仍发生 3 次下游尝试。
