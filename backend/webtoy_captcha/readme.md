# webtoy-captcha
图形验证码服务, 用于图形验证码的生成与校验

## 接口

### 生成
<b>url</b>: /captcha/load
<b>desc</b>: 生成图形验证码
<b>note</b>: 无
<b>url param</b>
| 字段名 | 类型 | 意义 | 取值 | 备注 |
| ---- | ---- | ---- | ---- | ---- |
| captcha_session | string | 当值不为空时, 更新captcha session的值, 而不会新生成. 这样防止用户频繁刷新在redis中生成许多captcha session | | |

### 验证
<b>url</b>: /captcha/verify
<b>desc</b>: 验证图形验证码
<b>note</b>: 无
<b>payload</b>
```
{
	"captcha_session": string,  // captcha session值
	"captcha_value": string,  // captcha value值
}
```
<b>return</b>
```
仅包含通用结构
```
