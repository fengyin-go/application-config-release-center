# Bug
批量发布的资源延迟到函数结束才释放，第三项未处理时错误又被提交结果覆盖。

# 触发
资源上限设为 2，依次发布 dev、staging、prod 三套环境。

# 错误信息
`success response covered only 2 of 3 environments`，审计显示 success 且事务已提交。
