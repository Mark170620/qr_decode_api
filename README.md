### 用Golang搭建一个解析二维码的api服务
- 框架：fiber

- 支持POST方式上传二维码图片解析

- 支持GET方式传送图片URL地址解析

- 支持一图多码解析

### 解码返回内容

```json
{
"code": 200,
"data": {
"qrText": [
"https://xxxx",
"https://xxxx"
]
},
"msg": "Decode Success！"
}
```