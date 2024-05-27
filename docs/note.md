## go-zero 操作 mysql

在 `go-zero` 中操作 `mysql` 数据库，先通过 `goctl` 生成 `model` 文件，然后通过 `model` 文件操作数据库

具体的操作方式如下：

1. 通过 `sql` 文件生成 `model` 文件
2. 数据库连接
3. 操作 `model`，修改业务逻辑

### 通过 sql 文件生成 model 文件

先定义 `user.sql` 文件：

```sql
create table user (
  id bigint AUTO_INCREMENT,
  username varchar(36) NOT NULL,
  password varchar(64) default "",
  UNIQUE name_index (username),
  PRIMARY KEY (id)
) ENGINE = InnoDB COLLATE utf8mb4_general_ci;
```

然后通过这个 `user.sql` 生成 `model` 文件

```bash
goctl model mysql ddl --src user.sql --dir .
```

### 数据库连接

在 `etcd/user.yaml` 中配置数据库连接信息

```yaml
Mysql:
  DataSource: root:123456@tcp(go-uccs:3306)/zero_study?charset=utf8mb4&parseTime=True&loc=Local
```

在 `internal/config/config.go` 中增加数据库的配置参数

```go
type Config struct {
  rest.RestConf
  Mysql struct {
    DataSource string
  }
}
```

在 `internal/srv/servicecontext.go` 中增加数据库的连接方式，这里的 `user` 是 `user.sql` 文件生成的 `model`

```go
type ServiceContext struct {
  Config     config.Config
  UsersModel user.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
  mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)
  return &ServiceContext{
    Config:     c,
    UsersModel: user.NewUserModel(mysqlConn),
  }
}
```

### 修改业务逻辑

在 `internal/handler/loginhandler.go` 文件中就可以调用 `model` 中的相关方法完成业务逻辑了

调用 `model` 中的方法是通过 `l.svcCtx.UsersModel` 调用

```go
func (l *LoginLogic) Login(req *types.LoginRequest) (resp string, err error) {
  user, err := l.svcCtx.UsersModel.FindOneByUsernameAndPassword(l.ctx, req.UserName, req.Password)
  if err != nil {
    return "", errors.New("登录失败")
  }
  return user.Username, nil
}
```

## go-zero 接入 gorm

`gorm` 是目前最流行的 `orm` 框架，相比 `go-zero` 再带的 `sqlx`，`gorm` 操作数据库更加方便

`go-zero` 支持接入 `gorm`，具体的操作方式如下：

1. 准备 `gorm` 连接数据库的配置，在 `etc/users.yaml` 中增加 `mysql` 的配置
   ```yaml
   Mysql:
     DataSource: root:123456@tcp(go-uccs-1:3306)/zero_study?charset=utf8mb4&parseTime=True&loc=Local
   ```
2. 准备 `gorm` 连接数据库的方法
   ```go
   func InitGorm(MysqlDataSource string) *gorm.DB {
     db, err := gorm.Open(mysql.Open(MysqlDataSource), &gorm.Config{})
     if err != nil {
       panic("连接 mysql 数据库失败：" + err.Error())
     } else {
       fmt.Println("连接数据库成功")
     }
     db.AutoMigrate(&user.UserModel{})
     return db
   }
   ```
3. 在配置文件中增加 `mysql` 的结构体
   ```go
   type Config struct {
     rest.RestConf
     Mysql struct {
       DataSource string
     }
   }
   ```
4. 在 `go-zero` 中增加 `gorm` 的连接方式
   ```go
   type ServiceContext struct {
     Config config.Config
     DB     *gorm.DB
   }
   func NewServiceContext(c config.Config) *ServiceContext {
     return &ServiceContext{
       Config: c,
       DB:     gorm_conn.InitGorm(c.Mysql.DataSource),
     }
   }
   ```

然后就可以表写相关的业务逻辑了

```go
func (l *LoginLogic) Login(req *types.LoginRequest) (resp string, err error) {
  var user user.UserModel
  err = l.svcCtx.DB.Take(&user, "username = ? and password = ?", req.UserName, req.Password).Error
  if err != nil {
    return "", errors.New("登录失败")
  }
  return "登录成功", nil
}
```

