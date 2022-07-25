
# 用户登录

main函数->路由->调用对应中间件加密密码->调用对应Handler

## Handler层：UserLoginHandler(c *gin.Context)

- 传入context，解析其中的name和password，解析失败，响应体写入错误
- 调用Service层QueryUserLogin，将上文解析的name和password送入
- Service返回包含user_id和Token的响应
- Handler在响应体写入状态码、状态信息、id和Token
  - 状态码1表示错误
  - 状态码0表示正常返回


## Service层：QueryUserLogin(username, password string) (*LoginResponse, error)

- 定义结构LoginResponse，包含user_id和Token
- 检测username和password是否合法
- 调用models层QueryUserLogin查找是否存在对应用户
- 若存在，颁发Token，打包数据写入到结构，返回给Handler

## Models层：QueryUserLogin(q.username, q.password, &login)
- 定义用户登录表
  - 用户id：外键，UserInfo表主键
  - UserInfoId:
  - Username：登陆表主键
  - Password

- 通过Gorm接口在MySQL中查找是否存在匹配该用户名和对应用户密码，返回错误信息


# 用户注册

main函数->路由->调用对应中间件加密密码->调用对应Handler

## Handles：UserRegisterHandler(c *gin.Context)

- 传入context，解析username和password，解析失败，响应体写入错误
- 调用Service层PostUserLogin(username, password string)
- Service返回结构体LoginResponse，包含user_id和Token的响应
- Handler在响应体写入状态码、状态信息、id和Token

## Service层：PostUserLogin(username, password string) (*LoginResponse, error)
- 检测传入username和password是否合法，构建用户信息表
- 调用models层AddUserInfo，添加该用户信息表
  - 用户登录表
  - 用户信息表UserInfo
- 添加成功，颁发Token，打包数据，返回给上层Handler

## Models层：AddUserInfo(&userinfo)
- 传入用户信息表
- 通过Gorm给MySQL添加用户信息
- 若出错则返回错误，否则不返回


# 用户发布视频 

main函数->router->调用对应中间件解析Token->调用对应Handler

## Handlers：PublishVideoHandler(c *gin.Context)
- 解析user_id和视频title，
- 提取文件后缀名，判断是否为视频文件格式；支持多文件上传
- 根据user_id与用户发布视频数量，构建唯一文件名，与文件拓展名结合，存于static文件下
- 截取一帧作为视频封面
- 调用Service层PostVideo(userId int64, videoName, coverName, title string)
- Handler在响应体写入状态码、状态信息

## Service：PostVideo(userId int64, videoName, coverName, title string) error 
- 准备参数，视频名和封面
- 构建视频表
- 调用models的NewVideoDAO().AddVideo(video)，添加视频表

## Models层：NewVideoDAO().AddVideo(video)
- 通过Gorm接口在MySQL中添加视频表

由于视频和userinfo有多对一的关系，所以传入的Video参数一定要进行id的映射处理

# 获得用户发布视频列表

main函数->router->调用对应中间件取出user_id->调用对应Handler

## Handler：QueryVideoListHandler(c *gin.Context) 
- 取出context中的user_id
- 调用Service层的QueryVideoListByUserId(userId),返回该id的视频列表
- 视频列表写入响应体中

## Service：
- 检测数据库中UserInfo表是否存在该id
- 调用Model层的QueryVideoListByUserId(q.userId, &q.videos)，视频列表存于q.videos
- 调用Model层的QueryUserInfoById(q.userId, &userInfo)，用户信息表中对应的用户信息存于userInfo
- 填充视频信息表的Auther字段，封装视频列表

## Model：
- 分别获取user_id对应的视频列表和用户信息表

# 用户点赞操作

main函数->router->调用对应中间件解析Token->调用对应Handler

## Handler：PostFavorHandler(c *gin.Context) 
- 取出context的user_id、video_id、动作类型（点赞1或取消点赞2）
- 调用Service的PostFavorState(p.userId, p.videoId, p.actionType)
- 操作成功或失败的状态码和错误消息写入到响应体中

## Service：PostFavorState(userId, videoId, actionType int64) error
- 传入的参数，封装为一个结构，定义其方法，用以调用Model层函数
- 检查动作类型：既非1也非2，不支持该操作
- 如果是点赞操作，调用Model层的PlusOneFavorByUserIdAndVideoId(p.userId, p.videoId)
  - 更新对应的用户点赞状态，该信息存于Redis中，key为(favor, userId)，value为视频id

- 如果是取消点赞操作，调用Model层的MinusOneFavorByUserIdAndVideoId(p.userId, p.videoId)
  - 更新对应用户点赞状态，同上

- 返回错误信息

## Model层：
- PlusOneFavorByUserIdAndVideoId(p.userId, p.videoId)
更新MySQL中video信息表中video_id对应的favorite_count，user_favorite_videos中间表的user_id及video_id
- MinusOneFavorByUserIdAndVideoId(p.userId, p.videoId)
先判断Video信息表中video_id对应的favorite_count大于，接下来才能减1
在user_favorite_videos中间表中移除user_id和video_id


## 用户评论

main函数->router->调用对应中间件解析Token->调用对应Handler

### Handler层：PostCommentHandler(c *gin.Context)

- 定义结构体，存储Context外还保存Context中的videoId、userId、commentId、actionType、commentText
- 从context解析参数，存储结构体
- 调用Service层的PostComment(p.userId, p.videoId, p.commentId, p.actionType, p.commentText)
- 若Service层处理失败，响应体写入错误代码和错误消息；否则返回操作成功代码和Comment表信息

### Service层：PostComment(userId int64, videoId int64, commentId int64, actionType int64, commentText string) (*Response, error)

- 调用Model层IsUserExistById(p.userId)和IsVideoExistById(p.videoId)，判断数据库是否存在对应信息
- 判断action是发布评论1还是删除评论2
- 如果是发布评论，调用CreateComment()函数
  - 该函数调用Model层的AddCommentAndUpdateCount(&comment)，这里comment为Comment表
- 如果是删除评论，调用DeleteComment()函数
  - 调用DeleteCommentAndUpdateCountById(p.commentId, p.videoId)


### Model层：
#### AddCommentAndUpdateCount(comment *Comment) error
1. 执行事务，数据库中添加该信息；返回任何错误，回滚事务
2. 增加Comment表中的comment_count
3. 执行事务成功，提交事务，返回nil

#### DeleteCommentAndUpdateCountById(commentId, videoId int64)
1. 执行事务，删除表中对应commentId的信息
2. 减少comment_count
3. 提交事务，返回nil