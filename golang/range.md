# Range

range 是 Go 语言用来遍历的一种方式，它可以操作数组、切片、map、channel 等

## 实现原理
###  range for slice
遍历 slice 前会先获取 slice 的长度 len_temp 来作为循环次数

每次循环会先获取元素值，如果接收index和value，对其赋值，发生一次数据拷贝

- 循环过程中新添加的元素无法遍历到

### range for map
遍历 map 时没有指定循环次数，循环体与遍历 slice 类似

- 由于map插入元素随机，故新加入的元素不能保证遍历到

### range for channel

channel 遍历是依次从 channel 中读取数据
- 没有元素，阻塞等待
- channel关闭，则会解除阻塞并退出循环