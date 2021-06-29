include .env # 引入环境变量,通过$(key)读取属性值

# 声明 ‘伪目标’即`.PHONY` 之后，make就不会去检查是否存在一个叫做 <target> 的文件，而是每次运行都执行对应的命令


# 🔨 准备阶段
.PHONY: env
# print variable within .env
env:
	@echo $(DOMAIN)
