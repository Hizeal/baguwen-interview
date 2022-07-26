# 分布式锁


分布式模型下，数据只有一份（或有限制），此时需要利用锁的技术控制某一时刻修改数据的进程数。与单机模式下的锁不同，分布式锁不仅需要保证进程可见，还需要考虑进程与锁之间的网络问题。

单机锁由于其能够共享堆内存，可以将内存作为标记存储位置。而分布式锁需要将标记存储在进程都能看见的地方，比如公共内存。

## 分布式锁的条件

- 互斥性：在任意时刻，对于同一个锁，只有一个客户端能持有
- 具备可重入特性
- 锁失效机制，防止死锁
- 非阻塞锁特性，即没有获取到锁将直接返回获取锁失败

## 基于 redis 的 setnx()、expire() 方法做分布式锁

首先解释函数。

setnx(key，value)，原子操作。如果key不存在，则设置其key，返回1；否则设置失败，返回0。

expire()：设置key的过期时间

具体过程：

1. setnx(lockkey, 1) 如果返回 0，则说明占位失败；如果返回 1，则说明占位成功

2. expire() 命令对 lockkey 设置超时时间，为的是避免死锁问题。
3. 执行完业务代码后，可以通过 delete 命令删除 key

可能存在的问题：

成功执行setnx()后，在expire()执行成功前，发生宕机，此时就还是死锁问题。

如何解决：setnx()、get() 和 getset() 方法来实现分布式锁

简单介绍函数

setset(key,value)：原子操作。对 key 设置 newValue 这个值，并且返回 key 原来的旧值。如果key之前不存在，返回的null。

过程；

1. setnx(lockkey, 当前时间+过期超时时间)，如果返回 1，则获取锁成功；如果返回 0 则没有获取到锁，转向 2
2. get(lockkey) 获取值 oldExpireTime ，并将这个 value 值与当前的系统时间进行比较，如果小于当前系统时间，则认为这个锁已经超时，可以允许别的请求重新获取，转向 3
3. 计算 newExpireTime = 当前时间+过期超时时间，然后 getset(lockkey, newExpireTime) 会返回当前 lockkey 的值currentExpireTime
4. 判断 currentExpireTime 与 oldExpireTime 是否相等，如果相等，说明当前 getset 设置成功，获取到了锁。如果不相等，说明这个锁又被别的请求获取走了，那么当前请求可以直接返回失败，或者继续重试
5. 获取到锁之后，当前线程可以开始自己的业务处理，当处理完毕后，比较自己的处理时间和对于锁设置的超时时间，如果小于锁设置的超时时间，则直接执行 delete 释放锁；如果大于锁设置的超时时间，则不需要再锁进行处理



## ETCD分布式锁

1. Prefix：目录机制，etcd支持前缀查找。具体形式：/锁名称/key的UUID。查询前缀"/锁名称"，返回 Key-Value 列表，同时也包含它们的 Revision，通过 Revision 大小，客户端可以判断自己是否获得锁，如果抢锁失败，则等待锁释放（对应的 Key 被删除或者租约过期），然后再判断自己是否可以获得锁
2. lease：在创建锁的时候绑定租约，并定期进行续约。如果租约到期，锁将被删除；如果获得锁期间持有者因故障不能主动释放锁，则持有的锁也会到期被自动删除，避免了死锁的产生
3. Revision：etcd内部维护了一个全局的Revision值，并会随着事务的递增而递增。可以用Revision值的大小来决定获取锁的先后顺序，实现公平锁。
4. Watch 机制支持监听某个固定的 Key，也支持监听一个范围（前缀机制），当被监听的 Key 或范围发生变化，客户端将收到通知。如果抢锁失败，可通过 Prefix 机制返回的 Key-Value 列表获得 Revision 比自己小且相差最小的 Key（称为 Pre-Key），对 Pre-Key 进行监听，因为只有它释放锁，自己才能获得锁，如果监听到 Pre-Key 的 DELETE 事件，则说明 Pre-Key 已经释放，自己已经持有