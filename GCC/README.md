# GCC相关

## gcc常用参数
1. -v/--version：查看gcc的版本
2. -I：编译的时候指定头文件路径，不然头文件找不到
3. -c：将汇编文件转换成二进制文件，得到.o文件
4. -g：gdb调试的时候需要加
5. -D：编译的时候指定一个宏（调试代码的时候需要使用例如printf函数，但是这种函数太多了对程序性能有影响，因此如果没有宏，则#ifdefine的内容不起作用）
6. -wall：添加警告信息
7. -On：-O是优化代码，n是优化级别：1，2，3

## 静态库制作
- 原材料：源代码（.c或.cpp文件）
- 将.c文件生成.o文件（ex：g++ a.c -c）
- 将.o文件打包
  - ar rcs 静态库名称 原材料(ex: ar rcs libtest.a a.0)


## 动态库制作

- 生成位置无关的目标文件.o，此外加编译器选项-fpic
  - `g++ -fPIC -c unite_time.cpp`
- 生成动态库，加链接器选项-shared
  - g++ -shared -o libunite_time.so unite_time.o


## Linux相关命令

### `top`：显示系统中各个进程的资源占用状况
P：按%CPU使用率排行
M：按%MEM排行
T： 根据时间/累计时间进行排序


### `scp`：不同linux主机之间复制文件
`scp  local_file    remote_username@remote_ip:remote_file `


### `find`：找文件名
`find . -name '[A-Z]*.txt' -print `

### `sar`

从多方面对系统的活动进行报告，包括：文件的读写情况、系统调用的使用情况、磁盘 I/O、CPU 效率、内存使用状况、进程活动及 IPC 有关的活动等

### `df`
用来检查linux服务器的文件系统的磁盘空间占用情况


### `free`
Linux系统中空闲的、已用的物理内存及swap内存,及被内核使用的buffer

### `netstat`

列出所有端口 #netstat -a

列出所有 tcp 端口 #netstat -at

列出所有 udp 端口 #netstat -au

