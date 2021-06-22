### 用Golang搭建一个解析二维码的api服务
- 框架：fiber

- 支持POST方式上传二维码解析 url/file

- 支持GET方式传送二维码URL地址解析 url/url

- 支持一图多码解析二维码

###请求示例
`curl http://127.0.0.1:9900/url\?http_url\=http://cdn.layui.com/upload/2021_4/168_1617643238223_53889.png`

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