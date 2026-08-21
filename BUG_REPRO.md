# Bug
后台任务切断请求取消信号，scheduler 和 worker 在取消后仍继续等待与调用。

# 触发
放行一次下游调用后取消请求，先尝试关闭服务，再放行剩余两个任务。

# 错误信息
`shutdown waited on a cancelled request`；`downstream calls grew to 3 after cancellation`。
