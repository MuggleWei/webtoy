# webtoy-gate
后端网关服务

## 通用消息
每个普通回复消息都包含一个外层结构
```
{
    "code": int,  // 错误码, 为0或不存在时表示正确
    "msg": string, // 错误消息, 仅当code不为0时需要检查
    "data": object,  // 附带的回复消息
}
```

每个前端请求应包含的头部信息
| 字段 | 意义 |
| ---- | ---- |
| uid | 用户id |
| session | 会话id |
| token | 会话令牌 |

## 接口

### 用户会话检测
<b>url</b>: /api/v1/user/check
<b>desc</b>: 检测用户会话是否可用
<b>note</b>: 无
<b>payload</b>
```
无
```
<b>return</b>
```
仅包含通用结构
```

### 用户登录
<b>url</b>: /api/v1/user/login
<b>desc</b>: 用户登录
<b>note</b>: 无
<b>payload</b>
```
{
    "name": string,  // 用户名
    "email": string,  // 邮箱
    "phone": string,  // 手机号
    "passwd": string, // 密码
    "captcha_session": string, // 验证码session
    "captcha_value": string, // 验证码值
}
```
<b>return</b>
```
{
    "uid": string,  // 用户id
    "session": string, // 会话id
    "token": string,  // 会话令牌
}
```