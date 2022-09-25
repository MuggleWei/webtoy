# webtoy-auth
用户认证服务, 也充当用户信息服务

## 接口

### 用户认证
<b>url</b>: /user/auth
<b>desc</b>: 用户认证
<b>note</b>: 无
<b>payload</b>
```
{
    "name": string,  // 用户名
    "email": string,  // 邮箱
    "phone": string,  // 手机号
    "passwd": string, // 密码
}
```
<b>return</b>
```
{
    "uid": string,  // 用户id
}
```