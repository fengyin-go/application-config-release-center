# Bug
生产者、消费者、协调器与收集器的退出所有权错位，错误或取消后管线无法一致结束。

# 触发
批次先发送一条正常配置，再发送错误；同时用 gate 控制 worker 启动并取消收集。

# 错误信息
`producer error left consumer and collector blocked`；协调器提前完成；取消无法唤醒收集端。