## go-zero 接入 rpc

服务对给外部使用的，一般用 `http` 服务，服务之间的通信可以使用 `rpc` 进行操作

`go-zero` 支持 `rpc` 服务，具体的操作方式如下：

1. 编写 `proto` 文件，定义 `rpc` 的服务，以前讲过，具体可以看：[protocol 和 grpc 的基本使用](https://juejin.cn/post/7254039378195087417)
   ```proto
   syntax = "proto3";
   package user;
   option go_package = "./user";
   message UserInfoRequest{
    uint32 userId = 1;
   }
   message UserInfoResponse{
    uint32 userId = 1;
    string username = 2;
   }
   message UserCreateRequest{
    string username = 1;
    string password = 2;
   }
   message UserCreateResponse{
    string err = 1;
   }
   service user{
    rpc UserInfo(UserInfoRequest) returns(UserInfoResponse);
    rpc UserCreate(UserCreateRequest)returns(UserCreateResponse);
   }
   ```
2. 运行命令，就会生成 `go-zero` 相关的文件
   ```bash
   goctl rpc protoc user.proto --go_out=./types --go-grpc_out=./types --zrpc_out=.
   ```

然后就可以编写业务逻辑了

```go
func (l *UserInfoLogic) UserInfo(in *user.UserInfoRequest) (*user.UserInfoResponse, error) {
  return &user.UserInfoResponse{
    UserId:   in.UserId,
    Username: "uccs",
  }, nil
}
```

### rpc 服务分组

如果一个服务有很多的接口，就需要按照功能进行分组，不然会很难管理

`go-zero` 中如何进行分组呢？

分组还是先从 `proto` 文件开始，我们将 `user` 服务拆成 `UserInfo` 和 `UserAction` 两个服务，代码如下

```proto
service UserInfo{
  rpc UserInfo(UserInfoRequest) returns(UserInfoResponse);
}
service UserAction{
  rpc UserCreate(UserCreateRequest)returns(UserCreateResponse);
}
```

运行命令，只需在命令的最后加上 `-m` 即可：

```bash
goctl rpc protoc user.proto --go_out=./types --go-grpc_out=./types --zrpc_out=. -m
```

我们在生成的代码 `internal/logic` 和 `internal/server` 目录中会看到 `userinfo`和`useraction` 两个目录，说明已经对不同的功能进行了分组`````````````````````````````````````````

## 完整的 api 和 rpc 接入

学完 `go-zero`，知道了怎么接入 `api` 和 `rpc` 现在就写一版完整的代码：

过程主要分为两部分：

1. 创建 `api`
2. 创建 `rpc`

### 创建 api

`api` 部分分为 `4` 个步骤

1. 编写 `user.api` 文件，生成对应的 `go` 文件
2. 增加 `jwt` 相关的配置和代码
3. 增加 `rpc` 相关的配置和代码
4. 编写业务逻辑

#### 编写 user.api 文件

`api` 部分的接口有 `3` 个：

- `/api/user/login`
- `/api/user/create`
- `/api/user/:id`

这三个接口 `/create` 和 `/:id` 都需要验证用户是否登录，所以需要一个支持 `jwt`

在 `go-zero-study/api` 目录下创建 `user.api`

`user.api` 定义：

1. `:id` 形式的参数获取要定义 `json tag` 为 `path:"id"`
2. `prefix` 定义接口的前缀
3. `jwt` 定义接口是否需要进行 `jwt` 验证

```proto
type UserRequest {
  Username string `json:"username"`
  Password string `json:"password"`
}

type UserInfoResponse {
  Username string `json:"username"`
  UserId   string `json:"userId"`
}

type UserInfoRequest {
  id string `path:"id"`
}

@server (
  prefix: /api/user
)
service user {
  @handler login
  post /login (UserRequest) returns (string)
}

@server (
  prefix: /api/user
  jwt:    Auth
)
service user {
  @handler userInfo
  get /:id (UserInfoRequest) returns (UserInfoResponse)

  @handler createUser
  post /create (UserRequest) returns (UserInfoResponse)
}
```

然后在 `user.api` 的目录下运行命令，生成相关文件：

```bash
goctl api go -api user.api -dir .
```

#### 增加 jwt 相关的配置和代码

在 `api/etc/user.yaml` 文件中增加 `jwt` 的配置

```yaml
Name: user
Host: 0.0.0.0
Port: 8888
Auth:
  AccessSecret: 123456789qwert
  AccessExpire: 3600
```

配置 `api/internal/config/config.go` 文件中配置 `Auth` 解析方法

```go
type Config struct {
  rest.RestConf
  Auth struct {
    AccessSecret string
    AccessExpire int64
  }
}
```

写一个公共的 `jwt` 方法，包括生成 `jwt` 和 解析 `jwt`

```go
type JwtPayload struct {
  UserID   uint   `json:"userId"`
  UserName string `json:"username"`
}

type CustomClaims struct {
  JwtPayload
  jwt.RegisteredClaims
}

func GenToken(user JwtPayload, accessSecret string, expires int64) (string, error) {
  claims := CustomClaims{
    JwtPayload: user,
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expires))),
    },
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString([]byte(accessSecret))
}

