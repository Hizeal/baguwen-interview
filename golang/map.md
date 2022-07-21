# MAP

使用哈希表作为底层实现，一个哈希表里可以有多个哈希表节点，也即bucket，而每个bucket就保存了map中的一个或一组键值对

## 如何实现顺序读取

Go中map如果要实现顺序读取的话，可以先把map中的key，通过sort包排序
## 哈希冲突

当有两个或以上数量的键被哈希到了同一个bucket时，我们称这些键发生了冲突。Go使用链地址法来解决键冲突

由于每个bucket可以存放8个键值对，所以同一个bucket存放超过8个键值对时就会再创建一个键值对，用类似链表的方式将bucket连接起来。

- 为什么是8个：bucket的数据结构中，仅能存储哈希值的高8位

## 哈希扩容

扩容条件：

1. 负载因子>6.5

   ```go
   负载因子 = 键数量/bucket数量
   ```

2. overflow数量 > `2^15`

   bucket数据结构指示下一个bucket的指针称为overflow bucket，意为当前bucket盛不下而溢出的部分

### 增量扩容

当负载因子过大时，就新建一个bucket，新的bucket长度是原来的2倍，然后旧bucket数据搬迁到新的bucket。

考虑到如果map存储了数以亿计的key-value，一次性搬迁将会造成比较大的延时，Go采用逐步搬迁策略，即每次访问map时都会触发一次搬迁，每次搬迁2个键值对

### 等量扩容

buckets数量不变，重新做一遍类似增量扩容的搬迁动作，把松散的键值对重新排列一次，以使bucket的使用率更高，进而保证更快的存取。


### 为什么并发不安全
map 在扩缩容时，需要进行数据迁移，迁移的过程并没有采用锁机制防止并发操作，而是会对某个标识位标记为 1，表示此时正在迁移数据。如果有其他 goroutine 对 map 也进行写操作，当它检测到标识位为 1 时，将会直接 panic


并发安全则使用`sync.map`

如何实现并发安全：

sync.map底层数据结构：
```go
type Map struct {
 mu Mutex //互斥锁
 read atomic.Value // readOnly
 dirty map[interface{}]*entry //读写数据，原生map，需要加锁保证数据安全
 misses int //多少次读取read未命中
}
```

当读数据：
1. 先查看 read 中是否包含所需的元素
   1. 有，则通过 atomic 原子操作读取数据并返回
   2. 无，则会判断 read.readOnly 中的 amended 属性，他会告诉程序 dirty 是否包含 read.readOnly.m 中没有的数据。如果true，从dirty读取


因此读操作性能高

当写入：
1. 查 read：存在且未被标记删除状态，尝试存储；read 上没有，或者已标记删除状态，紧接下面
2. 上互斥锁
3. 操作dirty