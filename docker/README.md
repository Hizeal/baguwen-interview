# docker常用命令

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