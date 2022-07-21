# STL
## RAII
直译过来是“资源获取即初始化”，也就是说在构造函数中申请分配资源，在析构函数中释放资源

## 迭代器与指针区别
迭代器实际上是对“遍历容器”这一操作进行了封装，迭代器实际上是对“遍历容器”这一操作进行了封装

## 容器中迭代器失效

1. 当容器调用erase()方法后，当前位置到容器末尾元素的所有迭代器全部失效。
2. 当容器调用insert()方法后，当前位置到容器末尾元素的所有迭代器全部失效。
3. 如果容器扩容，在其他地方重新又开辟了一块内存。原来容器底层的内存上所保存的迭代器全都失效了。

- 对于序列式容器（如vector，deque，list等），删除当前的iterator会使后面所有元素的iterator都失效
   1. 原因：vector，deque使用了连续分配的内存，删除一个元素导致后面所有的元素会向前移动一个位置。不过erase方法可以返回下一个有效的iterator
   2. 当插入一个元素到vector中，由于引起了内存重新分配，所以指向原内存的迭代器全部失效
- 对于链表式容器和关联式容器
   1. 原因：链表的插入和删除节点不会对其他节点造成影响，因此只会使得当前的iterator失效
   2. 解决办法：利用erase可以返回下一个有效的iteratord的特性，或者直接iterator++


## 为何关联容器的插入删除效率一般比用其他序列容器高

关联容器一般指map,multimap,set,multiset这四种底层实现都是红黑树。对于关联容器来说，存储的只是节点。插入删除只是节点指针的换来换去，不需要做内存拷贝和内存移动
## 哈希表实现，如何解决哈希值冲突
- 解决方法
  - 线性探测
    - 使用hash函数计算出的位置如果已经有元素占用了，则向后依次寻找，找到表尾则回到表头，直到找到一个空位
  - 开链
    - 每个表格维护一个list，如果hash函数计算出的格子相同，则按顺序存在这个list中
  - 再散列
    - 发生冲突时使用另一种hash函数再计算一个地址，直到不冲突

## vector扩容为什么2倍
- 采用成倍方式扩容，可以保证常数的时间复杂度
- 增加指定大小容量，仅达到O(n)时间复杂度
## map、set如何实现，红黑树如何同时实现两种容器，为什么使用红黑树
- map与set底层都是以红黑树的结构实现，因此插入删除等操作都在O(logn）时间内完成，因此可以完成高效的插入删除

- map和set要求是自动排序的，红黑树能够实现这一功能，而且时间复杂度比较低
- 树高度越小越好，BST这种有特殊情况，比如只有左子树有值，导致O(n)复杂度
## unordered_map与map区别
- unordered_map不会根据key的大小进行排序，后者按照键值排序
- map中的元素是按照二叉搜索树存储，进行中序遍历会得到有序遍历
- unordered_map的底层实现是hash_table，后者红黑树
- map适用于有序数据的应用场景，unordered_map适用于高效查询的应用场景
## STL的heap实现
- binary heap本质是一种complete binary tree（完全二叉树），整棵binary tree除了最底层的叶节点之外，都是填满的，但是叶节点从左到右不会出现空隙

大端堆建立
```c++
class Solution {
public:
    void maxTrep(vector<int>& nums,int l,int treapsize){ //将最大元素调整到堆顶
        int largest = l;
        if(l*2+1 < treapsize && nums[l*2+1] > nums[l]){ //完全二叉树性质
            largest = l*2+1;
        }
        if(l*2+2 < treapsize && nums[l*2+2] > nums[l]){
            largest = l*2+2;
        }

        if(largest != l){
            swap(nums[largest],nums[l]);
            maxTrep(nums,largest,treapsize);
        }
    }

    void buildTreap(vector<int>& nums,int treapsize){
        for(int i=treapsize/2;i>=0;i--){
            maxTrep(nums,i,treapsize);
        }
    }

    int findKthLargest(vector<int>& nums, int k) {
        int treapsize = nums.size();
        buildTreap(nums,treapsize);//建堆
```
## 红黑树概念
- 首先是一个二叉排序树
  - 若左子树不空，则左子树上所有结点的值均小于或等于它的根结点的值
  - 若右子树不空，则右子树上所有结点的值均大于或等于它的根结点的值。
  - 左、右子树也分别为二叉排序树
- 满足以下请求
  - 树中所有节点非红即黑
  - 根节点必为黑节点
  - 红节点的子节点必为黑（黑节点子节点可为黑）
  - 从根到NULL的任何路径上黑结点数相同
  - 查找时间一定可以控制在O(logn)
## 常见容器
- **顺序容器**：vector、deque、list
  - vector 采用一维数组实现，元素在内存连续存放，不同操作的时间复杂度为： 插入: O(N) 查看: O(1) 删除: O(N)
  - deque 采用双向队列实现，元素在内存连续存放，不同操作的时间复杂度为： 插入: O(N) 查看: O(1) 删除: O(N)
  - list 采用双向链表实现，元素存放在堆中，不同操作的时间复杂度为： 插入: O(1) 查看: O(N) 删除: O(1)
- **关联式容器**：set、multiset、map、multimap
  - 均采用红黑树实现，红黑树是平衡二叉树的一种。不同操作的时间复杂度近似为: 插入: O(logN) 查看: O(logN) 删除: O(logN)
- unordered_map、unordered_set、unordered_multimap、 unordered_multiset 
  - 上述四种容器采用哈希表实现，不同操作的时间复杂度为： 插入: O(1)，最坏情况O(N) 查看: O(1)，最坏情况O(N) 删除: O(1)
- **容器适配器**：栈、队列、优先级队列

## deque底层存储

存储数据的空间是由一段一段等长的连续空间构成，各段空间之间并不一定是连续的，可以位于在内存的不同区域。

当 deque 容器需要在头部或尾部增加存储空间时，它会申请一段新的连续空间，同时在 map 数组的开头或结尾添加指向该空间的指针，由此该空间就串接到了 deque 容器的头部或尾部。