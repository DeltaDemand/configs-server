# configs-server
> refresh agents' configs in real time

拉取configs-server

```bash
docker build -t config-server https://github.com/DeltaDemand/configs-server.git#main
```
```bash
#开放8887端口
docker run -name config-s01 -p 8887:8887 -restart=always config-server -name=config-s01
```

### 修改配置前端页面

使用swagger实现接口文档以及测试，启动后访问以下链接即可：

```
#本地
localhost:8887/swagger/index.html

#阿里云
http://112.74.60.132:8887/swagger/index.html
```

点击`try it out`即可测试