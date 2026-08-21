# Bug
临时失败的发布尝试也计入成功事件并提前写成功审计，重试后产生重复结果。

# 触发
让第一次 publisher 调用返回临时错误，第二次重试成功。

# 错误信息
`commits=1 success events=2`；`audit entries=[success failed success]`。
