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

