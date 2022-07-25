# docker常用命令

## 介绍一下docker与虚拟机区别

二者目的不同

VMware 为应用提供虚拟的计算机（虚拟机）；Docker 为应用提供虚拟的空间，被称作容器（Container）

对于虚拟机来说
1. 用起来像一台真的机器那样，包括开机、关机，以及各种各样的硬件设备，需要安装操作系统。同时未来多个操作系统高效率同时进行，虚拟机很依赖底层硬件架构提供的虚拟化能力
   1. Hypter-V：代表WSL2
   2. Type-2：代表VMware
2. 一台实体机上的所有的虚拟机实例不会互相影响

而对于容器来说，
1. 容器是直接跑在操作系统之上的，容器内部是应用，应用执行起来就是进程

## 什么是docker

Docker对进程进行封装隔离，属于操作系统层面的虚拟化技术。 由于隔离的进程独立于宿主和其它的隔离的进 程，因此也称其为容器。

1. 镜像
- 特殊的文件系统，除了提供容器运行时所需的程序、库、资源、配置等文件外，还包含了一些为运行时准备的一些配置参数（如匿名卷、环境变量、用户等）
2. 容器
- 容器是镜像运行时的实体
3. 仓库
- 镜像仓库是Docker用来集中存放镜像文件的地方类似于我们之前常用的代码仓库
## images镜像

- docker pull ubuntu:20.04：拉取一个镜像
- docker images：列出本地所有镜像
- docker image rm ubuntu:20.04 或 docker rmi ubuntu:20.04：删除镜像ubuntu:20.04
- docker [container] commit CONTAINER IMAGE_NAME:TAG：创建某个container的镜像
- docker save -o ubuntu_20_04.tar ubuntu:20.04：将镜像ubuntu:20.04导出到本地文件ubuntu_20_04.tar中
- docker load -i ubuntu_20_04.tar：将镜像ubuntu:20.04从本地文件ubuntu_20_04.tar中加载出来

## 容器

- docker [container] create -it ubuntu:20.04：利用镜像ubuntu:20.04创建一个容器。
- docker ps -a：查看本地的所有容器
- docker [container] start CONTAINER：启动容器
- docker [container] stop CONTAINER：停止容器
- docker [container] restart CONTAINER：重启容器
- docker [contaienr] run -itd ubuntu:20.04：创建并启动一个容器


images和container是 类和对象 的关系