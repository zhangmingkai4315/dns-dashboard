# dns-dashboard
dns-dashboard 提供基于DNS日志的信息统计与展示平台


### 1. DNS服务运行环境部署(仅用于演示)

所有日志将展示在docker-envirments/logs目录下，该项目将分析该日志的内容，并将结果展示在web系统上

```
cd docker-envirments
docker build -t dns-bind-log .
docker run -it -d --name=bind9.9 -v "$PWD/log/dnslog":/var/cache/bind/dnslog dns-bind-log:latest 
```

测试,请先查询启动的docker的ip地址
```
> docker inspect -f "{{ .NetworkSettings.IPAddress}}" bind9.9
> 172.17.0.2
```
发送dns请求，并检查`docker-envirments/log/dnslog`目录下的query_log文件是否有日志输出情况

```
dig @172.17.0.2 ns

; <<>> DiG 9.10.3-P4-Ubuntu <<>> @172.17.0.2 ns
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 52193
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 2
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;.				IN	NS

;; ANSWER SECTION:
.			86400	IN	NS	a.dns.test.

;; ADDITIONAL SECTION:
a.dns.test.		86400	IN	A	192.168.0.2

;; Query time: 0 msec
;; SERVER: 172.17.0.2#53(172.17.0.2)
;; WHEN: Tue Nov 21 10:20:02 CST 2017
;; MSG SIZE  rcvd: 67

```


### 2. DNS发包测试

请参考链接 [dns-loader](https://github.com/zhangmingkai4315/dns-loader)

### 3. Dashboard展示