func ParseToken(tokenStr string, accessSecret string, expires uint32) (*CustomClaims, error) {
  token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(accessSecret), nil
  })
  if err != nil {
    return nil, err
  }
  claims, ok := token.Claims.(*CustomClaims)
  if !ok {
    return nil, errors.New("token 无效")
  }
  return claims, nil
}
```

#### 增加 rpc 相关的配置和代码

在 `api/etc/user.yaml` 文件中增加 `rpc` 的配置

```yaml
Name: user
Host: 0.0.0.0
Port: 8888
Auth:
  AccessSecret: 123456789qwert
  AccessExpire: 3600
UserRpc:
  Etcd:
    Hosts:
      - etcd-server:2379
    Key: user.rpc
```

在 `api/internal/svc/servicecontext.go` 文件中增加 `rpc` 的连接

```go
type ServiceContext struct {
  Config  config.Config
  UserRpc user.User
}
func NewServiceContext(c config.Config) *ServiceContext {
  return &ServiceContext{
    Config:  c,
    UserRpc: user.NewUser(zrpc.MustNewClient(c.UserRpc)),
  }
}
```

#### 编写业务逻辑

在编写业务逻辑之前，需要对响应进行改造

1. 需要进行 `jwt` 验证的接口如果没有传递 `Authorization`，`http` 状态码返回 `200`，自定义 `code` 返回 `401`
   - 在 `api/user.go` 中增加 `JwtUnauthorizedResult` 方法
   ```go
   func JwtUnauthorizedResult(w http.ResponseWriter, r *http.Request, err error) {
     httpx.WriteJson(w, http.StatusOK, response.Body{Code: 401, Data: nil, Msg: err.Error()})
   }
   ```
   - 修改 `main` 文件中的 `rest` 配置
   ```go
   server := rest.MustNewServer(c.RestConf, rest.WithUnauthorizedCallback(JwtUnauthorizedResult))
   ```
2. 统一响应格式，包括 `code` 和 `msg`
   ```go
   type Body struct {
     Code int    `json:"code"`
     Data any    `json:"data"`
     Msg  string `json:"msg"`
   }
   func Response(r *http.Request, w http.ResponseWriter, res any, err error) {
     if err != nil {
       httpx.WriteJson(w, http.StatusOK, &Body{
         Code: 500,
         Data: nil,
         Msg:  err.Error(),
       })
       return
     }
     httpx.WriteJson(w, http.StatusOK, &Body{
       Code: 200,
       Data: res,
       Msg:  "成功",
     })
   }
   ```
   - 修改 `api/internal/handler/loginhandler.go` 文件，每个 `api` 的响应部分都要进行修改
     ```go
     func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) {
         var req types.UserRequest
         if err := httpx.Parse(r, &req); err != nil {
           httpx.ErrorCtx(r.Context(), w, err)
           return
         }
         l := logic.NewLoginLogic(r.Context(), svcCtx)
         resp, err := l.Login(&req)
         // if err != nil {
         // 	httpx.ErrorCtx(r.Context(), w, err)
         // } else {
         // 	httpx.OkJsonCtx(r.Context(), w, resp)
         // }
         // 注释上面的代码，增加下面的代码
         response.Response(r, w, resp, err)
       }
     }
     ```

业务逻辑没啥好说的，就是调用 `rpc` 的方法，然后返回结果

这里记录一下，参数的获取方法

1. 拿取 `rpc` 的方法 `l.svcCtx.UserRpc.xxx`
2. `ctx` 使用 `l.ctx`
3. 获取配置文件的参数 `l.svcCtx.Config.Auth`
4. 解析 `jwt` 字符串上的参数
   ```go
   userId := l.ctx.Value("userId").(json.Number)
   uid, _ := userId.Int64()
   ```

### 创建 rpc

创建 `rpc` 部分为 `3` 个步骤

1. 编写 `user.proto` 文件，生成对应的 `go` 文件
2. 增加 `rpc` 配置
3. 配置 `gorm` 连接数据库
4. 编写业务逻辑

#### 编写 user.proto 文件

在 `rpc/user.proto` 文件中编写 `rpc` 的接口

```proto
syntax = "proto3";
package user;
option go_package = "./user";

