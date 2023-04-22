### k8s-client获取容器的日志  

```bash
k apply -f service-account.yaml         创建服务账号
k apply -f cluster-role.yaml            创建Role角色(这个角色是Cluster类型的)
k apply -f cluster-role-bind.yaml       将账号和Role绑定起来
```     
main.go 逻辑:
- 获取当前服务账号访问 API-Server 的 Token
- 从指定的端口获取特定容器的日志