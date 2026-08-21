# Bug
重试 attempt 改变幂等键，旧版本回调和缓存写入又允许成功终态倒退。

# 触发
同一任务执行两次 attempt，第二次成功后让第一次的 running 回调最后到达。

# 错误信息
外部发布执行 2 次，任务与缓存均回退为 `running version=1`。
