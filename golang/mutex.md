# Mutex 

## 面试题

### `mutex`是如何加锁的

分析：也是考察`mutex`的实现原理。基本上就是围绕`自旋-FIFO`来说的。简单理解就是，`mutex`先尝试自旋，自旋不行就所有`goroutine`步入`FIFO`，先到先得。

加锁的这个步骤，讲得非常详细。

答案：`mutex`加锁大概分成两种模式：
1. 在正常模式下，`goroutine`通过自旋来获得锁；
2. 但是如果存在一个`goroutine`等待锁超过了`1ms`，那么`mutex`就会进入饥饿模式，在饥饿模式下，互斥锁会直接交给等待队列最前面的 Goroutine。也就是从等待队列里面唤醒第一个等待者。新的 Goroutine 在该状态下不能获取锁、也不会进入自旋状态，它们只会在队列的末尾等待。如果一个 Goroutine 获得了互斥锁并且它在队列的末尾或者它等待的时间少于 1ms，那么当前的互斥锁就会切换回正常模式

（讨论一下公平性的问题）所以从严格意义上来说，它并不是公平锁，因为在正常状态下，一个新的请求锁的`goroutine`和等待的`goroutine`一起竞争锁。而严格意义的公平应该是永远遵循 `FIFO`。


1. 判断当前 Goroutine 能否进入自旋；

自旋是一种多线程同步机制，当前的进程在进入自旋的过程中会一直保持 CPU 的占用，持续检查某个条件是否为真。在多核的 CPU 上，自旋可以避免 Goroutine 的切换，使用恰当会对性能带来很大的增益，但是使用的不恰当就会拖慢整个程序，所以 Goroutine 进入自旋的条件非常苛刻：

1）互斥锁只有在普通模式才能进入自旋；

2）runtime.sync_runtime_canSpin需要返回 true

- 运行在多 CPU 的机器上；

- 当前 Goroutine 为了获取该锁进入自旋的次数小于四次；

- 当前机器上至少存在一个正在运行的处理器 P 并且处理的运行队列为空

2. 通过自旋等待互斥锁的释放

一旦当前 Goroutine 能够进入自旋就会调用runtime.sync_runtime_doSpin和runtime.procyield执行 30 次的 PAUSE 指令，该指令只会占用 CPU 并消耗 CPU 时间

3. 计算互斥锁的最新状态

处理了自旋相关的特殊逻辑之后，互斥锁会根据上下文计算当前互斥锁最新的状态

4. 更新互斥锁的状态并获取锁

计算了新的互斥锁状态之后，会使用 CAS 函数更新状态

如果没有通过 CAS 获得锁，会调用 runtime.sync_runtime_SemacquireMutex

runtime.sync_runtime_SemacquireMutex 会在方法中不断尝试获取锁并陷入休眠等待信号量的释放，一旦当前 Goroutine 可以获取信号量，它就会立刻返回，sync.Mutex.Lock的剩余代码也会继续执行
#### 类似问题
- 什么是饥饿状态
- 是不是公平锁


## Reference
[mutex is more fair](https://news.ycombinator.com/item?id=15096463)
[mutex fairness]()
[这可能是最容易理解的 Go Mutex 源码剖析](https://segmentfault.com/a/1190000039855697)
