
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