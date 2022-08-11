# 架构设计

Middleware：Token解析和颁发，密码加密

Handler：解析得到参数，开始调用下层逻辑

Service：
- 上层需要返回数据信息
  - 检查参数
  - 准备数据
  - 打包数据

- 不需要返回数据信息
  - 检查参数，执行上层指定动作

- Model层
  - 面向于数据库的增删改查，不需要考虑和上层的交互

# 用户登录

main函数->路由->调用对应中间件加密密码->调用对应Handler

## middleware

进入中间件SHAMiddleWare内的函数逻辑，得到password明文加密后再设置password。具体需要调用gin.Context的Set方法设置password。随后调用next()方法继续下层路由

## Handler层：UserLoginHandler(c *gin.Context)

- 传入context，解析其中的name和password，解析失败，响应体写入错误
- 调用Service层QueryUserLogin，将上文解析的name和password送入
- Service调用middleware中颁布Token函数，返回包含user_id和Token的响应
- Handler在响应体写入状态码、状态信息、id和Token（HS256）
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

## middleware

进入中间件SHAMiddleWare内的函数逻辑，得到password明文加密后再设置password。具体需要调用gin.Context的Set方法设置password。随后调用next()方法继续下层路由

## Handles：UserRegisterHandler(c *gin.Context)

- 传入context，解析username和password，解析失败，响应体写入错误
- 调用Service层PostUserLogin(username, password string)
- Service返回结构体LoginResponse，包含user_id和Token的响应
- Handler在响应体写入状态码、状态信息、id和Token

## Service层：PostUserLogin(username, password string) (*LoginResponse, error)
- 检测传入username和password是否合法
  - 用户名是否为空，用户名长度限制，密码是否空
- 构建用户信息表
  - 如果用户名已存在，返回错误
- 调用models层AddUserInfo，数据库中添加该用户信息表
  - 用户登录表
  - 用户信息表UserInfo
- 添加成功，颁发Token，打包数据，返回给上层Handler

## Models层：AddUserInfo(&userinfo)
- 传入用户信息表
- 通过Gorm给MySQL添加用户信息
- 若出错则返回错误，否则不返回


# 用户发布视频 

main函数->router->调用对应中间件解析Token->调用对应Handler

此时处于登录状态，只需要传入用户id和Token，验证Token。根据用户Id生成Token

## middleware

解析token，是否超时或token不正确

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

## middleware

解析token，是否超时或token不正确

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


## 用户关注/取消关注

## middleware

解析token，是否超时或token不正确

## Handler：PostFollowActionHandler(c *gin.Context)
- 用户id，他要关注的id，关注/取消关注的动作
- 调用startAction()，startAction()调用Service层的PostFollowAction(p.userId, p.followId, p.actionType)
- Service返回错误
  - 判断是否是未定义操作或是数据库中找不到用户，自己关注自己
  - 否则就是model层错误，说明是重复键值插入

## Service：PostFollowAction(p.userId, p.followId, p.actionType)
- 调用Model层IsUserExistById，检查要关注的用户id是否在数据库当中
- 判断actionType是否为关注或取消关注
- 关注
  - 调用Models层的AddUserFollow(userId, userToId int64)
  - 更新Redis中信息
- 取消关注
  - 调用Models层的CancelUserFollow(p.userId, p.userToId)
  - 更新Redis中信息

## Model层：

### AddUserFollow(p.userId, p.userToId)

更新用户信息表中对应的关注人数和关注者人数
在用户关系表中插入(用户id，关注id)

更新Redis中
### CancelUserFollow(p.userId, p.userToId)

更新用户信息表中对应的关注人数和关注者人数
在用户关系表中删除(用户id，关注id)


# 用户信息表
```go
type UserInfo struct {
	Id            int64       `json:"id" gorm:"id,omitempty"`
	Name          string      `json:"name" gorm:"name,omitempty"`
	FollowCount   int64       `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount int64       `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow      bool        `json:"is_follow" gorm:"is_follow,omitempty"`
	User          *UserLogin  `json:"-"`                                     //用户与账号密码之间的一对一
	Videos        []*Video    `json:"-"`                                     //用户与投稿视频的一对多
	Follows       []*UserInfo `json:"-" gorm:"many2many:user_relations;"`    //用户之间的多对多
	FavorVideos   []*Video    `json:"-" gorm:"many2many:user_favor_videos;"` //用户与点赞视频之间的多对多
	Comments      []*Comment  `json:"-"`                                     //用户与评论的一对多
}
```

# 用户登录表,与用户信息表一对一
```go
type UserLogin struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:200;notnull"`
}
```
# 视频表

```go
type Video struct {
	Id            int64       `json:"id,omitempty"`
	UserInfoId    int64       `json:"-"`
	Author        UserInfo    `json:"author,omitempty" gorm:"-"` //这里应该是作者对视频的一对多的关系，而不是视频对作者，故gorm不能存他，但json需要返回它
	PlayUrl       string      `json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url,omitempty"`
	FavoriteCount int64       `json:"favorite_count,omitempty"`
	CommentCount  int64       `json:"comment_count,omitempty"`
	IsFavorite    bool        `json:"is_favorite,omitempty"`
	Title         string      `json:"title,omitempty"`
	Users         []*UserInfo `json:"-" gorm:"many2many:user_favor_videos;"`
	Comments      []*Comment  `json:"-"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
}
```

# 评论表
```go
type Comment struct {
	Id         int64     `json:"id"`
	UserInfoId int64     `json:"-"` //用于一对多关系的id
	VideoId    int64     `json:"-"` //一对多，视频对评论
	User       UserInfo  `json:"user" gorm:"-"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"-"`
	CreateDate string    `json:"create_date" gorm:"-"`
}
```

# 中间表

user_infos和videos的多对多关系，创建一张user_favor_videos中间表，然后将该表的字段均设为外键，分别存下user_infos和videos对应行的id。如id为1的用户对id为2的视频点了个赞，那么就把这个1和2存入中间表user_favor_videos即可

user_relations中间表，存下用户id和该用户关注者id

# Redis

key：点赞：用户id
value：用户点赞的视频id集合


key：关注：用户id
value：用户关注的id集合

