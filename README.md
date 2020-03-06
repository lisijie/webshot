# 网页截图服务

提供网页截图服务，接收一个URL，返回该URL的截图。

## 运行

```
$ docker run -d --shm-size 2G --name webshot --rm -p 8080:80 lixiaoxie/webshot
```

## 使用

打开浏览器，访问 [http://localhost:8080/?url=https://www.baidu.com&size=500x500]() 



