## Linux相关命令

1. 查看当前进程？怎么执行退出？怎么查看当前路径？

查看当前进程：ps

执行退出：exit

查看当前路径：pwd

2. ls命令执行什么功能？ 可以带哪些参数，有什么区别？

列出指定目录中的目录，以及文件哪些参数以及区别：a所有文件l详细信息，包括大小字节数，可读可写可执行的权限等

3. 查看文件命令

- vi/vim ：编辑方式查看，可修改
- cat：显示全部文件内容
- more：分页显示文件内容
- less：功能类似more，可以往前翻页
- tail：仅查看尾部，还可以指定行数
- head：仅查看头部,还可以指定行数
- awk：打印符合条件的文本内容

4. 文档编辑
   - 创建目录和移除目录：mkdir rmdir
   - tar：打包
   - grep：在当前目录中中查找包含指定字符串的文件
   - ag：在目录中搜索相应关键字的文件

5. 文件管理
   - which：查找命令在哪个目录下
   - cp：复制文件或目录
   - chmod：控制文件权限
   - scp：不同linux主机之间复制文件
   - tree：列出目录内容
   - whereis：用于程序名的搜索，而且只搜索二进制文件
   - find：找到指定文件名的文件/目录

6. 磁盘管理
   - df：Linux系统上的文件系统的磁盘空间使用情况统计
   - mount：挂载Linux系统外的文件

7. 系统管理
   - groupadd：创建一个新的工作组，新工作组的信息将被添加到系统文件
   - ps：当前进程 (process) 的状态
     - 查看进程CPU占用率和内存占用率：ps -aux
     - 查看进程的父进程id：ps -ef
   - ifconfig：IP 地址配置
   - uname：显示电脑以及操作系统的相关信息。
   - who：当前用户
   - free：内存使用情况
   - top：显示进程的资源占用情况
     -  %CPU、%MEM 、load average(任务队列的平均长度)
   - service：打印/启动/停止指定服务
   - chkconfig：检查，设置系统的各种服务
   - kill：杀掉指定进程
   - fuser：查询文件、目录、socket端口和文件系统的使用进程
   - vmstat：查看CPU性能指令
     - 查看所有CPU核信息：mpstat
     - 每个进程使用CPU的用量分解信息：pidstat

8. 网络
   - lsof：列出当前系统打开文件的信息
     - lsof -i：端口号
     - netstat -npl
   - ifstat：查看网络IO情况的指令
   - netstat： 查看机器已建立的TCP连接的指令
   - ping：测试网络连通性


## 功能实现题：

1.  如何利用linux的指令来查询一个文件的行数/字节数/字符数

wc -l

2. linux下统计一个文件中每个id的出现次数

grep -o "xxx" a | wc -l

3. 如何查看占用cpu最多进程

ps aux|grep -v PID|sort -rn -k +3|head


### `sar`

从多方面对系统的活动进行报告，包括：文件的读写情况、系统调用的使用情况、磁盘 I/O、CPU 效率、内存使用状况、进程活动及 IPC 有关的活动等


### 如何查看主机 CPU 核数？如何查看内存还剩多少

cat /proc/cpuinfo

cat /proc/meminfo

## 查看进程上下文切换

pidstat

