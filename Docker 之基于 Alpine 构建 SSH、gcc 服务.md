# Docker 之基于 Alpine 构建 SSH、gcc 服务

## 一、参考资源

> - [不使用 gcc images 原因 - 太大，不支持 ssh](https://www.jianshu.com/p/30fa6448bb63)
> - [alpine linux images](https://blog.csdn.net/qq_39629343/article/details/81513110)

> - [SSH 服务](https://blog.csdn.net/AMimiDou_212/article/details/106502393)
> - [alpinelinux openssh](https://wiki.alpinelinux.org/wiki/Setting_up_a_ssh-server)
> - [ssh 配置参考，但应该在镜像 build 阶段配置](https://blog.csdn.net/qq_39626154/article/details/82856865)
> - [docker hub alpine-sshd](https://hub.docker.com/r/danielguerra/alpine-sshd/dockerfile/)
> - [demo1: ssh-keygen -A](https://blog.csdn.net/highlevels/article/details/99626949)
> - [demo1 github](https://github.com/tfwcn/docker-alpine-server)
> - [demo2: ssh-keygen -t x -P x -f x](https://blog.csdn.net/AMimiDou_212/article/details/106502393)

<br />
<br />

## 二、Docker

### 2.1 设置下载源

```bash
### 在/etc/apk/repositories插入官方链接
echo "https://dl-cdn.alpinelinux.org/alpine/v3.11/main" > /etc/apk/repositories
echo "https://dl-cdn.alpinelinux.org/alpine/v3.11/community" >> /etc/apk/repositories
### 官方链接替换成国内TUNA源
sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories


### 直接插入国内TUNA源
echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.8/main" > /etc/apk/repositories
echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.8/community" >> /etc/apk/repositories
```

### 2.2 安装 bash

> http://www.joycode.com.cn/archives/768

### 2.3 gcc 环境

```bash
# 在 alpine 之中，libc 是要单独安装的
apk add libc-dev

# 安装 gcc 编译器
apk add gcc
```

### 2.4 设置时区

> https://www.cnblogs.com/huoqs/p/10975375.html

ln 命令：创建文件链接

文件权限就被修改为 777（可读可写可执行）
&& chmod 777 /root/start.sh \

<br />
<br />

## 三、Dockerfile

### 3.1 Command
<!-- docker run --name gcc-sys -p 7777:22 -it alpine /bin/sh -->

```bash
# 插入下个代码块的内容
vi Dockerfile
# 根据Dockerfile构建自定义镜像 (path为`.`时为相对路径，读取当前路径下的Dockerfile文件)
docker build -t <name>:<tag> <path>
# 根据自定义镜像创建容器
docker run -d --name <container name> -p <本机端口,需打开防火墙>:22 <image name>:<image tag>
# 若为运行可手动运行
docker start <container name>
# 进入容器 (此方式进入无验证)
docker exec -it <container name> /bin/sh

# 拷贝文件并执行
docker cp /tmp/docker-cp/test.c <container name>:/tmp/test.c
gcc test.c -o a.out
./a.out
```

### 3.2 ssh-keygen -A

```Dockerfile
FROM alpine:latest

# env
ENV TZ=Asia/Shanghai
ENV PASSWORD=210713

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.9/main/" > /etc/apk/repositories \
	&& echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.9/community/" >> /etc/apk/repositories \
	&& apk update \
	&& apk upgrade \
	&& apk add --no-cache bash bash-doc bash-completion \
	&& apk add --no-cache openssh tzdata libc-dev gcc \
	&& ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
  && echo $TZ > /etc/timezone \
	&& sed -i 's/#Port 22/Port 22/g' /etc/ssh/sshd_config \
	&& sed -i 's/#ListenAddress 0.0.0.0/ListenAddress 0.0.0.0/g' /etc/ssh/sshd_config \
	&& sed -i 's/#ListenAddress ::/ListenAddress ::/g' /etc/ssh/sshd_config \
	&& sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/g' /etc/ssh/sshd_config \
	&& sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/g' /etc/ssh/sshd_config \
	&& ssh-keygen -A \
	&& echo root:$PASSWORD | chpasswd \
	&& rm -rf /var/cache/apk/* /tmp/*

# 开放端口
EXPOSE 22

# 容器启动时执行ssh启动命令
CMD ["/usr/sbin/sshd", "-D"]
```

### 3.3 四种

```Dockerfile
FROM alpine:latest

# env
ENV TZ=Asia/Shanghai
ENV PASSWORD=210713

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.9/main/" > /etc/apk/repositories \
	&& echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.9/community/" >> /etc/apk/repositories \
	&& apk update \
	&& apk upgrade \
	&& apk add --no-cache bash bash-doc bash-completion \
	&& apk add --no-cache openssh tzdata libc-dev gcc \
	&& ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
  && echo $TZ > /etc/timezone \
	&& sed -i 's/#Port 22/Port 22/g' /etc/ssh/sshd_config \
	&& sed -i 's/#ListenAddress 0.0.0.0/ListenAddress 0.0.0.0/g' /etc/ssh/sshd_config \
	&& sed -i 's/#ListenAddress ::/ListenAddress ::/g' /etc/ssh/sshd_config \
	&& sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/g' /etc/ssh/sshd_config \
	&& sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/g' /etc/ssh/sshd_config \
	&& ssh-keygen -t dsa -P "" -f /etc/ssh/ssh_host_dsa_key \
  && ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key \
  && ssh-keygen -t ecdsa -P "" -f /etc/ssh/ssh_host_ecdsa_key \
  && ssh-keygen -t ed25519 -P "" -f /etc/ssh/ssh_host_ed25519_key \
	&& echo root:$PASSWORD | chpasswd \
	&& rm -rf /var/cache/apk/* /tmp/*

# 开放端口
EXPOSE 22

# 容器启动时执行ssh启动命令
CMD ["/usr/sbin/sshd", "-D"]
```



# 搭建步骤

1. 安装docker：[runoob](https://www.runoob.com/docker/centos-docker-install.html)

2. 新建任意名称文件夹，存放Dockerfile文件
```bash
$ docker pull alpine #安装最小linux镜像
$ docker images # 查看已安装镜像

$ cd /
$ mkdir source/gcc-sys # 存放Dockerfile文件
$ vi source/gcc-sys/Dockerfile
# 拷贝 3.1 或 3.2 中的代码块于这个文件中
```

3. 根据Dockerfile构建自定义镜像 (path为`.`时为相对路径，读取当前路径下的Dockerfile文件)
```bash
docker build -t <name>:<tag> <path>
```

4.  根据自定义镜像创建容器
```bash
docker run -d --name <container name> -p <本机端口,需打开防火墙>:22 -e xxxx=xxxx  <image name>:<image tag>
```

5. 若未运行可手动运行
```bash
docker start <container name>
```

6. 进入容器 (此方式进入无验证)
```bash
docker exec -it <container name> /bin/sh
```

8. 拷贝文件并执行测试
```bash
docker cp /tmp/docker-cp/test.c <container name>:/tmp/test.c
gcc test.c -o a.out
./a.out
```