include .env # 引入环境变量,通过$(key)读取属性值

# 在定义时扩展(静态扩展),非动态
postgres_database := "postgresql://$(POSTGRES_USER):$(POSTGRES_PWD)@$(DOMAIN):$(POSTGRES_PORT)/$(DATABASE_NAME)?sslmode=disable"
mysql_database := "mysql://$(MYSQL_USER):$(MYSQL_PWD)@tcp($(DOMAIN):$(POSTGRES_PORT))/$(DATABASE_NAME)"

# 声明 ‘伪目标’即`.PHONY` 之后，make就不会去检查是否存在一个叫做 <target> 的文件，而是每次运行都执行对应的命令


# 🔨 准备阶段
.PHONY: env postgres mysql createdb dropdb
# print variable within .env
env:
	@echo $(postgres_database)
# 终端执行"docker run"的所有记录
history:
	@history | grep "docker run"
# docker 运行 postgres 镜像
postgres:
	docker run --name postgres12 --network bank-network -p $(POSTGRES_PORT):$(POSTGRES_PORT) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PWD) -d postgres:12-alpine
# docker 运行 mysql 镜像
mysql:
	@docker run --name mysql8 -p $(MYSQL_PORT):$(MYSQL_PORT)  -e MYSQL_ROOT_PASSWORD=$(MYSQL_PWD) -d mysql:8
# Postgres容器在本地设置为信任身份验证，无需密码
createdb:
	@docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	@docker exec -it postgres12 dropdb simple_bank


# 🔨 数据库迁移
.PHONY: migrateup migrateup1 migratedown migratedown1
# 开始迁移（all version）
migrateup:
	migrate -path db/migration -database $(postgres_database) -verbose up
# 仅迁移version 1的schema
migrateup1:
	migrate -path db/migration -database $(postgres_database) -verbose up 1
# 清空迁移
migratedown:
	migrate -path db/migration -database $(postgres_database) -verbose down
# version 1
migratedown1:
	migrate -path db/migration -database $(postgres_database) -verbose down 1