message UserRequest{
  string username = 1;
  string password = 2;
}

message UserResponse{
  string username = 1;
  uint32 userId = 2;
}

message UserInfoRequest{
  uint32 userId = 1;
}

service user {
  rpc UserInfo(UserInfoRequest)returns(UserResponse);
  rpc Login (UserRequest) returns(UserResponse);
  rpc UserCreate(UserRequest)returns(UserResponse);
}
```

然后在 `user.proto` 的目录下运行命令，生成相关文件：

```bash
goctl rpc protoc user.proto --go_out=./types --go-grpc_out=./types --zrpc_out=. -m
```

#### 增加 rpc 相关的配置和代码

在 `rpc/etc/user.yaml` 文件中增加 `rpc` 的配置

```yaml
Name: user.rpc
ListenOn: 0.0.0.0:8080
Etcd:
  Hosts:
    - etcd-server:2379
  Key: user.rpc
```

其他的啥也不用管，`go-zero` 会自动帮我们完成

#### 增加 gorm 连接数据库

在 `rpc/etc/user.yaml` 文件中增加数据的连接配置

```yaml
Name: user.rpc
ListenOn: 0.0.0.0:8080
Etcd:
  Hosts:
    - etcd-server:2379
  Key: user.rpc
Mysql:
  DataSource: root:123456@tcp(go-uccs:3306)/zero_study?charset=utf8mb4&parseTime=True&loc=Local
```

在根目录下创建一个 `model` 的文件夹，用来管理数据库 `model` 文件

然后在里面增加一个 `user.go` 的文件，用来管理 `user` 的 `model`

```go
type User struct {
  gorm.Model
  ID       uint32 `gorm:"index,primarykey" json:"id"`
  Username string `gorm:"size:32" json:"username"`
  Password string `gorm:"size:64" json:"password"`
}
```

写一个公共的方法，用来连接数据库

```go
func GormConn(MysqlDataSource string) *gorm.DB {
  db, err := gorm.Open(mysql.Open(MysqlDataSource), &gorm.Config{})
  if err != nil {
    panic("数据库连接失败，" + err.Error())
  }
  fmt.Println("数据库连接成功")
  db.AutoMigrate(&model.User{})
  return db
}
```

修改 `/rpc/internal/svc/servicecontext.go`，增加 `gorm` 的连接

```go
type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DB:     gorm_conn.GormConn(c.Mysql.DataSource),
	}
}
```

#### 编写业务逻辑

调用 `model` 中的方法，完成业务逻辑

```go
db := l.svcCtx.DB
u := &model.User{}
err := db.Take(u, "username = ? and password = ?", in.Username, in.Password).Error
```
