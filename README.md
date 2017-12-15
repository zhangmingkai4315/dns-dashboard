# dns-dashboard
dns-dashboard 提供基于DNS日志的信息统计与展示平台


### 1. DNS服务运行环境部署

(仅用于演示,用于开发或者测试使用)

```
cd docker-envirments
docker build -t dns-bind-log .
docker run -it -d --name=bind9.9 -v "$PWD/log/dnslog":/var/cache/bind/dnslog dns-bind-log
```
所有日志将展示在docker-envirments/logs目录下

查询启动的docker的ip地址
```
> docker inspect -f "{{ .NetworkSettings.IPAddress}}" bind9.9
> 172.17.0.2
```
发送dns请求，并检查`docker-envirments/log/dnslog`目录下的query_log文件是否有日志输出情况

```
dig @172.17.0.2 ns +short
a.dns.test.
```


### 2. DNS发包测试

请参考链接 [dns-loader](https://github.com/zhangmingkai4315/dns-loader)

### 3. Dashboard展示

请手动修改`config.ini`配置中加载日志的路径DNS:Path和解释日志的DNS:grok,并确保grok中至少包含domain, ip和type三个字段

比如下面的针对DNS bind日志的解释grok格式：

```
^%{NOTSPACE:date} %{TIME:time} queries: info: client %{IPORHOST:ip}%{GREEDYDATA:message} query: %{NOTSPACE:domain} IN %{WORD:type} %{NOTSPACE:message2} \(%{IPORHOST:server}\)$

```

配置完成后即可直接运行

```
./dns-dashboard -c config.ini
```
可以通过配置文件修改IP的监听端口和IP地址，默认情况下仅仅在本地可访问，打开浏览器访问http://localhost:9889/ 



